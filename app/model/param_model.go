package model

import "time"

type Param struct {
	Key         string    `json:"key" gorm:"column:id"`
	Value       string    `json:"value" gorm:"column:value"`
	Description string    `json:"description" gorm:"column:description"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy   string    `json:"updated_by" gorm:"column:updated_by"`
}
