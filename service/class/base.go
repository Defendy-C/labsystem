package class

import (
	"labsystem/dao"
	classDao "labsystem/dao/class"
	"labsystem/model"
)

type InternalClassSrv interface {
	CheckClass(name string) error
	CreateClass(class *model.Class) error
	List(opt *ListOpt, page, pageSize uint) (list []*model.Class, totalPage, totalCount uint)
}

type ServiceClass interface {

}


type service struct {
	dao classDao.DaoClass
}

func NewClassService() *service {
	return &service{
		dao: classDao.NewClassDao(),
	}
}

type OrderField string
const (
	CLASSNO = "class_no"
)
func (f OrderField)ToString() string {
	return string(f)
}

type ListOpt struct {
	OrderFiled OrderField
	OrderDesc  bool
}
func (opt *ListOpt) OrderPad() (string, dao.OrderTyp) {
	var typ dao.OrderTyp
	if opt.OrderDesc {
		typ = dao.DESC
	} else {
		typ = dao.ASC
	}

	return opt.OrderFiled.ToString(), typ
}

