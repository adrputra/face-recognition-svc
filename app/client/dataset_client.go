package client

import (
	"context"
	"encoding/json"
	"face-recognition-svc/app/config"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"fmt"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"gorm.io/gorm"
)

type InterfaceDatasetClient interface {
	GetDatasetList(ctx context.Context, user string) ([]*model.Dataset, error)
	TrainModel(ctx context.Context, request *model.RequestAPITrainModel) (res *model.ResponseAPITrainModel, err error)
	GetLastTrainModel(ctx context.Context, institutionID string) (string, error)
	GetModelTrainingHistory(ctx context.Context, req *model.FilterModelTraining) ([]*model.ModelTraining, error)
	InsertTrainedModel(ctx context.Context, req *model.ModelTraining, tx *gorm.DB) error
}

type DatasetClient struct {
	db  *gorm.DB
	cfg *config.Config
	mq  *amqp.Channel
}

func NewDatasetClient(db *gorm.DB, cfg *config.Config, mq *amqp.Channel) *DatasetClient {
	return &DatasetClient{
		db:  db,
		cfg: cfg,
		mq:  mq,
	}
}

func (c *DatasetClient) GetDatasetList(ctx context.Context, user string) ([]*model.Dataset, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetDatasetList")
	defer span.Finish()

	var result []*model.Dataset

	var sb strings.Builder
	sb.WriteString("SELECT username, dataset, created_at FROM face_datasets")

	if user != "" {
		sb.WriteString(fmt.Sprintf(" WHERE username = '%s'", user))
	}

	query := sb.String()
	err := c.db.Debug().WithContext(ctx).Raw(query).Scan(&result).Error

	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", result)

	return result, nil
}

func (d *DatasetClient) TrainModel(ctx context.Context, request *model.RequestAPITrainModel) (res *model.ResponseAPITrainModel, err error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: Processing TrainModel")
	defer span.Finish()

	utils.LogEvent(span, "Request", request)

	// Declare a queue
	_, err = d.mq.QueueDeclare(
		"TrainModel", // Queue name
		true,         // Durable
		true,         // Delete when unused
		false,        // Exclusive
		false,        // No-wait
		nil,          // Arguments
	)

	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	messageJSON, err := json.Marshal(request)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	err = d.mq.Publish(
		"",           // Exchange (default)
		"TrainModel", // Routing key (queue name)
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         messageJSON,
			DeliveryMode: amqp.Persistent, // Make message persistent
		},
	)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	// endpoint := fmt.Sprintf("%s:%d%s", d.cfg.API.ProcessingSVC.Host, d.cfg.API.ProcessingSVC.Port, d.cfg.API.ProcessingSVC.Endpoint)

	out := &model.ResponseAPITrainModel{}
	// err = utils.RequestAPI("POST", endpoint, request, &out)
	// if err != nil {
	// 	utils.LogEventError(span, err)
	// 	return nil, err
	// }

	// utils.LogEvent(span, "Response API", out)

	// if out.Code != 200 {
	// 	utils.LogEventError(span, errors.New(out.Message))
	// 	return nil, errors.New(out.Message)
	// }

	return out, nil
}

func (d *DatasetClient) GetLastTrainModel(ctx context.Context, institutionID string) (string, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetLastTrainModel")
	defer span.Finish()

	utils.LogEvent(span, "Request", institutionID)

	var res string

	query := "SELECT created_at FROM model_training WHERE institution_id = ? ORDER BY created_at DESC LIMIT 1"

	err := d.db.Debug().Raw(query, institutionID).Scan(&res).Error
	if err != nil {
		utils.LogEventError(span, err)
		return "", err
	}

	utils.LogEvent(span, "Response", res)

	return res, nil
}

func (d *DatasetClient) GetModelTrainingHistory(ctx context.Context, req *model.FilterModelTraining) ([]*model.ModelTraining, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetModelTrainingHistory")
	defer span.Finish()

	utils.LogEvent(span, "Request", req)

	var res []*model.ModelTraining

	sb := strings.Builder{}
	if req.InstitutionID != "" {
		sb.WriteString(fmt.Sprintf(" WHERE institution_id = '%s'", req.InstitutionID))
		if req.Status != "" {
			sb.WriteString(fmt.Sprintf(" AND status = '%s'", req.Status))
		}
		if req.IsUsed != "" {
			sb.WriteString(fmt.Sprintf(" AND is_used = '%s'", req.IsUsed))
		}
	}

	if req.OrderBy != "" {
		sb.WriteString(fmt.Sprintf(" ORDER BY %s %s", req.OrderBy, req.SortType))
	} else {
		sb.WriteString(" ORDER BY created_at DESC")
	}

	query := "SELECT * FROM model_training"

	err := d.db.Debug().Raw(query + sb.String()).Scan(&res).Error
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", res)

	return res, nil
}

func (d *DatasetClient) InsertTrainedModel(ctx context.Context, req *model.ModelTraining, tx *gorm.DB) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: InsertTrainedModel")
	defer span.Finish()

	var args []interface{}

	args = append(args, req.ID, req.InstitutionID, req.Status, time.Now(), req.CreatedBy)
	query := "INSERT INTO model_training (id, institution_id, status, created_at, created_by) VALUES (?, ?, ?, ?, ?)"
	result := tx.Debug().Exec(query, args...)

	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return result.Error
	}

	return nil
}
