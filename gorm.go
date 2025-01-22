package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	gorm_zerolog "github.com/wei840222/gorm-zerolog"
)

func NewGorm() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: gorm_zerolog.New(),
	})
}
