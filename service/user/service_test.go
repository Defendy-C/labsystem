package user

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	classSrv "labsystem/service/class"
	"testing"

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
	s.Suite.SetT(t)
	s.Assertions = require.New(t)
}

func (s *TestUserSrvSuite) SetupTest() {
	s.classDao = classDao.NewClassDao()
	s.dao = userDao.NewUserDao()
	s.adminDao = adminDao.NewAdminDao()
	s.userSrv = NewUserService(classSrv.NewClassService())
	s.commonSrv = commonSrv.NewService()
	s.NoError(s.dao.Truncate())
	s.NoError(s.adminDao.Truncate())
	s.NoError(s.classDao.Truncate())
	s.NoError(s.dao.FlushDB())
}

func (s *TestUserSrvSuite) TearDownTest() {

}

func (s *TestUserSrvSuite) TestRegister() {
	s.NoError(s.adminDao.Create(&model.Admin{
		BaseModel:model.BaseModel{ID: 1},
		NickName: "root",
		Password: "123456",
		Power: model.PowerAll,
		CreatedBy: "root",
	}))
	s.NoError(s.classDao.Create(&model.Class{
		BaseModel: model.BaseModel{ID: 1},
		ClassNo: "1701",
	}))
	s.NoError(s.dao.Create(&model.User{
		BaseModel:model.BaseModel{ID: 1},
		UserNo: "17111000",
		RealName: "zhang san",
		Password: "123456",
		Class: 1,
		Status: model.Teacher,
		ProfileUrl: "123.com",
		CreatedBy: 1,
	}))
	s.NoError(s.userSrv.RegisterStudent(&model.User{
		UserNo: "17111001",
		RealName: "li si",
		Password: "123456",
		Class: 1,
		ProfileUrl: "123.com",
	}, 1))
	us, err := s.userSrv.CheckList(1)
	s.NoError(err)
	s.Len(us, 1)
	s.Equal("17111001", us[0].UserNo)
}

func (s *TestUserSrvSuite) TestCheckStudent() {
	s.NoError(s.adminDao.Create(&model.Admin{
		BaseModel:model.BaseModel{ID: 1},
		NickName: "root",
		Password: "123456",
		Power: model.PowerAll,
		CreatedBy: "root",
	}))
	s.NoError(s.classDao.Create(&model.Class{
		BaseModel: model.BaseModel{ID: 1},
		ClassNo: "1701",
	}))
	s.NoError(s.dao.Create(&model.User{
		BaseModel:model.BaseModel{ID: 1},
		UserNo: "17111000",
		RealName: "zhang san",
		Password: "123456",
		Class: 1,
		Status: model.Teacher,
		ProfileUrl: "123.com",
		CreatedBy: 1,
	}))
	s.NoError(s.userSrv.RegisterStudent(&model.User{
		BaseModel: model.BaseModel{ID: 2},
		UserNo: "17111001",
		RealName: "li si",
		Password: "123456",
		Class: 1,
		ProfileUrl: "123.com",
	}, 1))
	s.NoError(s.userSrv.CheckStudent("1", "2"))
	us, err := s.dao.Query(&userDao.FilterUser{
		ID: []int{1},
	})
	s.NoError(err)
	s.Len(us.([]*model.User), 1)
}

func (s *TestUserSrvSuite) TestLogin() {
	s.NoError(s.adminDao.Create(&model.Admin{
		BaseModel:model.BaseModel{ID: 1},
		NickName: "root",
		Password: "123456",
		Power: model.PowerAll,
		CreatedBy: "root",
	}))
	s.NoError(s.classDao.Create(&model.Class{
		BaseModel: model.BaseModel{ID: 1},
		ClassNo: "1701",
	}))
	s.NoError(s.dao.Create(&model.User{
		BaseModel:model.BaseModel{ID: 1},
		UserNo: "17111000",
		RealName: "zhang san",
		Password: "123456",
		Class: 1,
		Status: model.Teacher,
		ProfileUrl: "123.com",
		CreatedBy: 1,
	}))
	s.NoError(s.userSrv.RegisterStudent(&model.User{
		BaseModel: model.BaseModel{ID: 2},
		UserNo: "17111001",
		RealName: "li si",
		Password: "123456",
		Class: 1,
		ProfileUrl: "123.com",
	}, 1))
	s.NoError(s.userSrv.CheckStudent("1", "2"))
	us, err := s.dao.Query(&userDao.FilterUser{
		ID: []int{1},
	})
	s.NoError(err)
	s.Len(us.([]*model.User), 1)
	key := s.commonSrv.GenerateVCode()
	vc := s.dao.CGet(key)
	t, err := s.userSrv.Login("17111001", "123456", key, vc)
	s.NoError(err)
	s.NotEqual("", t)
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(TestUserSrvSuite))
}