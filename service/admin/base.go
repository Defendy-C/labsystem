package admin

import (
	"go.uber.org/zap"
	"labsystem/dao"
	adminDao "labsystem/dao/admin"
	"labsystem/model"
	common "labsystem/service"
	"labsystem/service/user"
	"labsystem/util/logger"
	"time"
)

type ServiceAdmin interface {
	CreateAdmin(admin *model.Admin) error
	CreateTeacher(user *model.User) error
	Login(nickName, password, key string, vcode int) (token string, err error)
	QueryAdminById(id uint) *model.Admin
	QueryAdminByName(name string) *model.Admin
	AdminList(opt *ListOpt, page, pageSize uint) (list []*model.Admin, totalPage, count uint)
	DeleteAdmin(operatorId uint, adminId uint) bool
	UpdatePower(operatorId, adminId uint, add, remove int) bool
	UpdateAdmin(uid uint, nickName string, password string) bool
}

type OrderField string
const (
	CREATEDAT = "created_at"
)
func (f OrderField)ToString() string {
	return string(f)
}

type ListOpt struct {
	CreatedBy  string
	CreatedMin *time.Time
	CreatedMax *time.Time
	OrderFiled OrderField
	OrderDesc  bool
}
func (opt *ListOpt) OrderPad() (string, dao.OrderTyp) {
	var typ dao.OrderTyp
	if opt.OrderDesc {
		typ = dao.DESC
	} else {
		typ = dao.ASC
	}

	return opt.OrderFiled.ToString(), typ
}

type service struct {
	dao adminDao.DaoAdmin
	userSrv user.InternalUserSrv
	commonSrv common.Service
}

// check this operator of admin is existing
func (s *service) checkAdmin(operatorId, adminId uint) (operator, admin *model.Admin) {
	list, err := s.dao.Query(&adminDao.FilterAdmin{
		ID: []uint{operatorId, adminId},
	})
	if err != nil {
		logger.Log.Warn("query operator,admin failed", zap.Error(err))
		return
	}
	admins, ok := list.([]*model.Admin)
	if !ok || len(admins) != 2 {
		logger.Log.Warn("query operator,admin failed", zap.Any("list", list))
		return
	}
	for _, v := range admins {
		switch v.ID {
		case operatorId:
			operator = v
		case adminId:
			admin = v
		default:
			logger.Log.Warn("query operator,admin failed")
			return nil, nil
		}
	}

	if admin == nil || operator == nil || admin.CreatedBy != operator.NickName {
		return nil, nil
	}

	return
}