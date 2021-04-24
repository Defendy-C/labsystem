package main

import (
	"github.com/gin-gonic/gin"
	ginServer "labsystem/server"
	"labsystem/server/handler"
	adminHandler "labsystem/server/handler/admin"
	userHandler "labsystem/server/handler/user"
	"labsystem/server/middleware"
	"labsystem/service"
	adminSrv "labsystem/service/admin"
	classSrvPkg "labsystem/service/class"
	userSrv "labsystem/service/user"
)

const (
	COMMONAPI = "/api/common/"
	ADMINAPI = "/api/admin/"
	USERAPI = "/api/user/"
)

func main() {
	// srv obj
	srv := service.NewService()
	classService := classSrvPkg.NewClassService()
	userService := userSrv.NewUserService(srv, classService)
	adminService := adminSrv.NewAdminService(userService, classService, srv)
	// handles
	common := handler.CommonHandler{Srv: srv}
	admin := adminHandler.HandlerAdmin{Srv: adminService}
	user := userHandler.HandlerUser{Srv: userService}
	// register srv
	server := ginServer.NewGinServer(middleware.ReqLogger)
	authRgMw := []gin.HandlerFunc{middleware.VerifyToken}
	common.RegisterCommonHandles(server.GinRouterGroup(COMMONAPI, nil), server.GinRouterGroup(COMMONAPI, authRgMw))
	admin.RegisterAdminHandles(server.GinRouterGroup(ADMINAPI, nil), server.GinRouterGroup(ADMINAPI, authRgMw))
	user.RegisterUserHandles(server.GinRouterGroup(USERAPI, nil), server.GinRouterGroup(USERAPI, authRgMw))

	server.Run()
}
