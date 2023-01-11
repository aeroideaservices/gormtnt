package testing

import (
	"testing"

	"gorm.io/gorm"
)

func TestCreate(t *testing.T) {
	db := initDB(t)
	teardown := setupCreate(t, db)
	_ = teardown
	// defer teardown()
	u := User{
		ID:   1,
		Name: "Alice",
	}

	t.Run("Singe", func(t *testing.T) {
		err := db.Create(&u).Error
		if err != nil {
			t.Fatalf("unexpected error for Create: %v", err)
		}
	})

	us := []User{
		{
			ID:   2,
			Name: "Bob",
		},
		{
			ID:   3,
			Name: "Martha",
		},
	}

	t.Run("Multiple", func(t *testing.T) {
		err := db.Create(us).Error
		if err != nil {
			t.Fatalf("unexpected error for Create: %v", err)
		}
	})
}

func setupCreate(t *testing.T, db *gorm.DB) (teardown func()) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		t.Fatalf("unexpected error for AutoMigrate: %v", err)
	}
	teardown = func() {
		db.Migrator().DropTable(&User{})
	}
	return
}
