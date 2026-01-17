package client

import (
	"context"
	"errors"
	"face-recognition-svc/gateway/app/config"
	"face-recognition-svc/gateway/app/model"
	"face-recognition-svc/gateway/app/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type InterfaceUserClient interface {
	CreateNewUser(ctx context.Context, user *model.User) error
	GetUserDetail(ctx context.Context, username string, institutionID string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, username string) error
	CreateAccessToken(ctx context.Context, user *model.User, isLogout bool) (t string, expired int64, err error)
	GetAllUser(ctx context.Context, scope string, institutionID string, pagination *model.Pagination, filter *model.Filter) ([]*model.User, *model.Pagination, error)
	GetInstitutionList(ctx context.Context) ([]string, error)
	UpdateProfilePhoto(ctx context.Context, url string, username string) error
	UpdateCoverPhoto(ctx context.Context, url string, username string) error
}

type UserClient struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewUserClient(db *gorm.DB, cfg *config.Config) *UserClient {
	return &UserClient{
		db:  db,
		cfg: cfg,
	}
}

func (r *UserClient) CreateNewUser(ctx context.Context, req *model.User) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: CreateNewUser")
	defer span.Finish()

	utils.LogEvent(span, "Request", req)

	if req.InstitutionID == "" {
		return model.ThrowError(http.StatusBadRequest, errors.New("institution_id is required"))
	}
	roleIDs := req.RoleIDs
	if len(roleIDs) == 0 && req.RoleID != "" {
		roleIDs = []string{req.RoleID}
	}
	if len(roleIDs) == 0 {
		return model.ThrowError(http.StatusBadRequest, errors.New("role assignment is required"))
	}

	createdAt := utils.LocalTime()
	var userID string
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		query := `
			INSERT INTO "user" (username, email, password_hash, full_name, short_name, is_active, created_at, updated_at, profile_photo, cover_photo)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			RETURNING id`
		if err := tx.Raw(query, req.Username, req.Email, req.PasswordHash, req.Fullname, req.Shortname, true, createdAt, createdAt, req.ProfilePhoto, req.CoverPhoto).Scan(&userID).Error; err != nil {
			return err
		}

		query = `
			INSERT INTO user_institution (user_id, institution_id, status, joined_at, created_at, updated_at)
			VALUES (?, ?, 'active', ?, ?, ?)`
		if err := tx.Exec(query, userID, req.InstitutionID, createdAt, createdAt, createdAt).Error; err != nil {
			return err
		}

		for _, roleID := range roleIDs {
			query = `
				INSERT INTO user_role (user_id, institution_id, role_id, created_at)
				VALUES (?, ?, ?, ?)`
			if err := tx.Exec(query, userID, req.InstitutionID, roleID, createdAt).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		utils.LogEventError(span, err)
		return model.ThrowError(http.StatusInternalServerError, err)
	}

	return nil
}

func (r *UserClient) GetUserDetail(ctx context.Context, username string, institutionID string) (*model.User, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetUserDetail")
	defer span.Finish()

	utils.LogEvent(span, "Request", username)

	var user model.User

	query := `
		SELECT u.id, u.username, u.email, u.full_name, u.short_name, u.is_active,
			u.profile_photo, u.cover_photo, u.created_at, u.updated_at,
			i.id AS institution_id, i.name AS institution_name
		FROM "user" u
		LEFT JOIN user_institution ui ON ui.user_id = u.id
		LEFT JOIN institution i ON i.id = ui.institution_id
		WHERE u.username = ?`
	args := []interface{}{username}
	if institutionID != "" {
		query += " AND ui.institution_id = ?"
		args = append(args, institutionID)
	}
	result := r.db.Debug().WithContext(ctx).Raw(query, args...).Scan(&user)

	if result.RowsAffected == 0 {
		utils.LogEventError(span, errors.New("user not found"))
		return nil, model.ThrowError(http.StatusBadRequest, errors.New("user not found"))
	}

	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return nil, model.ThrowError(http.StatusInternalServerError, result.Error)
	}

	if institutionID != "" {
		var roleIDs []string
		var roleName string
		roleQuery := `
			SELECT r.id, r.name
			FROM role r
			JOIN user_role ur ON ur.role_id = r.id
			WHERE ur.user_id = ? AND ur.institution_id = ?`
		rows, err := r.db.WithContext(ctx).Raw(roleQuery, user.ID, institutionID).Rows()
		if err != nil {
			utils.LogEventError(span, err)
			return nil, model.ThrowError(http.StatusInternalServerError, err)
		}
		defer rows.Close()
		for rows.Next() {
			var id string
			var name string
			if err := rows.Scan(&id, &name); err != nil {
				utils.LogEventError(span, err)
				return nil, model.ThrowError(http.StatusInternalServerError, err)
			}
			roleIDs = append(roleIDs, id)
			if roleName == "" {
				roleName = name
				user.RoleID = id
				user.RoleName = name
			}
		}
		user.RoleIDs = roleIDs
	}

	utils.LogEvent(span, "Response", user)

	return &user, nil
}

