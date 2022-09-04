package model

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
)

type TransactionFilter struct {
	UserId int
	Limit  int
	Offset int
	Order  string
}

type Transactions struct {
	TransactionList []Transaction
}

func (t *Transactions) GetItems(filter *TransactionFilter, db *gorm.DB) error {
	rows, err := db.Model(&Transaction{}).Where("user_id = ?", filter.UserId).Limit(filter.Limit).Order(filter.Order).Offset(filter.Offset).Rows()
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

		t.TransactionList = append(t.TransactionList, transaction)
	}

	if t.TransactionList == nil {
		return errors.New("у пользователя нет транзакций")
	}

	return nil
}
