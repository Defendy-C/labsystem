package admin

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"labsystem/dao"
	adminDao "labsystem/dao/admin"
	"labsystem/model"
	srv "labsystem/service"
	classSrv "labsystem/service/class"
	userSrv "labsystem/service/user"
	"labsystem/util/rsa"
	"testing"
	"time"
)

type TestAdminDao interface {
	adminDao.DaoAdmin
	dao.TestDao
}

type AdminSrvSuite struct {
	suite.Suite
	rootId uint
	*require.Assertions
	srv ServiceAdmin
	dao TestAdminDao
}

func (s *AdminSrvSuite) SetT(t *testing.T) {
	s.Suite.SetT(t)
	s.Assertions = require.New(t)
	s.dao = adminDao.NewAdminDao()
	base := srv.NewService()
	classService := classSrv.NewClassService()
	userService := userSrv.NewUserService(base, classService)
	s.srv = NewAdminService(userService, classService, base)
}

func (s *AdminSrvSuite) SetupTest() {
	s.NoError(s.dao.Truncate())
	enPwd, err :=  rsa.Encrypt("123456")
	s.NoError(err)
	root := &model.Admin{
		NickName:  "root",
		Password:  enPwd,
		Power:     model.PowerAll,
	}
	s.NoError(s.dao.Create(root))
	s.rootId = root.ID
}

func TestAdminSrv(t *testing.T) {
	suite.Run(t, new(AdminSrvSuite))
}

func (s *AdminSrvSuite) TestCreateAdmin() {
	operator1 := &model.Admin{
		NickName:  "operator1",
		Password:  "123456",
		Power:     model.PowerCreateAdmin,
		CreatedBy: "root",
	}
	operator2 := &model.Admin{
		NickName:  "operator2",
		Password:  "123456",
		Power:     model.PowerCreateTeacher,
		CreatedBy: "root",
	}
	s.NoError(s.srv.CreateAdmin(operator1))
	s.NoError(s.srv.CreateAdmin(operator2))
	s.NoError(s.srv.CreateAdmin(&model.Admin{
		NickName: "admin",
		Password: "123456",
		Power: model.PowerCreateTeacher,
		CreatedBy: "operator2",
	}))
	s.Error(s.srv.CreateAdmin(&model.Admin{
		NickName: "admin",
		Password: "123456",
		Power: model.PowerCreateTeacher,
		CreatedBy: "operator1",
	}))
	s.NoError(s.srv.CreateAdmin(&model.Admin{
		NickName: "admin2",
		Password: "123456",
		Power: model.PowerCreateAdmin,
		CreatedBy: "operator1",
	}))
}

func (s *AdminSrvSuite) TestLogin() {
	s.NoError(s.dao.CSet("1", "120", time.Minute))
	s.NoError(s.srv.CreateAdmin(&model.Admin{
		BaseModel: model.BaseModel{ID: 3},
		NickName:  "test02",
		Password:  "MpJucmZlKGm5CCqKKxrQNUhr8zT+lMvNKowpo7LKfKJVcPcwpVu52DhAxdFkxcqrlkcH3B5fwOqfPLrYwFl0W60JIYYJ/kHw6n7wKEby9Bw0nXuBhEbtHb3o3eQMLBS1LEJ4HVyzZrB9byno+2DE3NTf+HcN7DqfU8ev3GyZpVA=",
		Power:     model.PowerCreateAdmin,
		CreatedBy: "root",
	}))
	t, err := s.srv.Login("test02", "aj6kNOi+Sja6N2UHYRqtHiK21WVqCid3VcJujNpOLaTfr+L4tP9MVE+QR8mvjRsnpXa5kY57LpN7IwMs/yyI8CA4gvKa+/L5f/N5os+ZHdZCXj2j5Ku8kAh93k3BhYESg/tiQ6++YtiyDVIe5802Hs7KKjummwXtStgm5DB4rvA=", "1", 20)
	s.NoError(err)
	s.NotEqual("", t)
	s.Equal("", s.dao.CGet("1"))
}

func (s *AdminSrvSuite) TestQueryAdminByName() {
	admin := &model.Admin{
		NickName: "test",
		CreatedBy: "root",
		Power: model.PowerAll,
		Password: "123456",
	}
	s.NoError(s.dao.Create(admin))
	res := s.srv.QueryAdminByName("test")
	s.NotNil(res)
	s.Equal("test", res.NickName)
}