func (r *UserClient) UpdateUser(ctx context.Context, user *model.User) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: UpdateUser")
	defer span.Finish()

	utils.LogEvent(span, "Request", user)

	query := "UPDATE \"user\" SET full_name = ?, short_name = ?, email = ?, is_active = ?, updated_at = ? WHERE username = ?"
	result := r.db.Debug().WithContext(ctx).Exec(query, user.Fullname, user.Shortname, user.Email, user.IsActive, utils.LocalTime(), user.Username)
	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return model.ThrowError(http.StatusInternalServerError, result.Error)
	}

	if user.InstitutionID != "" && len(user.RoleIDs) > 0 {
		roleIDs := user.RoleIDs
		err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var userID string
			if err := tx.Raw("SELECT id FROM \"user\" WHERE username = ?", user.Username).Scan(&userID).Error; err != nil {
				return err
			}
			createdAt := utils.LocalTime()
			ensureQuery := `
				INSERT INTO user_institution (user_id, institution_id, status, joined_at, created_at, updated_at)
				VALUES (?, ?, 'active', ?, ?, ?)
				ON CONFLICT (user_id, institution_id) DO NOTHING`
			if err := tx.Exec(ensureQuery, userID, user.InstitutionID, createdAt, createdAt, createdAt).Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM user_role WHERE user_id = ? AND institution_id = ?", userID, user.InstitutionID).Error; err != nil {
				return err
			}
			for _, roleID := range roleIDs {
				insertQuery := `
					INSERT INTO user_role (user_id, institution_id, role_id, created_at)
					VALUES (?, ?, ?, ?)`
				if err := tx.Exec(insertQuery, userID, user.InstitutionID, roleID, createdAt).Error; err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			utils.LogEventError(span, err)
			return model.ThrowError(http.StatusInternalServerError, err)
		}
	}

	return nil
}

func (r *UserClient) DeleteUser(ctx context.Context, username string) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: DeleteUser")
	defer span.Finish()

	utils.LogEvent(span, "Request", username)

	query := "DELETE FROM \"user\" WHERE username = ?"
	result := r.db.Debug().WithContext(ctx).Exec(query, username)

	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return model.ThrowError(http.StatusInternalServerError, result.Error)
	}

	if result.RowsAffected == 0 {
		utils.LogEventError(span, errors.New("user not found"))
		return model.ThrowError(http.StatusBadRequest, errors.New("user not found"))
	}

	return nil
}

func (r *UserClient) CreateAccessToken(ctx context.Context, user *model.User, isLogout bool) (t string, expired int64, err error) {
	span, _ := utils.SpanFromContext(ctx, "Client: CreateAccessToken")
	defer span.Finish()

	utils.LogEvent(span, "Request", user)

	ExpireCount, _ := strconv.Atoi(r.cfg.Auth.AccessExpiry)
	if isLogout {
		ExpireCount = 0
	}

	utils.LogEvent(span, "Expiry", ExpireCount)

	exp := utils.LocalTime().Add(time.Hour * time.Duration(ExpireCount))
	claims := &model.JwtCustomClaims{
		UserID:   user.ID,
		Username: user.Username,
		RoleIDs:  user.RoleIDs,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		InstitutionID: user.InstitutionID,
	}
	expired = exp.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err = token.SignedString([]byte(r.cfg.Auth.AccessSecret))
	if err != nil {
		utils.LogEventError(span, err)
		return "", 0, err
	}

	utils.LogEvent(span, "Token", t)

	return t, expired, nil
}

