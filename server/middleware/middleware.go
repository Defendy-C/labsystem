package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"labsystem/model"
	"labsystem/model/srverr"
	"labsystem/server/handler"
	"labsystem/util/jwt"
	"labsystem/util/logger"
	"net/http"
)

func ReqLogger(ctx *gin.Context) {
	req, err := ctx.GetRawData()
	if err != nil {
		logger.Log.Warn("don't get request data")
	}
	if len(req) == 0 {
		req = []byte("{}")
	}
	// limit log length
	if len(req) > 1024 {
		logger.Log.Info(ctx.ClientIP() + " send request: "+ string(req[:1024]) + "...... to " + ctx.Request.RequestURI)

	} else {
		logger.Log.Info(ctx.ClientIP() + " send request: "+ string(req) + " to " + ctx.Request.RequestURI)
	}

	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(req))
}

func VerifyToken(ctx *gin.Context) {
	rawToken := ctx.Request.Header.Get("token")
	if rawToken == "" {
		logger.Log.Warn("not found token")
		ctx.JSON(http.StatusUnauthorized, handler.NewResp(srverr.ErrTokenNotFound, nil))
		ctx.Abort()
		return
	}
	payload, err := jwt.ParseToken(rawToken)
	if err != nil {
		logger.Log.Warn("parse token error", zap.String("token", rawToken), zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, handler.NewResp(srverr.ErrInvalidToken, nil))
		ctx.Abort()
		return
	}
	ctx.Set("uid", payload["uid"])
	ctx.Set("rid", payload["rid"])
	if payload["rid"] == model.Visitor {
		ctx.Set("u_rid", payload["u_rid"])
	}
}