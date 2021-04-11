package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"labsystem/util/jwt"
	"labsystem/util/logger"
)

func ReqLogger(ctx *gin.Context) {
	req, err := ctx.GetRawData()
	if err != nil {
		logger.Log.Warn("don't get request data")
	}
	if len(req) == 0 {
		req = []byte("{}")
	}
	logger.Log.Info(ctx.ClientIP() + " send request: "+ string(req) + " to " + ctx.Request.RequestURI)

	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(req))
}

func VerifyToken(ctx *gin.Context) {
	rawToken := ctx.Request.Header.Get("token")
	if rawToken == "" {
		logger.Log.Warn("not found token")
		ctx.Abort()
		return
	}
	payload, err := jwt.ParseToken(rawToken)
	if err != nil {
		logger.Log.Warn("parse token error", zap.String("token", rawToken), zap.Error(err))
		ctx.Abort()
		return
	}
	ctx.Set("uid", payload["uid"])
	ctx.Set("rid", payload["rid"])
}