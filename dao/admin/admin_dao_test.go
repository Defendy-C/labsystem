package admin

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"labsystem/dao"
	"labsystem/model"
	"testing"
	"time"
)

type AdminDaoTestSuite struct {
	db *adminDao
	suite.Suite
	*require.Assertions
}

func (s *AdminDaoTestSuite) SetT(t *testing.T) {
	s.db = NewAdminDao()
	s.Suite.SetT(t)
	s.Assertions = require.New(t)
}

func (s *AdminDaoTestSuite) SetupTest() {
	s.db.SQL.Exec("truncate " + admin)
}

func (s *AdminDaoTestSuite) TearDownTest() {
	s.db.Clear()
}

func (s *AdminDaoTestSuite) TestBase() {
	testcase := []*model.Admin{
		{
			BaseModel: model.BaseModel{
				CreatedAt: time.Now().Add(-1 * model.Day),
			},
			NickName:  "testcase01",
			Password:  "123456",
			Power:     0,
			CreatedBy: "root",
		}, {
			BaseModel: model.BaseModel{
				CreatedAt: time.Now().Add(-2 * model.Day),
			},
			NickName:  "testcase02",
			Password:  "123456",
			Power:     0,
			CreatedBy: "root",
		}, {
			BaseModel: model.BaseModel{
				CreatedAt: time.Now().Add(-3 * model.Day),
			},
			NickName:  "testcase03",
			Password:  "123456",
			Power:     0,
			CreatedBy: "root",
		},
	}

	// test create
	for i := 0; i < 3; i++ {
		s.NoError(s.db.Create(testcase[i]))
	}
	data, err := s.db.Query(&FilterAdmin{
		BaseFilter: dao.BaseFilter{Sort: dao.NewOrderBy("id", dao.DESC)},
	})
	s.NoError(err)
	for i := 1; i < len(data.([]*model.Admin)); i++ {
		s.True(data.([]*model.Admin)[i].ID < data.([]*model.Admin)[i-1].ID)
	}

	// test page
	filter := &FilterAdmin{
		BaseFilter: dao.BaseFilter{Sort: dao.NewOrderBy("id", dao.DESC), Page: 2, PerPage: 2},
	}
	data, err = s.db.Query(filter)
	s.NoError(err)
	s.Equal(1, len(data.([]*model.Admin)))
	s.Equal(testcase[0].ID, data.([]*model.Admin)[0].ID)
	s.Equal(2, int(filter.TotalPage))
	s.Equal(3, int(filter.TotalCount))

	// test update
	s.db.Update(map[string]interface{}{
		"nick_name": "testcase01",
	}, map[string]interface{}{
		"password": "234567",
	})
	// test query
	data, err = s.db.Query(&FilterAdmin{
		NickName: []string{"testcase01"},
	})
	s.NoError(err)
	s.Equal(1, len(data.([]*model.Admin)))
	s.Equal("234567", data.([]*model.Admin)[0].Password)
	s.Equal("root", data.([]*model.Admin)[0].CreatedBy)
	timeMin, err := time.Parse("2006-01-02", "2021-04-03")
	s.NoError(err)
	timeMax, err := time.Parse("2006-01-02", "2021-04-04")
	f := &FilterAdmin{}
	f.SetCreatedAtRange(&timeMin, &timeMax)
	data, err = s.db.Query(f)
	s.Len(data, 1)
	admin := data.([]*model.Admin)[0]
	s.True(admin.CreatedAt.After(timeMin) && admin.CreatedAt.Before(timeMax))
	// test delete
	s.db.Delete(nil)
	data, err = s.db.Query(&FilterAdmin{})
	s.NoError(err)
	s.Equal(0, len(data.([]*model.Admin)))
}

func (s *AdminDaoTestSuite)TestCache() {
	err := s.db.CSet("1", "hello", 0)
	s.NoError(err)
	res := s.db.CGet("1")
	s.Equal("hello", res)
}

func TestAdminDaoSuite(t *testing.T) {
	suite.Run(t, new(AdminDaoTestSuite))
}
