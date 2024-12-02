package model

type Dataset struct {
	ID        string  `json:"id" gorm:"column:id"`
	Username  string  `json:"username" gorm:"column:username" validate:"required"`
	Bucket    string  `json:"bucket" gorm:"column:bucket" validate:"required"`
	Dataset   string  `json:"dataset" gorm:"column:dataset" validate:"required"`
	File      []*File `json:"file" gorm:"-"`
	CreatedAt string  `json:"created_at" gorm:"column:created_at"`
}

type ModelTraining struct {
	ID            string `json:"id" gorm:"column:id"`
	InstitutionID string `json:"institution_id" gorm:"column:institution_id"`
	Status        string `json:"status" gorm:"column:status"`
	IsUsed        string `json:"is_used" gorm:"column:is_used"`
	CreatedAt     string `json:"created_at" gorm:"column:created_at"`
	CreatedBy     string `json:"created_by" gorm:"column:created_by"`
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
