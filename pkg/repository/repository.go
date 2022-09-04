package repository

import (
	"gorm.io/gorm"
	model2 "testovoe/pkg/model"
)

type Transaction interface {
	Create(*model2.Transaction) (int, map[string]any)
	Transfer(*model2.Transaction, *model2.Transaction) (int, map[string]any)
	GetItems(*model2.TransactionFilter) (int, map[string]any)
}

type Balance interface {
	GetBalanceByUserId(balance *model2.Balance) (*model2.Balance, error)
	SetBalance(balance *model2.Balance) error
}

type Repository struct {
	Transaction
	Balance
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Transaction: NewTransactionGorm(db),
		Balance:     NewBalanceGorm(db),
	}
}
