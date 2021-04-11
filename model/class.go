package model

type Class struct {
	BaseModel
	ClassNo string `gorm:"unique"`
}