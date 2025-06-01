package utils

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB(models ...interface{}) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	if err := db.AutoMigrate(models...); err != nil {
		panic("failed to migrate database")
	}

	return db
}

func SeedTestData(db *gorm.DB, data ...interface{}) {
	for _, d := range data {
		if err := db.Create(d).Error; err != nil {
			panic("failed to seed test data")
		}
	}
}

func CleanupDB(db *gorm.DB, models ...interface{}) {
	for _, m := range models {
		db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(m)
	}
}
