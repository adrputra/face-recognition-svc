package client

import (
	"context"
	"errors"
	"face-recognition-svc/app/config"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type InterfaceUserClient interface {
	CreateNewUser(ctx context.Context, user *model.User) error
	GetUserDetail(ctx context.Context, username string) (*model.User, error)
	CreateAccessToken(ctx context.Context, user *model.User, isLogout bool, menuMapping map[string]string) (t string, expired int64, err error)
	GetAllUser(ctx context.Context) ([]*model.User, error)
	GetInstitutionList(ctx context.Context) ([]string, error)
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
	args = append(args, req.Username, req.Email, req.Password, req.Fullname, req.Shortname, req.RoleID, req.InstitutionID, time.Now())

	query := "INSERT INTO users (username, email, password, fullname, shortname, role_id, institution_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
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

	query := "SELECT * FROM users WHERE username = ?"
	result := r.db.Debug().WithContext(ctx).Raw(query, username).Scan(&user)

	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return nil, model.ThrowError(http.StatusInternalServerError, result.Error)
	}

	if result.RowsAffected == 0 {
		utils.LogEventError(span, errors.New("user not found"))
		return nil, model.ThrowError(http.StatusBadRequest, errors.New("user not found"))
	}

	utils.LogEvent(span, "Response", user)

	return &user, nil
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

	exp := time.Now().Add(time.Hour * time.Duration(ExpireCount))
	claims := &model.JwtCustomClaims{
		Name: user.Username,
		Role: user.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		MenuMapping: menuMapping,
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

func (r *UserClient) GetAllUser(ctx context.Context) ([]*model.User, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetAllUser")
	defer span.Finish()

	var response []*model.User

	query := "SELECT * FROM users"
	result := r.db.Debug().WithContext(ctx).Raw(query).Scan(&response)

	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return nil, model.ThrowError(http.StatusInternalServerError, result.Error)
	}

	utils.LogEvent(span, "Response", response)

	return response, nil
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
