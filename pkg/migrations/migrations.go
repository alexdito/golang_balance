package migrations

import (
	"gorm.io/gorm"
	model2 "testovoe/pkg/model"
)

func CreateTablesIfNotExist(db *gorm.DB) {
	if !db.Migrator().HasTable(&model2.Transaction{}) {
		err := db.Migrator().CreateTable(&model2.Transaction{})

		if err != nil {
			return
		}
	}

	if !db.Migrator().HasTable(&model2.Balance{}) {
		err := db.Migrator().CreateTable(&model2.Balance{})

		if err != nil {
			return
		}
	}
}
