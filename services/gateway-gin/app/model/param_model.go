package model

import "time"

type Param struct {
	Key         string     `json:"key" gorm:"column:id"`
	Value       string     `json:"value" gorm:"column:value"`
	Description string     `json:"description" gorm:"column:description"`
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	CreatedBy   string     `json:"created_by" gorm:"column:created_by"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedBy   string     `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:timestamp;index"`
	DeletedBy   *string    `json:"deleted_by" gorm:"column:deleted_by"`
}
