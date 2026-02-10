package router

import "github.com/labstack/echo/v4"

func InitUserRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.user

	route.GET("", service.GetAllUser, RequirePermission("gateway.user.read"))
	route.GET("/detail/:id", service.GetUserDetail, RequirePermission("gateway.user.read"))
	route.PUT("", service.UpdateUser, RequirePermission("gateway.user.update"))
	route.DELETE("/:id", service.DeleteUser, RequirePermission("gateway.user.delete"))
	route.GET("/institutions", service.GetInstitutionList, RequirePermission("gateway.user.read"))

	route.POST("/profile-photo", service.UploadProfilePhoto, RequirePermission("gateway.user.upload_profile_photo"))
	route.POST("/cover-photo", service.UploadCoverPhoto, RequirePermission("gateway.user.upload_cover_photo"))

}
