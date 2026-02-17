package router

import "github.com/gin-gonic/gin"

func InitUserRoute(prefix string, e *gin.RouterGroup) {
	route := e.Group(prefix)
	service := factory.Service.user

	route.GET("", RequirePermission("gateway.user.read"), wrapHandler(service.GetAllUser))
	route.GET("/detail/:id", RequirePermission("gateway.user.read"), wrapHandler(service.GetUserDetail))
	route.PUT("", RequirePermission("gateway.user.update"), wrapHandler(service.UpdateUser))
	route.DELETE("/:id", RequirePermission("gateway.user.delete"), wrapHandler(service.DeleteUser))
	route.GET("/institutions", RequirePermission("gateway.user.read"), wrapHandler(service.GetInstitutionList))

	route.POST("/profile-photo", RequirePermission("gateway.user.upload_profile_photo"), wrapHandler(service.UploadProfilePhoto))
	route.POST("/cover-photo", RequirePermission("gateway.user.upload_cover_photo"), wrapHandler(service.UploadCoverPhoto))

}
