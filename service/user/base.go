package user

import (
	userDao "labsystem/dao/user"
	classSrv "labsystem/service/class"
)

type service struct {
	dao userDao.DaoUser
	classSrv classSrv.InternalClassSrv
}

func NewUserService(classService classSrv.InternalClassSrv) *service {
	return &service{
		dao: userDao.NewUserDao(),
		classSrv: classService,
	}
}