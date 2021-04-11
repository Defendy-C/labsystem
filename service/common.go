package service

import (
	"go.uber.org/zap"
	"io"
	"labsystem/dao"
	"labsystem/model"
	"labsystem/util"
	"labsystem/util/logger"
	"labsystem/util/vcode"
	"os"
	"strconv"
	"time"
)

type Service interface {
	GenerateVCode() (key string)
	GetVImage(key string, typ vcode.VImageTyp) (outStream io.ReadCloser)
	VerifyCode(key string, vCode int) bool
}

type service struct {
	dao *dao.DAO
}

func NewService() Service {
	return &service{dao.NewDAO()}
}

var _ Service = &service{}

func (s *service)GenerateVCode() (key string) {
	key = strconv.Itoa(int(time.Now().UnixNano())) + strconv.Itoa(util.RandIntN(0))
	for i := 0;s.dao.CGet(key) != "" && i < model.MAXTryCount;i++ {
		key = strconv.Itoa(int(time.Now().UnixNano())) + strconv.Itoa(util.RandIntN(0))
		if i == model.MAXTryCount - 1 {
			logger.Log.Error("don't generate vcode")
			return ""
		}
	}
	no := util.RandIntN(4) + 1
	vc := vcode.GetVCode()

	if err := s.dao.CSet(key, strconv.Itoa(no) + strconv.Itoa(vc), time.Minute); err != nil {
		logger.Log.Error("verify code don't save in cache", zap.String("key", key))
	}
	return
}

func (s *service) GetVImage(key string, typ vcode.VImageTyp) (outStream io.ReadCloser) {
	fileTyp := ".jpeg"
	vc := s.dao.CGet(key)
	if vc == "" {
		logger.Log.Warn("no found key in the cache", zap.String("key", key))
		return nil
	}

	var f *os.File
	var err error
	switch typ {
	case vcode.Puzzle:
		f, err = os.Open( vcode.VImagePath + vc + vcode.PuzzleSuffix + fileTyp)
		if err != nil {
			logger.Log.Warn("invalid verify code", zap.String("vc", vc), zap.Error(err))
			return nil
		}
	case vcode.Image:
		f, err = os.Open( vcode.VImagePath + vc + fileTyp)
		if err != nil {
			logger.Log.Warn("invalid verify code", zap.String("vc", vc), zap.Error(err))
			return nil
		}
	}

	return f
}

func (s *service)VerifyCode(key string, vCode int) bool {
	vc := s.dao.CGet(key)
	if vc == "" {
		return false
	}
	vcCom, err := strconv.Atoi(vc[1:])
	if err != nil {
		logger.Log.Warn("convert string to int error:", zap.String("vCode", vc), zap.Error(err))
		return false
	}
	if val := vCode - vcCom; val < -5 || val > 5 {
		return false
	}

	return true
}