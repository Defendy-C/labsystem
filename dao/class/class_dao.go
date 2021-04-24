package class

import (
	"gorm.io/gorm"
	"labsystem/dao"
	"labsystem/model"
)

type DaoClass interface {
	dao.BaseDao
}

type classDao struct {
	*dao.DAO
}

type FilterClass struct {
	*dao.BaseFilter
	Id []uint
	ClassNo []string
}

func NewClassDao()*classDao {
	return &classDao{dao.NewDAO()}
}
var _ DaoClass = &classDao{}

func (c classDao) Create(i interface{}) error {
	return c.SQL.Create(i).Error
}

func (c classDao) Query(filter dao.Filter) (interface{}, error) {
	classFilter, ok := filter.(*FilterClass)
	if !ok {
		return nil, gorm.ErrInvalidData
	}
	db := c.SQL.Model(&model.Class{})
	if len(classFilter.ClassNo) > 0 {
		db = db.Where("class_no in (?)", classFilter.ClassNo)
	}
	switch len(classFilter.Id) {
	case 0:
	case 1:
		db = db.Where("id = ?", classFilter.Id[0])
	default:
		db = db.Where("id in (?)", classFilter.Id)
	}
	if classFilter.BaseFilter != nil {
		if classFilter.Sort != nil {
			db = db.Order(classFilter.OrderBy())
		}
		db = db.Scopes(filter.PageScope)
	}

	var data []*model.Class
	return data, db.Find(&data).Error
}

func (c classDao) Update(m map[string]interface{}, changed map[string]interface{}) error {
	panic("implement me")
}

func (c classDao) Delete(m map[string]interface{}) error {
	panic("implement me")
}

func (c classDao) Clear() {
	panic("implement me")
}

func (c classDao) Truncate() error {
	return c.SQL.Exec("truncate classes").Error
}