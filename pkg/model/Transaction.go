package model

import (
	"errors"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID      int
	Operation   string
	Amount      float64
	Description string
}

func (t Transaction) Create(db *gorm.DB) error {
	result := db.Create(&t)

	if result.Error != nil {
		return errors.New("ошибка при создании записи транзакции")
	}

	return nil
}
