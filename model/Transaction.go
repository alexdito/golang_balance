package model

import (
	"errors"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID      int64
	Operation   string
	Amount      float64
	Description string
}

func (t *Transaction) AddTransaction(db *gorm.DB) (float64, error) {
	tx := db.Begin()

	result := tx.Create(&t)

	if result.Error != nil {
		tx.Rollback()
		return 0, errors.New("ошибка при создании записи транзакции")
	}

	balance := Balance{UserId: t.UserID}
	balanceResult := balance.GetBalance(tx)

	switch t.Operation {
	case "withdrawal":
		balance.Balance = balance.Balance - t.Amount
	case "withdrawal-transfer":
		balance.Balance = balance.Balance - t.Amount
	case "additional":
		balance.Balance = balance.Balance + t.Amount
	case "additional-transfer":
		balance.Balance = balance.Balance + t.Amount
	}

	if !balanceResult {
		tx.Rollback()
		return 0, errors.New("ошибка при получении баланса")
	}

	editBalance := balance.SetBalance(tx)

	if !editBalance {
		tx.Rollback()
		return 0, errors.New("ошибка при изменении баланса")
	}

	tx.Commit()

	return balance.Balance, nil
}
