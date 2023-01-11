package testing

import (
	"reflect"
	"testing"

	"gorm.io/gorm"
)

type test struct {
	got  any
	want any
}

func TestSelect(t *testing.T) {
	db := initDB(t)
	teardown := setupSelect(t, db)
	defer teardown()
	// 1	Alice
	// 2	Bob
	// 3	Martha
	// 4	Sam
	// 5	John
	// 6	Tammy
	t.Run("basic", func(t *testing.T) {
		got := make([]User, 0)
		err := db.Model(&User{}).Select(`"id", "name"`).Where(`"id" > 3`).Find(&got).Error
		if err != nil {
			t.Fatalf("unexpected error for Find: %v", err)
		}
		// for _, u := range got {
		// 	fmt.Printf("User %d: %s\n", u.ID, u.Name)
		// }
		for _, tc := range []test{
			{
				got: got,
				want: []User{
					{
						ID:   4,
						Name: "Sam",
					},
					{
						ID:   5,
						Name: "John",
					},
					{
						ID:   6,
						Name: "Tammy",
					},
				},
			},
		} {
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got and want are not the same")
			}
		}
	})

}

func setupSelect(t *testing.T, db *gorm.DB) (teardown func()) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		t.Fatalf("unexpected error for AutoMigrate: %v", err)
	}

	us := []User{
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
			Name: "Martha",
		},
		{
			ID:   4,
			Name: "Sam",
		},
		{
			ID:   5,
			Name: "John",
		},
		{
			ID:   6,
			Name: "Tammy",
		},
	}

	err = db.Create(us).Error
	if err != nil {
		t.Fatalf("unexpected error for Create: %v", err)
	}

	teardown = func() {
		db.Migrator().DropTable(&User{})
	}
	return
}
