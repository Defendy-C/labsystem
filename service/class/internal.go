package class

import (
	"go.uber.org/zap"
	"labsystem/dao"
	classDao "labsystem/dao/class"
	"labsystem/model"
	"labsystem/model/srverr"
	"labsystem/util/logger"
)

var _ InternalClassSrv = &service{}

func (s service) CheckClass(name string) error {
	filter := classDao.FilterClass{
		ClassNo: []string{name},
	}

	class, err := s.dao.Query(&filter)
	if err != nil || len(class.([]*model.Class)) == 0 {
		logger.Log.Warn("found not record", zap.Any("filter", filter), zap.Error(err))
		return srverr.ErrInvalidClass
	}

	return nil
}

func (s *service) CreateClass(class *model.Class) error {
	if err := s.dao.Create(class); err != nil {
		return srverr.ErrSystemException
	}

	return nil
}

func (s *service) List(opt *ListOpt, page, pageSize uint) (list []*model.Class, totalPage, totalCount uint) {
	filter := &classDao.FilterClass{}
	filter.BaseFilter = &dao.BaseFilter{
		Page: page,
		PerPage: pageSize,
	}
	if opt != nil {
		filter.Sort = dao.NewOrderBy(opt.OrderPad())
	}

	raw, err := s.dao.Query(filter)
	if err != nil {
		logger.Log.Warn("query class list failed", zap.Error(err))
		return nil, 0, 0
	}

	return raw.([]*model.Class), filter.TotalPage, filter.TotalCount
}
