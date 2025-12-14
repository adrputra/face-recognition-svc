package client

import (
	"context"
	"errors"
	"face-recognition-svc/app/config"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type InterfaceUserClient interface {
	CreateNewUser(ctx context.Context, user *model.User) error
	GetUserDetail(ctx context.Context, username string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, username string) error
	CreateAccessToken(ctx context.Context, user *model.User, isLogout bool, menuMapping map[string]string) (t string, expired int64, err error)
	GetAllUser(ctx context.Context, roleLevel int, institutionID string, pagination *model.Pagination, filter *model.Filter) ([]*model.User, *model.Pagination, error)
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

	var args []interface{}
	args = append(args, req.Username, req.Email, req.Password, req.Fullname, req.Shortname, req.RoleID, req.InstitutionID, utils.LocalTime(), req.Address, req.PhoneNumber, req.Gender, req.Religion)

	query := "INSERT INTO users (username, email, password, fullname, shortname, role_id, institution_id, created_at, address, phone_number, gender, religion) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result := r.db.Debug().WithContext(ctx).Exec(query, args...)

	if result.Error != nil {
		if mysqlErr, ok := result.Error.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062: // Duplicate entry
				utils.LogEventError(span, errors.New("username or email already exists"))
				return model.ThrowError(http.StatusBadRequest, errors.New("username or email already exists"))
			}
		}
		utils.LogEventError(span, result.Error)
		return model.ThrowError(http.StatusInternalServerError, result.Error)
	}

	return nil
}

func (r *UserClient) GetUserDetail(ctx context.Context, username string) (*model.User, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetUserDetail")
	defer span.Finish()

	utils.LogEvent(span, "Request", username)

	var user model.User

	query := "SELECT u.*, i.name AS institution_name, r.role_name FROM users AS u LEFT JOIN institution AS i ON u.institution_id = i.id LEFT JOIN role AS r ON u.role_id = r.id WHERE username = ?"
	result := r.db.Debug().WithContext(ctx).Raw(query, username).Scan(&user)

	if result.RowsAffected == 0 {
		utils.LogEventError(span, errors.New("user not found"))
		return nil, model.ThrowError(http.StatusBadRequest, errors.New("user not found"))
	}

	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return nil, model.ThrowError(http.StatusInternalServerError, result.Error)
	}

	utils.LogEvent(span, "Response", user)

	return &user, nil
}

func (r *UserClient) UpdateUser(ctx context.Context, user *model.User) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: UpdateUser")
	defer span.Finish()

	utils.LogEvent(span, "Request", user)

	var args []interface{}
	args = append(args, user.Fullname, user.Shortname, user.Email, user.RoleID, user.InstitutionID, user.Address, user.PhoneNumber, user.Gender, user.Religion, user.Username)

	query := "UPDATE users SET fullname = ?, shortname = ?, email = ?, role_id = ?, institution_id = ?, address = ?, phone_number = ?, gender = ?, religion = ? WHERE username = ?"
	result := r.db.Debug().WithContext(ctx).Exec(query, args...)

	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return model.ThrowError(http.StatusInternalServerError, result.Error)
	}

	// if result.RowsAffected == 0 {
	// 	utils.LogEventError(span, errors.New("user not found"))
	// 	return model.ThrowError(http.StatusBadRequest, errors.New("user not found"))
	// }

	return nil
}

func (r *UserClient) DeleteUser(ctx context.Context, username string) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: DeleteUser")
	defer span.Finish()

	utils.LogEvent(span, "Request", username)

	query := "DELETE FROM users WHERE username = ?"
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

func (r *UserClient) CreateAccessToken(ctx context.Context, user *model.User, isLogout bool, menuMapping map[string]string) (t string, expired int64, err error) {
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
		Name: user.Username,
		Role: user.RoleID,
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

func (r *UserClient) GetAllUser(ctx context.Context, roleLevel int, institutionID string, pagination *model.Pagination, filter *model.Filter) ([]*model.User, *model.Pagination, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetAllUser")
	defer span.Finish()

	var response []*model.User

	// Build WHERE clause for role level
	whereConditions := []string{}
	if roleLevel == 2 {
		whereConditions = append(whereConditions, fmt.Sprintf("u.institution_id = '%s'", institutionID))
	}

	// Build WHERE clause for search
	searchFields := []string{"u.username", "u.fullname", "u.shortname", "u.email", "i.name", "r.role_name"}
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
	countQuery := "SELECT COUNT(*) FROM users AS u LEFT JOIN institution AS i ON u.institution_id = i.id LEFT JOIN role AS r ON u.role_id = r.id" + whereClause
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
		"fullname":         "u.fullname",
		"shortname":        "u.shortname",
		"email":            "u.email",
		"institution_name": "i.name",
		"role_name":        "r.role_name",
		"created_at":       "u.created_at",
	}
	orderByClause := utils.BuildOrderByClause(filter, allowedSortFields, "u.username")

	// Build query with pagination
	sb := strings.Builder{}
	sb.WriteString(whereClause)
	sb.WriteString(orderByClause)
	if pagination != nil {
		sb.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", pagination.Limit, (pagination.Page-1)*pagination.Limit))
	}

	query := "SELECT u.username, u.fullname, u.shortname, u.email, u.institution_id, u.role_id, u.address, u.phone_number, u.gender, u.religion, u.created_at, i.name AS institution_name, r.role_name FROM users AS u LEFT JOIN institution AS i ON u.institution_id = i.id LEFT JOIN role AS r ON u.role_id = r.id"
	result := r.db.Debug().WithContext(ctx).Raw(query + sb.String()).Scan(&response)

	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return nil, nil, model.ThrowError(http.StatusInternalServerError, result.Error)
	}

	utils.LogEvent(span, "Response", response)
	utils.LogEvent(span, "Pagination", pagination)

	return response, pagination, nil
}

func (r *UserClient) GetInstitutionList(ctx context.Context) ([]string, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetInstitutionList")
	defer span.Finish()

	var response []string

	query := "SELECT DISTINCT institution_id FROM users"
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
	query := "UPDATE users SET profile_photo = ? WHERE username = ?"

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
	query := "UPDATE users SET cover_photo = ? WHERE username = ?"

	result = r.db.Debug().WithContext(ctx).Exec(query, args...)

	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return result.Error
	}

	return nil
}
