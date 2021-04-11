package user

import (
	"gorm.io/gorm"
	"labsystem/dao"
	"labsystem/model"
)

type DaoUser interface {
	dao.BaseDao
	dao.BaseCacheDao
}

var _ DaoUser = &userDao{}

type userDao struct {
	*dao.DAO
}

func NewUserDao() *userDao {
	return &userDao{dao.NewDAO()}
}

type FilterUser struct {
	dao.BaseFilter
	ID    []int
	UserNo [] string
	RealName []string
}
const user = "users"

func (u userDao) Create(i interface{}) error {
	return u.SQL.Create(i).Error
}

func (u userDao) Query(filter dao.Filter) (interface{}, error) {
	if filter == nil {
		return nil, gorm.ErrInvalidData
	}
	userFilter, ok := filter.(*FilterUser)
	if !ok {
		return nil, gorm.ErrInvalidData
	}
	db := u.SQL
	switch len(userFilter.ID) {
	case 0:
	case 1:
		db = db.Where("id = ?", userFilter.ID[0])
	default:
		db = db.Where("id in (?)", userFilter.ID)
	}
	if len(userFilter.UserNo) > 0 {
		db = db.Where("user_no in (?)", userFilter.UserNo)
	}
	if len(userFilter.RealName) > 0 {
		db = db.Where("real_name in (?)", userFilter.RealName)
	}
	if userFilter.Sort != nil {
		db = db.Order(userFilter.OrderBy())
	}

	db = db.Scopes(userFilter.PageScope)

	var users []*model.User
	return users, db.Find(&users).Error
}

func (u userDao) Update(cond map[string]interface{}, changed map[string]interface{}) error {
	panic("implement me")
}

func (u userDao) Delete(m map[string]interface{}) error {
	panic("implement me")
}

func (u userDao) Clear() {
	u.DAO.Clear([]string{"user"})
}

func (u userDao) Truncate() error {
	return u.SQL.Exec("truncate " + user).Error
}