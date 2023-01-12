package testing

import (
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestCreate(t *testing.T) {
	db := initDB(t)

	t.Run("basic", func(t *testing.T) {
		teardown := setupCreate(t, db, &User{})
		defer teardown()
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
	})

	t.Run("uuid", func(t *testing.T) {
		teardown := setupCreate(t, db, &Uuid{})
		defer teardown()

		// _ = setupCreate(t, db, &Uuid{})

		u := Uuid{
			ID: func() uuid.UUID {
				x, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")
				return x
			}(),
			AnotherUUID: func() *uuid.UUID { x := uuid.New(); return &x }(),
			Name:        "Alice",
		}

		t.Run("Single", func(t *testing.T) {
			err := db.Create(&u).Error
			if err != nil {
				t.Fatalf("unexpected error for Create: %v", err)
			}
		})

		us := []Uuid{
			{
				ID: func() uuid.UUID {
					x, _ := uuid.Parse("00000000-0000-0000-0000-000000000002")
					return x
				}(),
				AnotherUUID: func() *uuid.UUID { x := uuid.New(); return &x }(),
				Name:        "Bob",
			},
			{
				ID: func() uuid.UUID {
					x, _ := uuid.Parse("00000000-0000-0000-0000-000000000003")
					return x
				}(),
				AnotherUUID: func() *uuid.UUID { x := uuid.New(); return &x }(),
				Name:        "Martha",
			},
			{
				ID: func() uuid.UUID {
					x, _ := uuid.Parse("00000000-0000-0000-0000-000000000004")
					return x
				}(),
				AnotherUUID: nil,
				Name:        "John",
			},
		}

		t.Run("Multiple", func(t *testing.T) {
			err := db.Create(us).Error
			if err != nil {
				t.Fatalf("unexpected error for Create: %v", err)
			}
		})
	})
}

func setupCreate(t *testing.T, db *gorm.DB, model any) (teardown func()) {
	err := db.AutoMigrate(model)
	if err != nil {
		t.Fatalf("unexpected error for AutoMigrate: %v", err)
	}
	teardown = func() {
		db.Migrator().DropTable(model)
	}
	return
}
