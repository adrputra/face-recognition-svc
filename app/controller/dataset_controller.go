package controller

import (
	"context"
	"face-recognition-svc/app/client"
	"face-recognition-svc/app/config"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InterfaceDatasetController interface {
	UploadUserDataset(ctx context.Context, req *model.Dataset) error
	GetDatasetList(ctx context.Context) ([]*model.Dataset, error)
	DeleteDataset(ctx context.Context, username string) error
	TrainModel(ctx context.Context, institutionID string) (*model.ResponseTrainModel, error)
	GetLastTrainModel(ctx context.Context, institutionID string) (string, error)
	GetModelTrainingHistory(ctx context.Context, req *model.FilterModelTraining) ([]*model.ModelTraining, error)
	GetDatasetsByUsername(ctx context.Context, username string) ([]string, error)
}

type DatasetController struct {
	storageClient client.InterfaceStorageClient
	db            *gorm.DB
	userClient    client.InterfaceUserClient
	cfg           *config.Config
	datasetClient client.InterfaceDatasetClient
}

func NewDatasetController(storageClient client.InterfaceStorageClient, db *gorm.DB, userClient client.InterfaceUserClient, cfg *config.Config, datasetClient client.InterfaceDatasetClient) *DatasetController {
	return &DatasetController{
		storageClient: storageClient,
		db:            db,
		userClient:    userClient,
		cfg:           cfg,
		datasetClient: datasetClient,
	}
}

func (c *DatasetController) UploadUserDataset(ctx context.Context, req *model.Dataset) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: UploadUserDataset")
	defer span.Finish()

	span.LogKV("Request", req.Username, req.Dataset)

	user, err := c.userClient.GetUserDetail(ctx, req.Username)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	bucket := fmt.Sprintf("%s/%s", user.InstitutionID, req.Username)
	req.Bucket = bucket
	req.CreatedAt = time.Now().String()

	tx := c.db.Begin()

	dataset, err := c.datasetClient.GetDatasetList(ctx, req.Username)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	if len(dataset) == 0 {
		err = c.storageClient.StoreFileData(ctx, tx, req)
		if err != nil {
			utils.LogEventError(span, err)
			return err
		}
	}

	err = c.storageClient.UploadFile(ctx, req.File, "face-dataset", bucket)
	if err != nil {
		utils.LogEventError(span, err)
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	return nil
}

func (c *DatasetController) GetDatasetList(ctx context.Context) ([]*model.Dataset, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetDatasetList")
	defer span.Finish()

	datasets, err := c.datasetClient.GetDatasetList(ctx, "")
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", datasets)

	return datasets, nil
}

func (c *DatasetController) DeleteDataset(ctx context.Context, username string) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: DeleteDataset")
	defer span.Finish()

	span.LogKV("Request", username)

	user, err := c.userClient.GetUserDetail(ctx, username)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	tx := c.db.Begin()

	err = c.storageClient.DeleteDatasetDB(ctx, tx, username)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	prefix := fmt.Sprintf("%s/%s/", user.InstitutionID, username)

	utils.LogEvent(span, "Request", prefix)

	err = c.storageClient.DeleteObject(ctx, c.cfg.MinioProfile.Bucket, prefix)
	if err != nil {
		utils.LogEventError(span, err)
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Delete Dataset")

	return nil
}

func (c *DatasetController) TrainModel(ctx context.Context, institutionID string) (*model.ResponseTrainModel, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: TrainModel")
	defer span.Finish()

	utils.LogEvent(span, "Request", institutionID)

	session, err := utils.GetMetadata(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "session", session)

	modelReq := &model.ModelTraining{
		ID:            uuid.New().String(),
		InstitutionID: institutionID,
		Status:        "STARTED",
		CreatedAt:     time.Now().String(),
		CreatedBy:     session.Username,
	}

	tx := c.db.Begin()

	err = c.datasetClient.InsertTrainedModel(ctx, modelReq, tx)
	if err != nil {
		utils.LogEventError(span, err)
		tx.Rollback()
		return nil, err
	}

	req := &model.RequestAPITrainModel{
		BucketName: c.cfg.MinioProfile.Bucket,
		Prefix:     institutionID,
		CreatedBy:  session.Username,
		ID:         modelReq.ID,
	}
	res, err := c.datasetClient.TrainModel(ctx, req)
	if err != nil {
		utils.LogEventError(span, err)
		tx.Rollback()
		return nil, err
	}

	utils.LogEvent(span, "Response", res)

	result := &model.ResponseTrainModel{
		ID: res.Data.ID,
	}

	err = tx.Commit().Error
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	return result, nil
}

func (c *DatasetController) GetLastTrainModel(ctx context.Context, institutionID string) (string, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetLastTrainModel")
	defer span.Finish()

	utils.LogEvent(span, "Request", institutionID)

	res, err := c.datasetClient.GetLastTrainModel(ctx, institutionID)
	if err != nil {
		utils.LogEventError(span, err)
		return "", err
	}

	utils.LogEvent(span, "Response", res)

	return res, nil
}

func (c *DatasetController) GetModelTrainingHistory(ctx context.Context, req *model.FilterModelTraining) ([]*model.ModelTraining, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetModelTrainingHistory")
	defer span.Finish()

	utils.LogEvent(span, "Request", req)

	res, err := c.datasetClient.GetModelTrainingHistory(ctx, req)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", res)

	return res, nil
}

func (c *DatasetController) GetDatasetsByUsername(ctx context.Context, username string) ([]string, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetDatasetByUsername")
	defer span.Finish()

	utils.LogEvent(span, "Request", username)

	res, err := c.storageClient.GetDatasetsByUsername(ctx, c.cfg.MinioProfile.Bucket, username)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	utils.LogEvent(span, "Response", res)

	return res, nil
}
