package testing

import (
	"testing"

	"gorm.io/gorm"
)

/*
	Тут планировалось добавить тесты на прелоады, но времени не было и их уже потестили в рабочем окружении,
	так что это еще одна задача на будущее
*/

func TestBelongsTo(t *testing.T) {
	db := initDB(t)

	teardown := setupBelongsTo(t, db)
	defer teardown()

	// _ = setupBelongsTo(t, db)

	err := db.Create(&Client{
		ID:        3,
		Name:      "Homeless",
		CompanyID: 3,
	}).Error
	if err == nil {
		t.Fatalf("homeless client was created without an error")
	}

}

func setupBelongsTo(t *testing.T, db *gorm.DB) (teardown func()) {
	err := db.AutoMigrate(&Company{}, &Client{})
	if err != nil {
		t.Fatalf("auto migration error: %v", err)
	}

	for _, company := range []Company{
		{
			ID:   1,
			Name: "Ascaro Inc.",
		},
		{
			ID:   2,
			Name: "Aero",
		},
	} {
		tc := company
		err = db.Create(&tc).Error
		if err != nil {
			t.Fatalf("create error: %v", err)
		}
	}

	for _, client := range []Client{
		{
			ID:        1,
			Name:      "Ilya Scaro",
			CompanyID: 1,
		},
		{
			ID:        2,
			Name:      "Napas Lampas",
			CompanyID: 1,
		},
		{
			ID:        3,
			Name:      "Vova Medved",
			CompanyID: 2,
		},
	} {
		tc := client
		err = db.Create(&tc).Error
		if err != nil {
			t.Fatalf("create error: %v", err)
		}
	}
	teardown = func() {
		db.Migrator().DropTable(&Company{}, &Client{})
	}
	return
}
