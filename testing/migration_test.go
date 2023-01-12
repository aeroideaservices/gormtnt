package testing

import "testing"

func TestMigration(t *testing.T) {
	db := initDB(t)

	t.Run("basic", func(t *testing.T) {
		m := &User{}
		err := db.AutoMigrate(m)
		if err != nil {
			t.Fatalf("unexpected error for AutoMigrate: %v", err)
		}
		err = db.Migrator().DropTable(&m)
		if err != nil {
			t.Fatalf("unexpected error for DropTable: %v", err)
		}
	})

	t.Run("uuid", func(t *testing.T) {
		m := &Uuid{}
		err := db.AutoMigrate(m)
		if err != nil {
			t.Fatalf("unexpected error for AutoMigrate: %v", err)
		}
		err = db.Migrator().DropTable(&m)
		if err != nil {
			t.Fatalf("unexpected error for DropTable: %v", err)
		}
	})
}
