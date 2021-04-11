package class

import (
	"labsystem/model"
	"labsystem/model/srverr"
)

type ServiceClass interface {
	CreateClass(class *model.Class) error
}

var _ ServiceClass = &service{}

func (s service) CreateClass(class *model.Class) error {
	if err := s.dao.Create(class); err != nil {
		return srverr.ErrSystemException
	}

	return nil
}
