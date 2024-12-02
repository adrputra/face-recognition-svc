package model

import "time"

type MenuRoleMapping struct {
	Id           string    `gorm:"column:id" json:"id"`
	MenuID       string    `gorm:"column:menu_id" json:"menu_id"`
	RoleID       string    `gorm:"column:role_id" json:"role_id"`
	RoleName     string    `gorm:"column:role_name" json:"role_name"`
	MenuName     string    `gorm:"column:menu_name" json:"menu_name"`
	MenuRoute    string    `gorm:"column:menu_route" json:"menu_route"`
	AccessMethod string    `gorm:"column:access_method" json:"access_method"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
	CreatedBy    string    `gorm:"column:created_by" json:"created_by"`
	UpdatedBy    string    `gorm:"column:updated_by" json:"updated_by"`
}

type Menu struct {
	Id        string    `gorm:"column:id" json:"id"`
	MenuName  string    `gorm:"column:menu_name" json:"menu_name"`
	MenuRoute string    `gorm:"column:menu_route" json:"menu_route"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	CreatedBy string    `gorm:"column:created_by" json:"created_by"`
	UpdatedBy string    `gorm:"column:updated_by" json:"updated_by"`
}

type Role struct {
	Id        string    `gorm:"column:id" json:"id"`
	RoleName  string    `gorm:"column:role_name" json:"role_name"`
	RoleDesc  string    `gorm:"column:role_desc" json:"role_desc"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	CreatedBy string    `gorm:"column:created_by" json:"created_by"`
	UpdatedBy string    `gorm:"column:updated_by" json:"updated_by"`
	IsActive  bool      `gorm:"column:is_active" json:"is_active"`
}
