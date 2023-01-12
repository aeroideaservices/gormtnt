package testing

import "github.com/google/uuid"

type User struct {
	ID   int
	Name string
}

func (User) TableName() string {
	return "UserTable"
}

type Uuid struct {
	ID          uuid.UUID  `gorm:"type:uuid"`
	AnotherUUID *uuid.UUID `gorm:"type:uuid"`
	Name        string
}

func (Uuid) TableName() string {
	return "UuidTest"
}
