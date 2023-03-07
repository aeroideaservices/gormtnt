# GORM tnt

## Example of usage

```go
package main

import (
	"log"
	"reflect"

	"github.com/aeroideaservices/gormtnt"
	"gorm.io/gorm"
)

type Post struct {
	ID       int64  `json:"id"`
	Text     string `json:"text"`
	AuthorID int64  `json:"author_id"`
	Author   *User  `josn:"author"`
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func main() {
	db, _ := gorm.Open(gormtnt.Open("tarantool://admin:pass@localhost:3301"), &gorm.Config{})

	db.AutoMigrate(&User{}, &Post{})
	defer db.Migrator().DropTable(&User{}, &Post{})

	if err := db.Create(users).Error; err != nil {
		log.Fatal(err)
	}

	if err := db.Create(posts).Error; err != nil {
		log.Fatal(err)
	}

	var p Post
	// Не забываем про двойные кавычки для идентификатора
	if err := db.Model(&Post{}).Select("*").Preload("Author").Where(`"id" = ?`, 1).First(&p).Error; err != nil {
		log.Fatal(err)
	}

	log.Print(reflect.DeepEqual(p, Post{
		ID:       1,
		Text:     "New database/sql driver for tarantool!",
		AuthorID: 1,
		Author: &User{
			ID:   1,
			Name: "Alice",
		},
	}))

}

var posts = []Post{
	{
		ID:       1,
		Text:     "New database/sql driver for tarantool!",
		AuthorID: 1,
	},
	{
		ID:       2,
		Text:     "New gorm driver for tarantool!",
		AuthorID: 2,
	},
	{
		ID:       3,
		Text:     "Another post",
		AuthorID: 2,
	},
}

var users = []User{
	{
		ID:   1,
		Name: "Alice",
	},
	{
		ID:   2,
		Name: "Bob",
	},
	{
		ID:   3,
		Name: "John",
	},
}
```