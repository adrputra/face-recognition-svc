package service

import (
	"net/http"
	"strings"

	"github.com/adrputra/face-recognition-svc/gateway/app/controller"
	"github.com/adrputra/face-recognition-svc/gateway/app/domain"
	"github.com/adrputra/face-recognition-svc/gateway/app/utils"
	"github.com/labstack/echo/v4"
)

type InterfacePresenceService interface {
	CreatePresence(e echo.Context) error
	GetPresence(e echo.Context) error
	ListPresence(e echo.Context) error
	UpdatePresence(e echo.Context) error
	DeletePresence(e echo.Context) error
}

type PresenceService struct {
	presenceController controller.InterfacePresenceController
}

func NewPresenceService(presenceController controller.InterfacePresenceController) InterfacePresenceService {
	return &PresenceService{
		presenceController: presenceController,
	}
}

func (s *PresenceService) CreatePresence(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "CreatePresence")
	defer span.Finish()

	var request domain.Presence
	if err := e.Bind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	res, err := s.presenceController.CreatePresence(ctx, &request)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, domain.Response{
		Code:    200,
		Message: "Success Create Presence",
		Data:    res,
	})
}

func (s *PresenceService) GetPresence(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetPresence")
	defer span.Finish()

	id := strings.TrimSpace(e.Param("id"))
	res, err := s.presenceController.GetPresence(ctx, id)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, domain.Response{
		Code:    200,
		Message: "Success Get Presence",
		Data:    res,
	})
}

func (s *PresenceService) ListPresence(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "ListPresence")
	defer span.Finish()

	filter := domain.PresenceFilter{
		UserID:        strings.TrimSpace(e.QueryParam("user_id")),
		InstitutionID: strings.TrimSpace(e.QueryParam("institution_id")),
		Status:        strings.TrimSpace(e.QueryParam("status")),
	}

	res, err := s.presenceController.ListPresence(ctx, filter)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, domain.Response{
		Code:    200,
		Message: "Success Get Presence List",
		Data:    res,
	})
}

func (s *PresenceService) UpdatePresence(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "UpdatePresence")
	defer span.Finish()

	id := strings.TrimSpace(e.Param("id"))
	var request domain.Presence
	if err := e.Bind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	res, err := s.presenceController.UpdatePresence(ctx, id, &request)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, domain.Response{
		Code:    200,
		Message: "Success Update Presence",
		Data:    res,
	})
}

func (s *PresenceService) DeletePresence(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "DeletePresence")
	defer span.Finish()

	id := strings.TrimSpace(e.Param("id"))
	if err := s.presenceController.DeletePresence(ctx, id); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, domain.Response{
		Code:    200,
		Message: "Success Delete Presence",
		Data:    nil,
	})
}
