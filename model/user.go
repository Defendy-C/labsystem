package model


type User struct {
	BaseModel
	UserNo     string `gorm:"unique"`
	RealName   string `gorm:"unique"`
	Password   string
	Status     UserStatus
	Class      uint
	ProfileUrl string `go:"unique"` // profile url
	CreatedBy  uint
}

// status
type UserStatus uint

const (
	Teacher = iota
	Student
)

func StatusToString(s UserStatus) string {
	switch s {
	case 0:
		return "teacher"
	case 1:
		return "student"
	}

	return ""
}
