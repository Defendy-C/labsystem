package admin

import (
	"gorm.io/gorm"
	"labsystem/dao"
	"labsystem/model"
)

type FilterAdmin struct {
	dao.BaseFilter
	ID        []uint
	NickName  []string
	CreatedBy []string
}
type DaoAdmin interface {
	dao.BaseDao
	dao.BaseCacheDao
}

type adminDao struct {
	*dao.DAO
}

func NewAdminDao() *adminDao {
	return &adminDao{dao.NewDAO()}
}

var _ DaoAdmin = &adminDao{}

const admin = "admins"

func (d *adminDao) Create(obj interface{}) error {
	return d.SQL.Table(admin).Create(obj).Error
}

func (d *adminDao) Query(filter dao.Filter) (interface{}, error) {
	db := d.SQL.Table(admin)
	if filter == nil {
		return nil, gorm.ErrInvalidData
	}
	adminFilter, ok := filter.(*FilterAdmin)
	if !ok {
		return nil, gorm.ErrInvalidData
	}
	switch len(adminFilter.ID) {
	case 0:
	case 1:
		db = db.Where("id = ?", adminFilter.ID[0])
	default:
		db = db.Where("id in (?)", adminFilter.ID)
	}
	switch len(adminFilter.NickName) {
	case 0:
	case 1:
		db = db.Where("nick_name = ?", adminFilter.NickName)
	default:
		db = db.Where("nick_name in (?)", adminFilter.NickName)
	}
	if len(adminFilter.CreatedBy) > 0 {
		db = db.Where("created_by in (?)", adminFilter.CreatedBy)
	}
	if adminFilter.Sort != nil {
		db = db.Order(adminFilter.OrderBy())
	}
	db = db.Scopes(adminFilter.PageScope)

	var admins []*model.Admin
	return admins, db.Find(&admins).Error
}

func (d *adminDao) Update(cond map[string]interface{}, changed map[string]interface{}) error {
	return d.SQL.Table(admin).Where(cond).Updates(changed).Error
}

func (d *adminDao) Delete(fields map[string]interface{}) error {
	db := d.SQL.Where("1 = 1") // soft delete require where
	if fields != nil {
		db = db.Where(fields)
	}
	return db.Delete(&model.Admin{}).Error
}

func (d *adminDao) Truncate() error {
	return d.SQL.Exec("truncate " + admin).Error
}

func (d *adminDao) Clear() {
	d.DAO.Clear([]string{"admin"})
}
