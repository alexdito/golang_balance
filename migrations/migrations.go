package migrations

import (
	"gorm.io/gorm"
	"testovoe/model"
)

func CreateTablesIfNotExist(db *gorm.DB) {
	if !db.Migrator().HasTable(&model.Transaction{}) {
		err := db.Migrator().CreateTable(&model.Transaction{})

		if err != nil {
			return
		}
	}

	if !db.Migrator().HasTable(&model.Balance{}) {
		err := db.Migrator().CreateTable(&model.Balance{})

		if err != nil {
			return
		}
	}
}
