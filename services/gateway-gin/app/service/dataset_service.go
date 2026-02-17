package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/adrputra/face-recognition-svc/gateway-gin/app/controller"
	"github.com/adrputra/face-recognition-svc/gateway-gin/app/model"
	"github.com/adrputra/face-recognition-svc/gateway-gin/app/utils"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InterfaceDatasetService interface {
	UploadUserDataset(e *gin.Context) error
	GetDatasetList(e *gin.Context) error
	DeleteDataset(e *gin.Context) error
	TrainModel(e *gin.Context) error
	GetLastTrainModel(e *gin.Context) error
	GetModelTrainingHistory(e *gin.Context) error
	GetDatasetsByUsername(e *gin.Context) error
}

type DatasetService struct {
	uc controller.InterfaceDatasetController
}

func NewDatasetService(uc controller.InterfaceDatasetController) InterfaceDatasetService {
	return &DatasetService{
		uc: uc,
	}
}

func (s *DatasetService) UploadUserDataset(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "UploadUserDataset")
	defer span.Finish()

	form, err := e.MultipartForm()
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	files := form.File["file"]

	utils.LogEvent(span, "Request", "")

	var attach []*model.File
	for _, file := range files {
		// Open the file
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		var buffer bytes.Buffer
		_, err = io.Copy(&buffer, src)
		if err != nil {
			utils.LogEventError(span, err)
			return utils.LogError(e, err, nil)
		}
		imageBytes := buffer.Bytes()
		attach = append(attach, &model.File{
			FileName:    file.Filename,
			BytesObject: imageBytes,
		})
	}

	request := &model.Dataset{
		Username: e.PostForm("username"),
		File:     attach,
	}

	err = s.uc.UploadUserDataset(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", "Upload Success")

	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Upload Success",
		Data:    nil,
	})
}

func (s *DatasetService) GetDatasetList(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "GetDatasetList")
	defer span.Finish()

	dataset, err := s.uc.GetDatasetList(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", dataset)

	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get Dataset List",
		Data:    dataset,
	})
}

func (s *DatasetService) DeleteDataset(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "DeleteDataset")
	defer span.Finish()

	id := e.Param("id")

	utils.LogEvent(span, "Request", id)

	if id == "" {
		utils.LogEventError(span, errors.New("id shouldn't be empty"))
		return utils.LogError(e, errors.New("id shouldn't be empty"), nil)
	}

	err := s.uc.DeleteDataset(ctx, id)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", "Delete Success")

	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Delete Success",
		Data:    nil,
	})
}

func (s *DatasetService) TrainModel(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "TrainModel")
	defer span.Finish()

	institutionID := e.Param("id")

	utils.LogEvent(span, "Request", institutionID)

	res, err := s.uc.TrainModel(ctx, institutionID)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", res)

	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Train Model",
		Data:    res,
	})
}

func (s *DatasetService) GetLastTrainModel(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "GetLastTrainModel")
	defer span.Finish()

	id := e.Param("id")

	utils.LogEvent(span, "Request", id)

	if id == "" {
		utils.LogEventError(span, errors.New("id shouldn't be empty"))
		return utils.LogError(e, errors.New("id shouldn't be empty"), nil)
	}

	res, err := s.uc.GetLastTrainModel(ctx, id)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", res)

	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get Last Train Model",
		Data:    res,
	})
}

func (s *DatasetService) GetModelTrainingHistory(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "GetModelTrainingHistory")
	defer span.Finish()

	var request model.FilterModelTraining

	if err := e.ShouldBind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Request", request)

	res, err := s.uc.GetModelTrainingHistory(ctx, &request)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", res)

	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get Model Training History",
		Data:    res,
	})
}

func (s *DatasetService) GetDatasetsByUsername(e *gin.Context) error {
	ctx, span := utils.StartSpan(e, "GetDatasetsByUsername")
	defer span.Finish()

	id := e.Param("id")
	institutionID := e.Param("institution-id")

	utils.LogEvent(span, "institutionID", institutionID)
	utils.LogEvent(span, "id", id)

	if institutionID == "" {
		utils.LogEventError(span, errors.New("institutionID shouldn't be empty"))
		return utils.LogError(e, errors.New("institutionID shouldn't be empty"), nil)
	}

	if id == "" {
		utils.LogEventError(span, errors.New("id shouldn't be empty"))
		return utils.LogError(e, errors.New("id shouldn't be empty"), nil)
	}

	res, err := s.uc.GetDatasetsByUsername(ctx, fmt.Sprintf("%s/%s", institutionID, id))
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Response", res)

	return utils.JSON(e, http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get Datasets By Username",
		Data:    res,
	})
}
