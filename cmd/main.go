package main

import (
	"github.com/gin-gonic/gin"
	ginServer "labsystem/server"
	"labsystem/server/handler"
	adminHandler "labsystem/server/handler/admin"
	"labsystem/server/middleware"
	"labsystem/service"
	adminSrv "labsystem/service/admin"
	classSrvPkg "labsystem/service/class"
	userSrv "labsystem/service/user"
)

const (
	COMMONAPI = "/api/common/"
	ADMINAPI = "/api/admin/"
)

func main() {
	// service obj
	service := service.NewService()
	classService := classSrvPkg.NewClassService()
	userService := userSrv.NewUserService(classService)
	adminService := adminSrv.NewAdminService(userService, service)
	// handles
	common := handler.CommonHandler{Srv: service}
	admin := adminHandler.HandlerAdmin{Srv: adminService}
	// register service
	server := ginServer.NewGinServer(middleware.ReqLogger)
	authRgMw := []gin.HandlerFunc{middleware.VerifyToken}
	common.RegisterCommonHandles(server.GinRouterGroup(COMMONAPI, nil))
	admin.RegisterAdminHandles(server.GinRouterGroup(ADMINAPI, nil), server.GinRouterGroup(ADMINAPI, authRgMw))
	server.Run()
}
