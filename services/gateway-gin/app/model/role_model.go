package model

import "time"

type MenuRoleMapping struct {
	ID        string    `gorm:"column:id" json:"id"`
	MenuID    string    `gorm:"column:menu_id" json:"menu_id"`
	RoleID    string    `gorm:"column:role_id" json:"role_id"`
	RoleName  string    `gorm:"column:role_name" json:"role_name"`
	MenuName  string    `gorm:"column:menu_name" json:"menu_name"`
	MenuRoute string    `gorm:"column:menu_route" json:"menu_route"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

type Menu struct {
	ID         string    `gorm:"column:id" json:"id"`
	MenuKey    string    `gorm:"column:menu_key" json:"menu_key"`
	MenuName   string    `gorm:"column:name" json:"menu_name"`
	MenuRoute  string    `gorm:"column:route" json:"menu_route"`
	Icon       string    `gorm:"column:icon" json:"icon"`
	ParentID   *string   `gorm:"column:parent_id" json:"parent_id"`
	SortOrder  int       `gorm:"column:sort_order" json:"sort_order"`
	FeatureKey *string   `gorm:"column:feature_key" json:"feature_key"`
	IsActive   bool      `gorm:"column:is_active" json:"is_active"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Role struct {
	ID              string       `gorm:"column:id" json:"id"`
	RoleName        string       `gorm:"column:name" json:"role_name"`
	RoleDesc        string       `gorm:"column:description" json:"role_desc"`
	Scope           string       `gorm:"column:scope" json:"scope"`
	InstitutionID   *string      `gorm:"column:institution_id" json:"institution_id"`
	IsActive        bool         `gorm:"column:is_active" json:"is_active"`
	IsAdministrator bool         `gorm:"column:is_administrator" json:"is_administrator"`
	CreatedAt       time.Time    `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time    `gorm:"column:updated_at" json:"updated_at"`
	Institution     *Institution `gorm:"foreignKey:InstitutionID" json:"institution"`
}

// TableName specifies the table name for Role model
func (Role) TableName() string {
	return "role"
}

// TableName specifies the table name for Menu model
func (Menu) TableName() string {
	return "menu"
}

type RoleMenu struct {
	ID        string    `gorm:"column:id" json:"id"`
	RoleID    string    `gorm:"column:role_id" json:"role_id"`
	MenuID    string    `gorm:"column:menu_id" json:"menu_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (RoleMenu) TableName() string {
	return "role_menu"
}
