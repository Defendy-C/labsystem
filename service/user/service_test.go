package user

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	classSrv "labsystem/service/class"
	"labsystem/util/rsa"
	"strconv"
	"testing"
	"time"

	"labsystem/dao"
	adminDao "labsystem/dao/admin"
	classDao "labsystem/dao/class"
	userDao "labsystem/dao/user"
	"labsystem/model"
	commonSrv "labsystem/service"
)

type TestUserDao interface {
	userDao.DaoUser
	dao.TestDao
}

type TestAdminDao interface {
	adminDao.DaoAdmin
	dao.TestDao
}

type TestClassDao interface {
	classDao.DaoClass
	dao.TestDao
}

type TestUserSrvSuite struct {
	suite.Suite
	*require.Assertions
	dao TestUserDao
	adminDao TestAdminDao
	userSrv ServiceUser
	commonSrv commonSrv.Service
	classDao TestClassDao
}

func (s *TestUserSrvSuite) SetT(t *testing.T) {
	s.classDao = classDao.NewClassDao()
	s.dao = userDao.NewUserDao()
	s.adminDao = adminDao.NewAdminDao()
	s.commonSrv = commonSrv.NewService()
	s.userSrv = NewUserService(s.commonSrv, classSrv.NewClassService())
	s.Suite.SetT(t)
	s.Assertions = require.New(t)
}

func (s *TestUserSrvSuite) SetupTest() {
	s.NoError(s.dao.Truncate())
	s.NoError(s.adminDao.Truncate())
	s.NoError(s.classDao.Truncate())
	s.NoError(s.dao.FlushDB())
}

func (s *TestUserSrvSuite) TearDownTest() {

}

func (s *TestUserSrvSuite) TestRegister() {
	s.NoError(s.dao.CSet("1", "120", time.Minute))
	pwd, err := rsa.Encrypt("123456")
	s.NoError(s.adminDao.Create(&model.Admin{
		BaseModel:model.BaseModel{ID: 1},
		NickName: "root",
		Password: pwd,
		Power: model.PowerAll,
		CreatedBy: 1,
	}))
	s.NoError(s.classDao.Create(&model.Class{
		BaseModel: model.BaseModel{ID: 1},
		ClassNo: "1701",
	}))
	s.NoError(s.dao.Create(&model.User{
		UserNo: "17111000",
		RealName: "zhang san",
		Password: pwd,
		Class: "1234",
		Status: model.Teacher,
		ProfileUrl: "123.com",
		CreatedBy: "1234",
	}))
	s.NoError(s.userSrv.RegisterStudent(&model.User{
		UserNo: "17111001",
		RealName: "li si",
		Password: pwd,
		Class: "1234",
		ProfileUrl: "123.com",
	}, "17111000", "1", 20))
	us, err := s.userSrv.CheckList(1)
	s.NoError(err)
	s.Len(us, 1)
	s.Equal("17111001", us[0].UserNo)
}

func (s *TestUserSrvSuite) TestCheckStudent() {
	s.NoError(s.dao.CSet("1", "120", time.Minute))
	pwd, err := rsa.Encrypt("123456")
	s.NoError(err)
	s.NoError(s.adminDao.Create(&model.Admin{
		BaseModel:model.BaseModel{ID: 1},
		NickName: "root",
		Password: pwd,
		Power: model.PowerAll,
		CreatedBy: 1,
	}))
	s.NoError(s.classDao.Create(&model.Class{
		BaseModel: model.BaseModel{ID: 1},
		ClassNo: "1701",
	}))
	s.NoError(s.dao.Create(&model.User{
		UserNo: "17111000",
		RealName: "zhang san",
		Password: pwd,
		Class: "1234",
		Status: model.Teacher,
		ProfileUrl: "123.com",
		CreatedBy: "1234",
	}))
	u := &model.User{
		UserNo: "17111001",
		RealName: "li si",
		Password: pwd,
		Class: "1234",
		ProfileUrl: "123.com",
	}
	s.NoError(s.userSrv.RegisterStudent(u, "zhang san", "1", 20))
	s.NoError(s.userSrv.CheckStudent("1", strconv.Itoa(int(u.ID))))
	us, err := s.dao.Query(&userDao.FilterUser{
		ID: []uint{1},
	})
	s.NoError(err)
	s.Len(us.([]*model.User), 1)
}

func (s *TestUserSrvSuite) TestLogin() {
	s.NoError(s.dao.CSet("1", "120", time.Minute))
	pwd, err := rsa.Encrypt("123456")
	s.NoError(err)
	s.NoError(s.adminDao.Create(&model.Admin{
		BaseModel:model.BaseModel{ID: 1},
		NickName: "root",
		Password: pwd,
		Power: model.PowerAll,
		CreatedBy: 1,
	}))
	s.NoError(s.classDao.Create(&model.Class{
		BaseModel: model.BaseModel{ID: 1},
		ClassNo: "1701",
	}))
	s.NoError(s.dao.Create(&model.User{
		UserNo: "17111000",
		RealName: "zhang san",
		Password: pwd,
		Class: "1234",
		Status: model.Teacher,
		ProfileUrl: "123.com",
		CreatedBy: "1234",
	}))
	u := &model.User{
		UserNo: "17111001",
		RealName: "li si",
		Password: pwd,
		Class: "1234",
		ProfileUrl: "123.com",
	}
	s.NoError(s.userSrv.RegisterStudent(u, "zhang san", "1", 20))
	s.NoError(s.userSrv.CheckStudent("1", strconv.Itoa(int(u.ID))))
	us, err := s.dao.Query(&userDao.FilterUser{
		ID: []uint{1},
	})
	s.NoError(err)
	s.Len(us.([]*model.User), 1)
	key := s.commonSrv.GenerateVCode()
	vc := s.dao.CGet(key)
	val, err := strconv.Atoi(vc)
	s.NoError(err)
	t, err := s.userSrv.Login("17111001", pwd, key, val)
	s.NoError(err)
	s.NotEqual("", t)
}

func (s *TestUserSrvSuite) TestInfo() {
	pwd, err := rsa.Encrypt("123456")
	s.NoError(err)
	user := &model.User{
		UserNo: "17111000",
		RealName: "zhang san",
		Password: pwd,
		Class: "1234",
		Status: model.Teacher,
		ProfileUrl: "123.com",
		CreatedBy: "1234",
	}
	s.NoError(s.dao.Create(user))
	res := s.userSrv.Info(user.ID)
	s.NotNil(res)
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(TestUserSrvSuite))
}