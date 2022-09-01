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
	balance := Balance{UserId: t.UserID}

	if !balance.GetBalance(db) {
		return 0, errors.New("ошибка при создании баланса")
	}

	switch t.Operation {
	case Additional:
		balance.Balance += t.Amount
	case Withdrawal:
		balance.Balance -= t.Amount
		if balance.Balance < 0 {
			return 0, errors.New("сумма списания больше чем баласн")
		}
	default:
		return 0, errors.New("неизвестная операция")
	}

	result := db.Create(&t)

	if result.Error != nil {
		return 0, result.Error
	}

	editBalance := balance.SetBalance(db)

	if !editBalance {
		return 0, errors.New("ошибка при изменении баланса")
	}

	return balance.Balance, nil
}

func GetBalanceByUserId(id int64, db *gorm.DB) (float64, error) {
	balanceModel := Balance{UserId: id}

	result := db.Where("user_id = ?", balanceModel.UserId).First(&balanceModel)

	if result.Error != nil {
		return 0, errors.New("нет записи о балансе")
	}

	return balanceModel.Balance, nil
}
