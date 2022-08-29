package model

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
)

const (
	Additional         = "additional"
	Withdrawal         = "withdrawal"
	Transfer           = "transfer"
	AdditionalTransfer = "additional-transfer"
	WithdrawalTransfer = "withdrawal-transfer"
)

type Transaction struct {
	gorm.Model
	UserID      int64
	Operation   string
	Amount      float64
	Description string
}

type Transactions struct {
	Transaction []Transaction
}

func (t *Transactions) GetTransactionList(userId int64, limit int, offset int, order string, db *gorm.DB) error {
	rows, err := db.Model(&Transaction{}).Where("user_id = ?", userId).Limit(limit).Order(order).Offset(offset).Rows()
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	if err != nil {
		return err
	}

	for rows.Next() {
		var transaction Transaction
		err := db.ScanRows(rows, &transaction)

		if err != nil {
			return err
		}

		t.Transaction = append(t.Transaction, transaction)
	}

	return nil
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

	if !balanceResult {
		tx.Rollback()
		return 0, errors.New("ошибка при получении баланса")
	}

	switch t.Operation {
	case Withdrawal:
		balance.Balance = balance.Balance - t.Amount
	case WithdrawalTransfer:
		balance.Balance = balance.Balance - t.Amount
	case Additional:
		balance.Balance = balance.Balance + t.Amount
	case AdditionalTransfer:
		balance.Balance = balance.Balance + t.Amount
	}

	editBalance := balance.SetBalance(tx)

	if !editBalance {
		tx.Rollback()
		return 0, errors.New("ошибка при изменении баланса")
	}

	tx.Commit()

	return balance.Balance, nil
}
