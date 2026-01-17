package model

import "time"

type Institution struct {
	ID          string    `json:"id" gorm:"column:id"`
	Name        string    `json:"name" gorm:"column:name"`
	Code        string    `json:"code" gorm:"column:code"`
	Address     string    `json:"address" gorm:"column:address"`
	PhoneNumber string    `json:"phone_number" gorm:"column:phone_number"`
	Email       string    `json:"email" gorm:"column:email"`
	IsActive    bool      `json:"is_active" gorm:"column:is_active"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// TableName specifies the table name for Institution model
func (Institution) TableName() string {
	return "institution"
}