func (r *UserClient) GetAllUser(ctx context.Context, scope string, institutionID string, pagination *model.Pagination, filter *model.Filter) ([]*model.User, *model.Pagination, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetAllUser")
	defer span.Finish()

	var response []*model.User

	// Build WHERE clause for role level
	whereConditions := []string{}
	if scope != "system" {
		whereConditions = append(whereConditions, fmt.Sprintf("ui.institution_id = '%s'", institutionID))
	}

	// Build WHERE clause for search
	searchFields := []string{"u.username", "u.full_name", "u.short_name", "u.email", "i.name"}
	searchClause := utils.BuildSearchWhereClause(filter.Search, searchFields)
	if searchClause != "" {
		// Remove " WHERE " prefix and add condition
		searchCondition := strings.TrimPrefix(searchClause, " WHERE ")
		whereConditions = append(whereConditions, "("+searchCondition+")")
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Get total count for pagination
	var totalCount int64
	countQuery := "SELECT COUNT(*) FROM user_institution ui JOIN \"user\" u ON u.id = ui.user_id JOIN institution i ON i.id = ui.institution_id" + whereClause
	countResult := r.db.Debug().WithContext(ctx).Raw(countQuery).Scan(&totalCount)
	if countResult.Error != nil {
		utils.LogEventError(span, countResult.Error)
		return nil, nil, model.ThrowError(http.StatusInternalServerError, countResult.Error)
	}

	// Calculate total pages
	pagination.Total = int(totalCount)
	if pagination.Limit > 0 {
		pagination.TotalPages = (pagination.Total + pagination.Limit - 1) / pagination.Limit
	} else {
		pagination.TotalPages = 1
	}

	// Define allowed sort fields with their actual column names
	allowedSortFields := map[string]string{
		"username":         "u.username",
		"fullname":         "u.full_name",
		"shortname":        "u.short_name",
		"email":            "u.email",
		"institution_name": "i.name",
		"created_at":       "u.created_at",
	}
	orderByClause := utils.BuildOrderByClause(filter, allowedSortFields, "u.username")

	query := `
		SELECT u.id, u.username, u.full_name, u.short_name, u.email, u.is_active, u.created_at,
			i.id AS institution_id, i.name AS institution_name,
			COALESCE(array_agg(r.id::text) FILTER (WHERE r.id IS NOT NULL), '{}') AS role_ids
		FROM "user" u
		JOIN user_institution ui ON ui.user_id = u.id
		JOIN institution i ON i.id = ui.institution_id
		LEFT JOIN user_role ur ON ur.user_id = u.id AND ur.institution_id = ui.institution_id
		LEFT JOIN role r ON r.id = ur.role_id` + whereClause + " GROUP BY u.id, i.id"

	type userRow struct {
		ID              string
		Username        string
		Fullname        string
		Shortname       string
		Email           string
		IsActive        bool
		CreatedAt       time.Time
		InstitutionID   string
		InstitutionName string
		RoleIDs         pq.StringArray `gorm:"type:text[]"`
	}

	var rows []userRow
	queryArgs := query + orderByClause
	if pagination != nil {
		queryArgs += fmt.Sprintf(" LIMIT %d OFFSET %d", pagination.Limit, (pagination.Page-1)*pagination.Limit)
	}
	result := r.db.Debug().WithContext(ctx).Raw(queryArgs).Scan(&rows)

	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return nil, nil, model.ThrowError(http.StatusInternalServerError, result.Error)
	}

	for _, row := range rows {
		user := &model.User{
			ID:              row.ID,
			Username:        row.Username,
			Fullname:        row.Fullname,
			Shortname:       row.Shortname,
			Email:           row.Email,
			IsActive:        row.IsActive,
			CreatedAt:       row.CreatedAt,
			InstitutionID:   row.InstitutionID,
			InstitutionName: row.InstitutionName,
			RoleIDs:         []string(row.RoleIDs),
		}
		if len(user.RoleIDs) > 0 {
			user.RoleID = user.RoleIDs[0]
		}
		response = append(response, user)
	}

	utils.LogEvent(span, "Response", response)
	utils.LogEvent(span, "Pagination", pagination)

	return response, pagination, nil
}

func (r *UserClient) GetInstitutionList(ctx context.Context) ([]string, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetInstitutionList")
	defer span.Finish()

	var response []string

	query := "SELECT DISTINCT institution_id FROM user_institution"
	result := r.db.Debug().WithContext(ctx).Raw(query).Scan(&response)

	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return nil, model.ThrowError(http.StatusInternalServerError, result.Error)
	}

	utils.LogEvent(span, "Response", response)

	return response, nil
}

func (r *UserClient) UpdateProfilePhoto(ctx context.Context, url string, username string) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: UpdateProfilePhoto")
	defer span.Finish()

	var args []interface{}
	args = append(args, url, username)

	var result *gorm.DB
	query := "UPDATE \"user\" SET profile_photo = ? WHERE username = ?"

	result = r.db.Debug().WithContext(ctx).Exec(query, args...)

	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return result.Error
	}

	return nil
}

func (r *UserClient) UpdateCoverPhoto(ctx context.Context, url string, username string) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: UpdateCoverPhoto")
	defer span.Finish()

	var args []interface{}
	args = append(args, url, username)

	var result *gorm.DB
	query := "UPDATE \"user\" SET cover_photo = ? WHERE username = ?"

	result = r.db.Debug().WithContext(ctx).Exec(query, args...)

	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return result.Error
	}

	return nil
}
