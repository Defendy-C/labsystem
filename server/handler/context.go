package handler

import (
	"labsystem/util/vcode"
)

type SendVerifyCodeKeyResp struct {
	Key string `json:"key"`
}

type GetVImageReq struct {
	Key string `json:"key"`
	Typ int    `json:"typ"`
}

func (req *GetVImageReq) Valid() bool {
	typ := vcode.VImageTyp(req.Typ)
	switch typ {
	case vcode.Puzzle:
		return true
	case vcode.Image:
		return true
	}

	return false
}

type VerifyCodeReq struct {
	Key string `json:"key"`
	Code int `json:"code"`
}

func (req *VerifyCodeReq)Valid() bool  {
	if req.Code <= 0 {
		return false
	}

	return true
}

type GetImageReq struct {
	Url string `json:"url"`
}