func (s *AdminSrvSuite) TestAdminList() {
	list := []*model.Admin {
		{
			NickName: "test01",
			CreatedBy: "root",
			Power: model.PowerAll,
			Password: "123456",
		},
		{
			NickName: "test02",
			CreatedBy: "root",
			Power: model.PowerAll,
			Password: "123456",
		},
		{
			NickName: "test03",
			CreatedBy: "root",
			Power: model.PowerAll,
			Password: "123456",
		},
		{
			NickName: "test04",
			CreatedBy: "root",
			Power: model.PowerAll,
			Password: "123456",
		},
		{
			NickName: "test05",
			CreatedBy: "root",
			Power: model.PowerAll,
			Password: "123456",
		},
		{
			NickName: "test06",
			CreatedBy: "root1",
			Power: model.PowerAll,
			Password: "123456",
		},
	}
	for _, v := range list {
		s.NoError(s.dao.Create(v))
	}
	res, totalPage, totalCount := s.srv.AdminList(nil, 2, 3)
	s.NotNil(res)
	s.Equal(3, int(totalPage))
	s.Equal(7, int(totalCount))
	res, totalPage, totalCount = s.srv.AdminList(&ListOpt{
		CreatedBy: "root",
		OrderFiled: CREATEDAT,
		OrderDesc: true,
	}, 2, 2)
	s.NotNil(res)
	s.Equal(3, int(totalPage))
	s.Equal(5, int(totalCount))
	for i := 1;i < len(res);i++ {
		s.True(res[i].CreatedAt.Before(res[i - 1].CreatedAt))
	}
}

func (s *AdminSrvSuite) TestDeleteAdmin() {
	admin := &model.Admin {
			NickName: "test01",
			CreatedBy: "root",
			Power: model.PowerAll,
			Password: "123456"}
	fateRoot := &model.Admin {
		NickName: "fake_root",
		CreatedBy: "",
		Power: model.PowerAll,
		Password: "123456"}
	s.NoError(s.dao.Create(&admin))
	s.NoError(s.dao.Create(&fateRoot))
	s.False(s.srv.DeleteAdmin(fateRoot.ID, admin.ID))
	s.True(s.srv.DeleteAdmin(s.rootId, admin.ID))
}

func (s *AdminSrvSuite) TestUpdatePower() {
	admins := []*model.Admin {
		{
			NickName: "operator1",
			CreatedBy: "root",
			Power: model.PowerCreateAdmin,
			Password: "123456",
		},
		{
			NickName: "operator2",
			CreatedBy: "root",
			Power: model.PowerAll,
			Password: "123456",
		},
		{
			NickName: "admin1",
			CreatedBy: "operator2",
			Power: model.PowerCreateTeacher,
			Password: "123456",
		},
		{
			NickName: "admin2",
			CreatedBy: "operator1",
			Power: model.PowerCreateAdmin,
			Password: "123456",
		},
	}
	for _, v := range admins {
		s.NoError(s.srv.CreateAdmin(v))
	}
	s.False(s.srv.UpdatePower(admins[0].ID, admins[2].ID, 2, 4))
	s.False(s.srv.UpdatePower(admins[1].ID, admins[2].ID, int(model.PowerNone), int(model.PowerCreateAdmin)))
	s.False(s.srv.UpdatePower(admins[0].ID, admins[3].ID, int(model.PowerCreateTeacher), int(model.PowerNone)))
	s.False(s.srv.UpdatePower(admins[1].ID, admins[2].ID, int(model.PowerCreateTeacher), int(model.PowerNone)))
	s.True(s.srv.UpdatePower(admins[1].ID, admins[2].ID, int(model.PowerCreateAdmin), int(model.PowerCreateTeacher)))
}

func (s *AdminSrvSuite) TestUpdateAdmin() {
	admin := model.Admin{
		NickName: "admin",
		CreatedBy: "root",
		Power: model.PowerCreateTeacher,
		Password: "123456",
	}
	s.NoError(s.dao.Create(&admin))
	s.True(s.srv.UpdateAdmin(admin.ID, "admin1", "234567"))
	list, err := s.dao.Query(&adminDao.FilterAdmin{
		ID: []uint{admin.ID},
	})
	s.NoError(err)
	s.Len(list, 1)
	obj := list.([]*model.Admin)[0]
	s.Equal("admin1", obj.NickName)
	s.Equal("234567", obj.Password)
	s.True(s.srv.UpdateAdmin(admin.ID, "admin2", ""))
	list, err = s.dao.Query(&adminDao.FilterAdmin{
		ID: []uint{admin.ID},
	})
	s.NoError(err)
	s.Len(list, 1)
	obj = list.([]*model.Admin)[0]
	s.Equal("admin2", obj.NickName)
	s.Equal("234567", obj.Password)
}