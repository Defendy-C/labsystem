package class

import (
	classDao "labsystem/dao/class"
)

type service struct {
	dao classDao.DaoClass
}

func NewClassService() *service {
	return &service{
		dao: classDao.NewClassDao(),
	}
}


