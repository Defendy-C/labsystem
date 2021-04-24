package model


type User struct {
	BaseModel
	UserNo     string `gorm:"unique"`
	RealName   string `gorm:"unique"`
	Password   string
	Status     UserStatus
	Class      string
	ProfileUrl string `go:"unique"` // profile url
	CreatedBy  string
}

// status
type UserStatus uint

func (s *UserStatus)Uint() uint {
	return uint(*s)
}

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
