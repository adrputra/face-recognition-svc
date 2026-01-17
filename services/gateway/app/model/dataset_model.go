package model

import "time"

type Dataset struct {
	ID        string     `json:"id" gorm:"column:id"`
	Username  string     `json:"username" gorm:"column:username" validate:"required"`
	Bucket    string     `json:"bucket" gorm:"column:bucket" validate:"required"`
	Dataset   string     `json:"dataset" gorm:"column:dataset" validate:"required"`
	File      []*File    `json:"file" gorm:"-"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	CreatedBy string     `json:"created_by" gorm:"column:created_by"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedBy string     `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:timestamp;index"`
	DeletedBy *string    `json:"deleted_by" gorm:"column:deleted_by"`
}

type ModelTraining struct {
	ID            string     `json:"id" gorm:"column:id"`
	InstitutionID string     `json:"institution_id" gorm:"column:institution_id"`
	Status        string     `json:"status" gorm:"column:status"`
	IsUsed        string     `json:"is_used" gorm:"column:is_used"`
	CreatedAt     time.Time  `json:"created_at" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	CreatedBy     string     `json:"created_by" gorm:"column:created_by"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedBy     string     `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:timestamp;index"`
	DeletedBy     *string    `json:"deleted_by" gorm:"column:deleted_by"`
}

type FilterModelTraining struct {
	InstitutionID string `json:"institution_id" gorm:"column:institution_id" validate:"required"`
	Status        string `json:"status" gorm:"column:status" validate:"required"`
	IsUsed        string `json:"is_used" gorm:"column:is_used" validate:"required"`
	OrderBy       string `json:"order_by" gorm:"column:order_by" validate:"required"`
	SortType      string `json:"sort_type" gorm:"column:sort_type" validate:"required"`
}

type DatasetURL struct {
	URL string `json:"url"`
}

type RequestTrainModel struct {
	InstitutionID string `json:"institution_id"`
}

type ResponseTrainModel struct {
	ID string `json:"id"`
}

type RequestAPITrainModel struct {
	BucketName string `json:"bucket_name" validate:"required"`
	Prefix     string `json:"prefix" validate:"required"`
	CreatedBy  string `json:"created_by" validate:"required"`
	ID         string `json:"id"`
}

type ResponseAPITrainModel struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ID string `json:"id"`
	}
}
