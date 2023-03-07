package testing

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"gorm.io/gorm"
)

/*
	Тест нужен был, когда решалась проблема с ограничением на горутины, когда не выдерживалось > 4 горутин
*/

type testConn struct {
	db        *gorm.DB
	num       int // кол. горутин
	c         chan User
	generator func(chan User)
	action    func(*gorm.DB, chan User, *sync.WaitGroup) error
}

const DATA_LIMIT = 1000

var timeout = 5 * time.Second

func TestConnection(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	action := func(db *gorm.DB, c chan User, wg *sync.WaitGroup) error {
		defer wg.Done()
		fmt.Println("action!")
		var u User
		for {
			select {
			case u = <-c:
				err := db.Create(&u).Error
				if err != nil {
					t.Fatalf("action error: %v", err)
					return err
				}
			case <-time.After(timeout):
				return nil
			}
		}
	}

	generator := func(c chan User) {
		for i := 0; i < DATA_LIMIT; i++ {
			alph := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
			rand.Shuffle(len(alph), func(i, j int) {
				alph[i] = alph[j]
			})
			u := User{
				ID:   rand.Int(),
				Name: string(alph[:10]),
			}
			c <- u
		}
	}

	db := initDB(t)

	teardown := setupConnection(t, db)
	defer teardown()

	// _ = setupConnection(t, db)

	tc := testConn{
		db:        db,
		num:       10,
		c:         make(chan User, DATA_LIMIT),
		action:    action,
		generator: generator,
	}
	tc.generator(tc.c)
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < tc.num; i++ {
		sdb, _ := db.DB()
		stats := sdb.Stats()
		t.Logf("Stats:\nMaxOpenConnections:%d\nOpenConnections:%d\nInUse:%d\nIdle:%d\n",
			stats.MaxOpenConnections,
			stats.OpenConnections,
			stats.InUse,
			stats.Idle,
		)
		wg.Add(1)
		go tc.action(tc.db, tc.c, &wg)
	}
	wg.Wait()
	t.Logf("duration: %s", time.Now().Sub(start))

}

func setupConnection(t *testing.T, db *gorm.DB) (teardown func()) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		t.Fatalf("unexpected error for AutoMigrate: %v", err)
	}

	teardown = func() {
		db.Migrator().DropTable(&User{})
	}
	return
}
