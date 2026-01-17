package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	UserID        string   `json:"user_id"`
	Username      string   `json:"username"`
	RoleIDs       []string `json:"role_ids"`
	InstitutionID string   `json:"institution_id"`
	jwt.RegisteredClaims
}

type MetadataUser struct {
	UserID        string   `json:"user_id"`
	Username      string   `json:"username"`
	RoleIDs       []string `json:"role_ids"`
	InstitutionID string   `json:"institution_id"`
}

type User struct {
	ID              string    `json:"id" gorm:"column:id"`
	Username        string    `json:"username" gorm:"column:username" validate:"required"`
	Email           string    `json:"email" gorm:"column:email" validate:"required"`
	PasswordHash    string    `json:"-" gorm:"column:password_hash"`
	Password        string    `json:"password,omitempty" gorm:"-"`
	Fullname        string    `json:"fullname" gorm:"column:full_name" validate:"required"`
	Shortname       string    `json:"shortname" gorm:"column:short_name"`
	IsActive        bool      `json:"is_active" gorm:"column:is_active"`
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"column:updated_at"`
	ProfilePhoto    string    `json:"profile_photo" gorm:"column:profile_photo"`
	CoverPhoto      string    `json:"cover_photo" gorm:"column:cover_photo"`
	InstitutionID   string    `json:"institution_id" gorm:"-"`
	InstitutionName string    `json:"institution_name" gorm:"-"`
	RoleID          string    `json:"role_id" gorm:"-"`
	RoleName        string    `json:"role_name" gorm:"-"`
	RoleIDs         []string  `json:"role_ids" gorm:"-"`
}

func (User) TableName() string {
	return "user"
}

type UserInstitution struct {
	ID            string     `json:"id" gorm:"column:id"`
	UserID        string     `json:"user_id" gorm:"column:user_id"`
	InstitutionID string     `json:"institution_id" gorm:"column:institution_id"`
	Status        string     `json:"status" gorm:"column:status"`
	JoinedAt      time.Time  `json:"joined_at" gorm:"column:joined_at"`
	LeftAt        *time.Time `json:"left_at" gorm:"column:left_at"`
	CreatedAt     time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"column:updated_at"`
}

func (UserInstitution) TableName() string {
	return "user_institution"
}

type UserRole struct {
	ID            string    `json:"id" gorm:"column:id"`
	UserID        string    `json:"user_id" gorm:"column:user_id"`
	InstitutionID string    `json:"institution_id" gorm:"column:institution_id"`
	RoleID        string    `json:"role_id" gorm:"column:role_id"`
	AssignedBy    *string   `json:"assigned_by" gorm:"column:assigned_by"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
}

func (UserRole) TableName() string {
	return "user_role"
}

type RequestLogin struct {
	Username      string `json:"username" gorm:"column:username" validate:"required"`
	Password      string `json:"password" gorm:"column:password" validate:"required"`
	InstitutionID string `json:"institution_id" validate:"required"`
}

type ResponseLogin struct {
	UserID          string             `json:"user_id"`
	Username        string             `json:"username"`
	Fullname        string             `json:"fullname"`
	Shortname       string             `json:"shortname"`
	RoleIDs         []string           `json:"role_ids"`
	Token           string             `json:"token"`
	InstitutionID   string             `json:"institution_id"`
	InstitutionName string             `json:"institution_name"`
	MenuMapping     []*MenuRoleMapping `json:"menu_mapping"`
}

type UploadPhoto struct {
	Photo File `json:"photo"`
}
