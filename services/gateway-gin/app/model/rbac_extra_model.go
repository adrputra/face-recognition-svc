package model

import "time"

type Permission struct {
	ID          string    `json:"id" gorm:"column:id"`
	Name        string    `json:"name" gorm:"column:name"`
	Service     string    `json:"service" gorm:"column:service"`
	Resource    string    `json:"resource" gorm:"column:resource"`
	Action      string    `json:"action" gorm:"column:action"`
	IsActive    bool      `json:"is_active" gorm:"column:is_active"`
	IsHighRisk  bool      `json:"is_high_risk" gorm:"column:is_high_risk"`
	Description string    `json:"description" gorm:"column:description"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Permission) TableName() string {
	return "permission"
}

type RolePermission struct {
	ID           string    `json:"id" gorm:"column:id"`
	RoleID       string    `json:"role_id" gorm:"column:role_id"`
	PermissionID string    `json:"permission_id" gorm:"column:permission_id"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
}

func (RolePermission) TableName() string {
	return "role_permission"
}

type Feature struct {
	ID             string    `json:"id" gorm:"column:id"`
	FeatureKey     string    `json:"feature_key" gorm:"column:feature_key"`
	Name           string    `json:"name" gorm:"column:name"`
	Description    string    `json:"description" gorm:"column:description"`
	FeatureType    string    `json:"feature_type" gorm:"column:feature_type"`
	DefaultEnabled bool      `json:"default_enabled" gorm:"column:default_enabled"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Feature) TableName() string {
	return "feature"
}

type InstitutionFeature struct {
	ID            string    `json:"id" gorm:"column:id"`
	InstitutionID string    `json:"institution_id" gorm:"column:institution_id"`
	FeatureKey    string    `json:"feature_key" gorm:"column:feature_key"`
	IsEnabled     bool      `json:"is_enabled" gorm:"column:is_enabled"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (InstitutionFeature) TableName() string {
	return "institution_feature"
}

type AuditLog struct {
	ID             string    `json:"id" gorm:"column:id"`
	ActorUserID    *string   `json:"actor_user_id" gorm:"column:actor_user_id"`
	InstitutionID  *string   `json:"institution_id" gorm:"column:institution_id"`
	PermissionName *string   `json:"permission_name" gorm:"column:permission_name"`
	Action         string    `json:"action" gorm:"column:action"`
	EntityType     string    `json:"entity_type" gorm:"column:entity_type"`
	EntityID       string    `json:"entity_id" gorm:"column:entity_id"`
	RequestID      string    `json:"request_id" gorm:"column:request_id"`
	IPAddress      string    `json:"ip_address" gorm:"column:ip_address"`
	UserAgent      string    `json:"user_agent" gorm:"column:user_agent"`
	Metadata       string    `json:"metadata" gorm:"column:metadata"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`
}

func (AuditLog) TableName() string {
	return "audit_log"
}

type RolePermissionAssignment struct {
	RoleID        string   `json:"role_id"`
	PermissionIDs []string `json:"permission_ids"`
}

type PermissionUpdateRequest struct {
	ID          string  `json:"id"`
	IsActive    *bool   `json:"is_active"`
	IsHighRisk  *bool   `json:"is_high_risk"`
	Description *string `json:"description"`
}

type InstitutionFeatureRequest struct {
	InstitutionID string `json:"institution_id"`
	FeatureKey    string `json:"feature_key"`
	IsEnabled     bool   `json:"is_enabled"`
}
