package handler

import (
	"github.com/gin-gonic/gin"
	"io"
	"labsystem/model"
	"labsystem/model/srverr"
	commonSrv "labsystem/service"
	"labsystem/util/vcode"
	"net/http"
)


type Resp struct {
	Success bool   `json:"success"`
	Err  string `json:"err"`
	Data    interface{} `json:"data"`
}

func NewResp(errMsg error, data interface{}) *Resp {
	if errMsg != nil {
		return &Resp{Err: errMsg.Error()}
	}

	return &Resp{Success: true, Data: data}
}

type Handle struct {
	RelativePath string
	HandlerFunc  gin.HandlerFunc
}

func Transfer(reader io.ReadCloser, ctx *gin.Context, speed int) {
	buf := make([]byte, speed)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			ctx.JSON(http.StatusBadRequest, NewResp(srverr.ErrDownload, nil))
			return
		}

		ctx.Data(200, "application/octet-stream", buf[:n])
	}
	reader.Close()
}

type CommonHandler struct {
	Srv     commonSrv.Service
}

func (h *CommonHandler) RegisterCommonHandles(rg *gin.RouterGroup, authRg *gin.RouterGroup) {
	// the RouterGroup rg mustn't be Authenticate
	{
		rg.POST("/vcode/key", h.sendVerifyCodeKey)
		rg.POST("/vcode/resource", h.getVImage)
		rg.POST("/vcode/verify", h.verifyCode)
	}
	// the RouterGroup hRg must be Authenticate
	{
		authRg.POST("/image", h.getImage)
	}
}

func (h *CommonHandler) sendVerifyCodeKey(ctx *gin.Context) {
	key := h.Srv.GenerateVCode()
	if key == "" {
		ctx.JSON(http.StatusBadRequest, NewResp(srverr.ErrSystemException, nil))
		ctx.Abort()
		return
	}
	resp := &SendVerifyCodeKeyResp{
		Key: key,
	}
	ctx.JSON(http.StatusOK, NewResp(nil, resp))
}

func (h *CommonHandler) getVImage(ctx *gin.Context) {
	var req GetVImageReq
	if err := ctx.ShouldBindJSON(&req); err != nil || !req.Valid() {
		ctx.JSON(http.StatusBadRequest, NewResp(srverr.ErrInvalidParams, nil))
		ctx.Abort()
		return
	}
	stream := h.Srv.GetVImage(req.Key, vcode.VImageTyp(req.Typ))
	if stream == nil {
		ctx.JSON(http.StatusBadRequest, NewResp(srverr.ErrSystemException, nil))
		return
	}

	Transfer(stream, ctx, 10 * model.KB)
}

func (h *CommonHandler) verifyCode(ctx *gin.Context) {
	var req VerifyCodeReq
	if err := ctx.ShouldBindJSON(&req); err != nil || !req.Valid() {
		ctx.JSON(http.StatusBadRequest, NewResp(srverr.ErrInvalidParams, nil))
		return
	}
	if !h.Srv.VerifyCode(req.Key, req.Code) {
		ctx.JSON(http.StatusBadRequest, NewResp(srverr.ErrVerify, nil))
		return
	}

	ctx.JSON(http.StatusOK, NewResp(nil, nil))
}

func (h *CommonHandler) getImage(ctx *gin.Context) {
	var req *GetImageReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, NewResp(srverr.ErrInvalidParams, nil))
		return
	}
	stream := h.Srv.GetImage(req.Url)
	if stream == nil {
		ctx.JSON(http.StatusBadRequest, NewResp(srverr.ErrInvalidParams, nil))
		return
	}

	Transfer(stream, ctx, 10 * model.KB)
}