package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Name        string            `json:"name"`
	Role        string            `json:"role"`
	MenuMapping map[string]string `json:"menu_mapping"`
	jwt.RegisteredClaims
	InstitutionID string `json:"institution_id"`
}

type MetadataUser struct {
	Username      string `json:"username"`
	RoleID        string `json:"role_id"`
	InstitutionID string `json:"institution_id"`
}

type User struct {
	Username        string `json:"username" gorm:"column:username" validate:"required"`
	Email           string `json:"email" gorm:"column:email" validate:"required"`
	Password        string `json:"password" gorm:"column:password" validate:"required"`
	Fullname        string `json:"fullname" gorm:"column:fullname" validate:"required"`
	Shortname       string `json:"shortname" gorm:"column:shortname" validate:"required"`
	RoleID          string `json:"role_id" gorm:"column:role_id" validate:"required"`
	RoleName        string `json:"role_name" gorm:"column:role_name"`
	InstitutionID   string `json:"institution_id" gorm:"column:institution_id" validate:"required"`
	InstitutionName string `json:"institution_name" gorm:"column:institution_name"`
	PhoneNumber     string `json:"phone_number" gorm:"column:phone_number"`
	Address         string `json:"address" gorm:"column:address"`
	Gender          string `json:"gender" gorm:"column:gender"`
	Religion        string `json:"religion" gorm:"column:religion"`
	CreatedAt       string `json:"created_at" gorm:"column:created_at"`
	RoleLevel       int    `json:"role_level" gorm:"-"`
	ProfilePhoto    string `json:"profile_photo" gorm:"column:profile_photo"`
	CoverPhoto      string `json:"cover_photo" gorm:"column:cover_photo"`
}

type RequestLogin struct {
	Username string `json:"username" gorm:"column:username" validate:"required"`
	Password string `json:"password" gorm:"column:password" validate:"required"`
}

type ResponseLogin struct {
	Username        string             `json:"username" gorm:"type:varchar(200);"`
	Fullname        string             `json:"fullname" gorm:"type:varchar(200);"`
	Shortname       string             `json:"shortname" gorm:"type:varchar(200);"`
	Role            string             `json:"role" gorm:"type:varchar(200);"`
	RoleName        string             `json:"role_name" gorm:"type:varchar(200);"`
	Token           string             `json:"token" gorm:"type:varchar(200);"`
	InstitutionID   string             `json:"institution_id" gorm:"type:varchar(200);"`
	InstitutionName string             `json:"institution_name" gorm:"type:varchar(200);"`
	MenuMapping     []*MenuRoleMapping `json:"menu_mapping" gorm:"-"`
}

type UploadPhoto struct {
	Photo File `json:"photo"`
}
