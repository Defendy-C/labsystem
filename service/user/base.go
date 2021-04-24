package user

import (
	userDao "labsystem/dao/user"
	"labsystem/model"
	commonSrv "labsystem/service"
	classSrv "labsystem/service/class"
)

type service struct {
	dao      userDao.DaoUser
	baseSrv  commonSrv.Service
	classSrv classSrv.InternalClassSrv
}

type ServiceUser interface {
	RegisterStudent(user *model.User, checker string, key string, vcode int) error
	CheckStudent(operator string, user string) error
	CheckList(operator int) ([]*model.User, error)
	Login(userNo, password, key string, vCode int) (token string, err error)
	Info(uid uint) *model.User
}

type InternalUserSrv interface {
	CreateUser(user *model.User) error
	List(page, pageSize uint) (list []*model.User, totalPage, count uint)
	DeleteUsers(ids []uint) error
}

func NewUserService(baseService commonSrv.Service, classService classSrv.InternalClassSrv) *service {
	return &service{
		dao:      userDao.NewUserDao(),
		classSrv: classService,
		baseSrv:  baseService,
	}
}