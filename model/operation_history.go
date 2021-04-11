package model

import (
	"strconv"
)

type operation string

const (
	Add    = "add"
	Del    = "delete"
	Modify = "modify"
)

type operator string

const (
	Users  operator = "user"
	Admins operator = "admin"
)

type OperationHistory struct {
	BaseModel
	Operator  string
	Operation string
}

// status
func NewOperationHistory(status UserStatus, id uint, object string, opt operation) *OperationHistory {
	operator := StatusToString(status) + " - " + strconv.Itoa(int(id))
	op := object + " - " + string(opt)

	return &OperationHistory{Operator: operator, Operation: op}
}
