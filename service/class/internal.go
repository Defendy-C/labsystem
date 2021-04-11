package class

import (
	"go.uber.org/zap"
	classDao "labsystem/dao/class"
	"labsystem/model/srverr"
	"labsystem/util/logger"
)

type InternalClassSrv interface {
	CheckClass(id uint) error
}

var _ InternalClassSrv = &service{}

func (s service) CheckClass(id uint) error {
	filter := classDao.FilterClass{
		Id: []uint{id},
	}

	class, err := s.dao.Query(&filter)
	if err != nil || class == nil {
		logger.Log.Warn("found not record", zap.Any("filter", filter), zap.Error(err))
		return srverr.ErrInvalidClass
	}

	return nil
}

