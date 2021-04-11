package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"labsystem/configs"
	"labsystem/server/handler"
	"labsystem/util/logger"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type ApiRouter struct {
	Api     string
	handles []*handler.Handle
}

type GinServer struct {
	engine           *gin.Engine
	server           *http.Server
	ctx              context.Context
}

func NewGinServer(globalMiddlewares ...gin.HandlerFunc) *GinServer {
	e := gin.New()
	e.Use(globalMiddlewares...)
	config := configs.NewHttpConfig()
	logger.Log.Info("init server config...", zap.Any("config", config))
	s := &http.Server{
		Handler: e,
		Addr:    config.Host + ":" + strconv.Itoa(config.Port),
	}
	ctx := context.Background()

	return &GinServer{server: s, engine: e, ctx: ctx}
}

func (s *GinServer) GinRouterGroup(apiPrefix string, middlewares []gin.HandlerFunc) *gin.RouterGroup {
	rg := s.engine.Group(apiPrefix)
	if middlewares != nil {
		rg.Use(middlewares...)
	}

	return rg
}

func (s *GinServer) Run() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGSEGV)
	logger.Log.Info("server is starting...")
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			panic("server starting failed")
		}
	}()
	for v := range ch {
		switch v {
		case syscall.SIGHUP:
			fallthrough
		case syscall.SIGINT:
			fallthrough
		case syscall.SIGTERM:
			logger.Log.Info("server is quiting...")
			s.server.Close()
			return
		case syscall.SIGKILL:
			fallthrough
		case syscall.SIGQUIT:
			logger.Log.Error("server caught the abnormal signal, server will be force quit!")
			s.server.Close()
			return
		}
	}
}
