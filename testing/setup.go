package testing

import (
	"os"
	"testing"

	"github.com/aeroideaservices/gormtnt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDB(t *testing.T) *gorm.DB {
	dsn, _ := os.LookupEnv("TEST_DB_DSN")
	cnf := &gorm.Config{}
	cnf.Logger = logger.Default.LogMode(logger.Info)
	db, err := gorm.Open(gormtnt.Open(dsn), cnf)
	if err != nil {
		t.Fatalf("can't init db: %v", err)
	}
	return db
}
