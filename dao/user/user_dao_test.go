package user

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"labsystem/model"
	"testing"
)

type UserDaoTestSuite struct {
	db *userDao
	suite.Suite
	*require.Assertions
}

func (s *UserDaoTestSuite) SetT(t *testing.T) {
	s.Suite.SetT(t)
	s.Assertions = require.New(t)
}

func (s *UserDaoTestSuite) SetupTest() {
	s.db = NewUserDao()
	s.db.SQL.Exec("truncate " + user)
}

func (s *UserDaoTestSuite) TearDownTest() {
	s.db.Clear()
}

func (s *UserDaoTestSuite) TestCreate() {
	users := []*model.User{
		{
			BaseModel: model.BaseModel{ID: 1},
			UserNo: "17111000",
			RealName: "Alice One",
			Password: "123456",
			Status: model.Student,
			Class: "1",
			ProfileUrl: "123.com",
		}, {
			BaseModel: model.BaseModel{ID: 2},
			UserNo: "17111001",
			RealName: "Ben Two",
			Password: "123456",
			Status: model.Teacher,
			Class: "1",
			ProfileUrl: "123.com",
		}, {
			BaseModel: model.BaseModel{ID: 3},
			UserNo: "17111002",
			RealName: "Chale Three",
			Password: "123456",
			Status: model.Student,
			Class: "1",
			ProfileUrl: "123.com",
		}, {
			BaseModel: model.BaseModel{ID: 4},
			UserNo: "17111003",
			RealName: "Ding Four",
			Password: "123456",
			Status: model.Teacher,
			Class: "1",
			ProfileUrl: "123.com",
		},
	}

	for i := 0;i < len(users);i++ {
		s.NoError(s.db.Create(users[i]))
	}

	filter := &FilterUser{ID: []uint{1, 2}}
	list, err := s.db.Query(filter)
	s.NoError(err)
	s.Len(list, 2)
	filter = &FilterUser{UserNo: []string{"17111000", "17111001"}}
	list, err = s.db.Query(filter)
	s.NoError(err)
	s.Len(list, 2)
	filter = &FilterUser{RealName: []string{"Ding Four"}}
	list, err = s.db.Query(filter)
	s.NoError(err)
	s.Len(list, 1)
	user := list.([]*model.User)
	s.Equal(4, int(user[0].ID))
}

func TestUserDaoSuite(t *testing.T) {
	suite.Run(t, new(UserDaoTestSuite))
}