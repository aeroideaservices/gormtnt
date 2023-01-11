package testing

type User struct {
	ID   int
	Name string
}

func (User) TableName() string {
	return "userTable"
}
