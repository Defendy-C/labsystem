package user

import (
	"go.uber.org/zap"
	"labsystem/dao"
	"labsystem/dao/user"
	"labsystem/model"
	"labsystem/model/srverr"
	"labsystem/util/logger"
)

var _ InternalUserSrv = &service{}

func (s *service) CreateUser(user *model.User) error {
	if err := s.classSrv.CheckClass(user.Class); err != nil {
		logger.Log.Warn("check class failed", zap.String("class", user.Class), zap.Error(err))
		return srverr.ErrInvalidClass
	}
	user.Status = model.Teacher
	if err := s.dao.Create(user); err != nil {
		logger.Log.Warn("create user failed", zap.Any("user", user), zap.Error(err))
		return srverr.ErrSystemException
	}

	return nil
}

func (s *service) DeleteUsers(ids []uint) error {
	return s.dao.Delete(map[string]interface{}{
		"id": ids,
	})
}

func (s *service) List(page, pageSize uint) (list []*model.User, totalPage, count uint) {
	filter := &dao.BaseFilter{
		Page: page,
		PerPage: pageSize,
	}
	raw, err := s.dao.Query(&user.FilterUser{
		BaseFilter: filter,
	})
	if err != nil {
		logger.Log.Warn("query user failed", zap.Error(err))
		return nil, 0, 0
	}
	totalPage = filter.TotalPage
	count = filter.TotalCount
	list = raw.([]*model.User)

	return
}

