package user

import (
	"encoding/json"
	"go.uber.org/zap"
	userDao "labsystem/dao/user"
	"labsystem/model"
	"labsystem/model/srverr"
	"labsystem/util/jwt"
	"labsystem/util/logger"
	"labsystem/util/rsa"
	"strconv"
)

var _ ServiceUser = &service{}

func (s service) CheckStudent(operator string, user string) error {
	key := "REGISTER_CHECK-" + operator
	j, err := s.dao.CHashGetV(key, user)
	if j == "" || err != nil {
		logger.Log.Warn("found not cache record", zap.String("key", key), zap.Error(err))
		return srverr.ErrSystemException
	}
	var u model.User
	err = json.Unmarshal([]byte(j), &u)
	if err != nil {
		logger.Log.Warn("json parse error", zap.String("json", j), zap.Error(err))
		return srverr.ErrSystemException
	}
	err = s.dao.CHashDelV(operator, user)
	if err != nil {
		return err
	}

	return s.dao.Create(&u)
}

func (s service) CheckList(operator int) (us []*model.User, err error) {
	var m map[string]string
	m, err = s.dao.CHashList("REGISTER_CHECK-" + strconv.Itoa(operator))
	if err != nil {
		logger.Log.Warn("query cache failed", zap.Int("key", operator), zap.Error(err))
		return nil, err
	}
	var index int
	us = make([]*model.User, len(m))
	for _, v := range m {
		us[index] = new(model.User)
		err = json.Unmarshal([]byte(v), &us[index])
		if err != nil {
			logger.Log.Warn("json parse failed", zap.String("json", v), zap.Error(err))
			return
		}
	}

	return
}

func (s service) RegisterStudent(user *model.User, checker string, key string, vcode int) error {
	// verify vcode
	if !s.baseSrv.VerifyCode(key, vcode) {
		return srverr.ErrVerify
	}
	if err := s.dao.CDelete(key); err != nil {
		logger.Log.Warn("cache key delete failed", zap.String("key", key), zap.Error(err))
	}
	// verify
	// verify checker user
	t, err := s.dao.Query(&userDao.FilterUser{
		UserNo: []string{user.UserNo, checker},
	})
	implement := t.([]*model.User)
	if len(implement) > 1 {
		return srverr.ErrRegisterExisted
	}
	if err != nil || len(implement) <= 0 || implement[0].UserNo != checker {
		logger.Log.Warn("student register: invalid checker", zap.String("checker", checker), zap.Error(err))
		return srverr.ErrRegisterChecker
	}
	u := implement[0]
	if u.Status != model.Teacher {
		return srverr.ErrRegisterChecker
	}
	// verify class
	if user == nil {
		logger.Log.Warn("structure user is nil")
		return srverr.ErrSystemException
	}
	err = s.classSrv.CheckClass(user.Class)
	if err != nil {
		logger.Log.Warn("class not exist", zap.String("class", user.Class), zap.Error(err))
	}
	// convert json
	user.Status = model.Student
	j, err := json.Marshal(user)
	if err != nil {
		logger.Log.Warn("student register: marshal failed", zap.Any("user", user), zap.Error(err))
		return srverr.ErrSystemException
	}

	// put register info to cache for checking
	cacheKey := "REGISTER_CHECK-" + strconv.Itoa(int(u.ID))
	field := strconv.Itoa(int(user.ID))
	return s.dao.CHashAdd(cacheKey, field, string(j))
}

func (s service) Login(userNo, password, key string, vCode int) (token string, err error) {
	// verify vCode
	if !s.baseSrv.VerifyCode(key, vCode) {
		return "", srverr.ErrVerify
	}
	// verify user existence
	filter := &userDao.FilterUser{
		UserNo: []string{userNo},
	}
	obj, err := s.dao.Query(filter)
	if err != nil || obj == nil {
		logger.Log.Warn("login failed", zap.Any("filter", filter), zap.Error(err))
		return "", srverr.ErrLoginFailed

	}
	user, ok := obj.([]*model.User)
	if !ok || len(user) <= 0 {
		logger.Log.Warn("login failed", zap.Any("filter", filter), zap.Error(err))
		return "", srverr.ErrLoginFailed
	}
	// verify password
	if ok := rsa.Compare(password, user[0].Password); !ok {
		return "", srverr.ErrLoginFailed
	}
	// generate token
	token, err = jwt.Token(map[string]interface{}{
		"uid": user[0].ID,
		"rid": model.Visitor,
		"u_rid": user[0].Status,

	})
	if err != nil {
		logger.Log.Warn("login failed", zap.Uint("userId", user[0].ID), zap.Error(err))
		return "", srverr.ErrSystemException
	}

	return
}

func (s *service) Info(uid uint) *model.User {
	raw, err := s.dao.Query(&userDao.FilterUser{
		ID: []uint{uid},
	})
	if err != nil {
		logger.Log.Warn("query user error", zap.Uint("id", uid), zap.Error(err))
		return nil
	}
	user, ok := raw.([]*model.User)
	if !ok || len(user) <= 0 {
		logger.Log.Warn("query user error", zap.Uint("id", uid), zap.Error(err))
		return nil
	}

	return user[0]
}