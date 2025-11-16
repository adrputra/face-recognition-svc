package model

type Institution struct {
	ID          string `json:"id" gorm:"type:varchar(200);"`
	Name        string `json:"name" gorm:"type:varchar(200);"`
	Address     string `json:"address" gorm:"type:varchar(200);"`
	PhoneNumber string `json:"phone_number" gorm:"type:varchar(200);"`
	Email       string `json:"email" gorm:"type:varchar(200);"`
	CreatedAt   string `json:"created_at" gorm:"type:varchar(200);"`
	CreatedBy   string `json:"created_by" gorm:"type:varchar(200);"`
	UpdatedAt   string `json:"updated_at" gorm:"type:varchar(200);"`
	UpdatedBy   string `json:"updated_by" gorm:"type:varchar(200);"`
}
