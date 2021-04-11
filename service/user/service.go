package user

import (
	"encoding/json"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	userDao "labsystem/dao/user"
	"labsystem/model"
	"labsystem/model/srverr"
	"labsystem/util/jwt"
	"labsystem/util/logger"
	"strconv"
)

type ServiceUser interface {
	RegisterStudent(user *model.User, checker int) error
	CheckStudent(operator string, user string) error
	CheckList(operator int) ([]*model.User, error)
	Login(userNo, password, key string, vCode string) (token string, err error)
}

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

func (s service) RegisterStudent(user *model.User, checker int) error {
	// verify checker
	t, err := s.dao.Query(&userDao.FilterUser{
		ID: []int{checker},
	})
	if err != nil || len(t.([]*model.User)) <= 0 {
		logger.Log.Warn("student register: invalid checker", zap.Int("checker", checker), zap.Error(err))
		return srverr.ErrRegisterChecker
	}
	// verify class
	if user == nil {
		logger.Log.Warn("structure user is nil")
		return srverr.ErrSystemException
	}
	err = s.classSrv.CheckClass(user.Class)
	if err != nil {
		logger.Log.Warn("class not exist", zap.Uint("class", user.Class), zap.Error(err))
	}
	// convert json
	user.Status = model.Student
	pwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Warn("password crypt failed", zap.String("pwd", user.Password))
		return srverr.ErrSystemException
	}
	user.Password = string(pwd)
	j, err := json.Marshal(user)
	if err != nil {
		logger.Log.Warn("student register: marshal failed", zap.Any("user", user), zap.Error(err))
		return srverr.ErrSystemException
	}

	// put register info to cache for checking
	return s.dao.CHashAdd("REGISTER_CHECK-" + strconv.Itoa(checker), strconv.Itoa(int(user.ID)), string(j))
}

func (s service) Login(userNo, password, key, vCode string) (token string, err error) {
	// verify vCode
	if v := s.dao.CGet(key); v != vCode {
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
	if err := bcrypt.CompareHashAndPassword([]byte(user[0].Password), []byte(password)); err != nil {
		logger.Log.Warn("login failed", zap.String("expected hash", user[0].Password), zap.String("actual password", password), zap.Error(err))
		return "", srverr.ErrLoginFailed
	}
	// generate token
	token, err = jwt.Token(map[string]interface{}{
		"Id": user[0].ID,
	})
	if err != nil {
		logger.Log.Warn("login failed", zap.Uint("userId", user[0].ID), zap.Error(err))
		return "", srverr.ErrSystemException
	}

	return
}
