package model

import (
	"labsystem/model/srverr"
)

type AdminStatus uint

type Admin struct {
	BaseModel
	NickName  string `gorm:"unique"`
	Password  string
	Power     Power
	CreatedBy uint
}

type Power int

const (
	PowerNone        Power = 0
	PowerCreateAdmin Power      = 1 << iota
	PowerCreateTeacher
	PowerDeleteUser
	PowerCreateClass
	PowerAll = 1<<iota - 1
)

func New(ps ...Power) Power {
	p := PowerNone
	for _, v := range ps {
		p.Add(v)
	}

	return p
}

func IntToPower(pForInt int) (Power, error) {
	if pForInt > PowerAll {
		return PowerNone, srverr.ErrInvalidPower
	}

	return Power(pForInt), nil
}

func (p Power) Own(owned Power) bool {
	return owned&p == owned
}

func (p Power) Disown(power Power) bool  {
	return power&p == 0
}

func (p Power) Add(extra Power) Power {
	return p | extra
}

func (p Power) Sub(owned Power) Power {
	return p ^ owned
}

var PowerList = []struct {
	No   Power
	Name string
}{
	{
		PowerCreateAdmin,
		"管理员",
	},
	{
		PowerCreateTeacher,
		"创建教师用户",
	},
	{
		PowerCreateClass,
		"创建班级",
	},
}
