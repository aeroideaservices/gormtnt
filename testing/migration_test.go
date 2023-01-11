package testing

import "testing"

func TestMigration(t *testing.T) {
	db := initDB(t)

	m := &User{}
	err := db.AutoMigrate(m)
	if err != nil {
		t.Fatalf("unexpected error for AutoMigrate: %v", err)
	}
	err = db.Migrator().DropTable(&m)
	if err != nil {
		t.Fatalf("unexpected error for DropTable: %v", err)
	}
}
