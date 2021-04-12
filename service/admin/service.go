package admin

import (
	"go.uber.org/zap"
	"labsystem/dao"
	adminDao "labsystem/dao/admin"
	"labsystem/model"
	"labsystem/model/srverr"
	commonSrv "labsystem/service"
	userSrv "labsystem/service/user"
	"labsystem/util/jwt"
	"labsystem/util/logger"
	"labsystem/util/rsa"
)

var _ ServiceAdmin = &service{}

func NewAdminService(userService userSrv.InternalUserSrv, commonService commonSrv.Service) *service {
	return &service{
		dao: adminDao.NewAdminDao(),
		userSrv: userService,
		commonSrv: commonService,
	}
}

func (s *service) CreateAdmin(admin *model.Admin) error {
	// verify createdBy
	creator := s.QueryAdminByName(admin.CreatedBy)
	if creator == nil {
		return srverr.ErrInvalidCreator
	}
	// verify power
	if !creator.Power.Own(model.PowerCreateAdmin) || !creator.Power.Own(admin.Power) {
		return srverr.ErrOwnPower
	}
	return s.dao.Create(admin)
}


func (s *service) QueryAdminById(id uint) *model.Admin {
	data, err := s.dao.Query(&adminDao.FilterAdmin{
		ID: []uint{id},
	})
	if err != nil {
		logger.Log.Warn("query admin failed", zap.Error(err))
		return nil
	}
	admins := data.([]*model.Admin)
	if len(admins) > 0 {
		return admins[0]
	}

	return nil
}

func (s *service) QueryAdminByName(name string) *model.Admin {
	data, err := s.dao.Query(&adminDao.FilterAdmin{
		NickName: []string{name},
	})
	if err != nil {
		logger.Log.Warn("query admin failed", zap.Error(err))
		return nil
	}
	admins := data.([]*model.Admin)
	if len(admins) > 0 {
		return admins[0]
	}

	return nil
}

func (s *service) Login(nickName, password, key string, vcode int) (token string, err error) {
	// verify vcode
	if !s.commonSrv.VerifyCode(key, vcode) {
		return "", srverr.ErrVerify
	}
	if err := s.dao.CDelete(key); err != nil {
		logger.Log.Warn("cache key delete failed", zap.String("key", key), zap.Error(err))
	}
	// verify nickname, password
	filter := &adminDao.FilterAdmin{
		NickName: []string{nickName},
	}
	obj, err := s.dao.Query(filter)
	com := obj.([]*model.Admin)
	if err != nil || len(com) <= 0 {
		logger.Log.Warn("query admin failed", zap.Any("filter", filter), zap.Error(err))
		return "", srverr.ErrLoginFailed
	}
	if ok := rsa.Compare(com[0].Password, password); !ok {
		logger.Log.Warn("password error", zap.String("expected", com[0].Password), zap.String("actual", password))
		return "", srverr.ErrLoginFailed
	}
	//generate token
	m := map[string]interface{}{
		"uid": com[0].ID,
		"rid": model.Administrator,
	}
	token, err = jwt.Token(m)
	if err != nil {
		logger.Log.Warn("token generating failed", zap.Any("map", m), zap.Error(err))
		return "", srverr.ErrSystemException
	}

	return
}

func (s *service) CreateTeacher(user *model.User) error {
	user.Status = model.Teacher
	return s.userSrv.CreateUser(user)
}

func (s *service) AdminList(opt *ListOpt, page, pageSize uint) (list []*model.Admin, totalPage, count uint) {
	filter := &adminDao.FilterAdmin{
		BaseFilter: dao.BaseFilter {
			Page: page,
			PerPage: pageSize,
		},
	}
	if opt != nil {
		if opt.CreatedBy != "" {
			filter.CreatedBy = []string{opt.CreatedBy}
		}
		filter.SetCreatedAtRange(opt.CreatedMin, opt.CreatedMax)
		if opt.OrderFiled.ToString() != "" {
			filter.Sort = dao.NewOrderBy(opt.OrderPad())
		}
	}
	if obj, err := s.dao.Query(filter); err != nil {
		logger.Log.Warn("admin list failed", zap.Any("filter", filter), zap.Error(err))
	} else {
		list = obj.([]*model.Admin)
		totalPage = filter.TotalPage
		count = filter.TotalCount
	}

	return
}

func (s *service) DeleteAdmin(operatorId uint, adminId uint) bool {
	 if operator, admin := s.checkAdmin(operatorId, adminId);operator == nil || admin == nil {
		 logger.Log.Warn("operator is irrelevant to the admin", zap.Any("operator", operator), zap.Any("admin", admin))
	 	return false
	 }
	if err := s.dao.Delete(map[string]interface{}{
		"id": adminId,
	}); err != nil {
		logger.Log.Warn("delete admin failed", zap.Error(err))
		return false
	}

	return true
}

func (s *service) UpdatePower(operatorId, adminId uint, add, remove int) bool {
	// check id and relation
	var operator, admin *model.Admin
	if operator, admin = s.checkAdmin(operatorId, adminId); admin == nil || operator == nil {
		logger.Log.Warn("operator is irrelevant to the admin", zap.Any("operator", operator), zap.Any("admin", admin))
		return false
	}
	// check power
	removedPower, err := model.IntToPower(remove)
	if err != nil {
		logger.Log.Warn("invalid power", zap.Any("power", remove))
		return false
	}
	addPower, err := model.IntToPower(add)
	if err != nil {
		logger.Log.Warn("invalid power", zap.Any("power", add))
		return false
	}
	if !admin.Power.Own(removedPower) {
		logger.Log.Warn("invalid power", zap.Any("power", remove))
		return false
	}
	if !admin.Power.Disown(addPower) || !operator.Power.Own(addPower) {
		logger.Log.Warn("invalid power", zap.Any("power", add))
		return false
	}
	// remove power
	admin.Power = admin.Power.Sub(removedPower)
	// add power
	admin.Power = admin.Power.Add(addPower)
	// update
	if err := s.dao.Update(map[string]interface{}{
		"id": adminId,
	}, map[string]interface{}{
		"power": admin.Power,
	}); err != nil {
		logger.Log.Warn("update power error", zap.Error(err), zap.Any("power", admin.Power))
		return false
	}

	return true
}

func (s *service) UpdateAdmin(uid uint, nickName string, password string) bool {
	if err := s.dao.Update(map[string]interface{}{
		"id": uid,
	}, map[string]interface{}{
		"nick_name": nickName,
		"password": password,
	});err != nil {
		logger.Log.Warn("update admin failed", zap.Error(err))
		return false
	}

	return true
}