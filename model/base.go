package model

import (
	"gorm.io/gorm"
	"time"
)

// model
type BaseModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

var Models = map[string]interface{}{
	"admin":             Admin{},
	"user":              User{},
	"operation_history": OperationHistory{},
	"classes":           Class{},
}

// common
const (
	Day         = time.Hour * 24
	DateFormat  = "2006-01-02"
	TimeFormat  = DateFormat + " " + "15:04:05"
	MAXTryCount = 5
	KB          = 1024 * 8
)

// token
const (
	TokenExp = Day
)

// role
type Role uint
const (
	_ Role = iota
	Administrator
	Visitor
)
func (r Role)Int()int {
	return int(r)
}

// regexp
const (
	RegExpUserNo = `^\d{9}$`
	RegExpUserName = `^[0-9a-zA-Z]{4,10}$`
	RegExpPassword = `^\S{6,16}$`
)
