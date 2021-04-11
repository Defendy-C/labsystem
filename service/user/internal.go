package user

import (
	"go.uber.org/zap"
	"labsystem/model"
	"labsystem/model/srverr"
	"labsystem/util/logger"
)

type InternalUserSrv interface {
	CreateUser(user *model.User) error
}

var _ InternalUserSrv = &service{}

func (s *service) CreateUser(user *model.User) error {
	if err := s.classSrv.CheckClass(user.Class); err != nil {
		logger.Log.Warn("check class failed", zap.Uint("class", user.Class), zap.Error(err))
		return err
	}
	if err := s.dao.Create(user); err != nil {
		logger.Log.Warn("create user failed", zap.Any("user", user), zap.Error(err))
		return srverr.ErrSystemException
	}

	return nil
}