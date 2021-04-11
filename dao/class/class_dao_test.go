package class

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"labsystem/dao"
	"labsystem/model"
	"testing"
)

type TestClassDao interface {
	DaoClass
	dao.TestDao
}

type TestClassDaoSuite struct {
	suite.Suite
	*require.Assertions
	dao TestClassDao
}

func (s *TestClassDaoSuite) SetT(t *testing.T) {
	s.Suite.SetT(t)
	s.Assertions = require.New(t)
}

func (s *TestClassDaoSuite) SetupTest() {
	s.dao = NewClassDao()
	s.NoError(s.dao.Truncate())
}

func (s *TestClassDaoSuite) TearDownTest() {

}

func (s *TestClassDaoSuite) TestClassCreate() {
	classes := []*model.Class {
		{
			BaseModel: model.BaseModel{ID: 1},
			ClassNo: "1701",
		},
		{
			BaseModel: model.BaseModel{ID: 2},
			ClassNo: "1702",
		},
		{
			BaseModel: model.BaseModel{ID: 3},
			ClassNo: "1703",
		},
		{
			BaseModel: model.BaseModel{ID: 4},
			ClassNo: "1704",
		},
		{
			BaseModel: model.BaseModel{ID: 5},
			ClassNo: "1705",
		},
	}

	for _, v := range classes {
		s.NoError(s.dao.Create(v))
	}

	list, err := s.dao.Query(&FilterClass{
		ClassNo: []string{"1703", "1705"},
	})
	s.NoError(err)
	s.Len(list, 2)
	list, err = s.dao.Query(&FilterClass{
		Id: []uint{1, 999},
	})
	s.NoError(err)
	s.Len(list, 1)
}

func TestClassSuite(t *testing.T) {
	suite.Run(t, new(TestClassDaoSuite))
}
