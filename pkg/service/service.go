package service

import (
	"testovoe/pkg/model"
	"testovoe/pkg/repository"
)

type Transaction interface {
	Create(userId int, amount float64, operation string, description string) (int, map[string]any)
	Transfer(senderId int, recipientId int, amount float64) (int, map[string]any)
	GetItems(userId int, limit int, offset int, order string) (int, map[string]any)
}

type Balance interface {
	GetBalanceByUserId(userId int, currency string) (*model.Balance, error)
	SetBalance(userId int, amount float64) (int, map[string]any)
}

type Service struct {
	Transaction
	Balance
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		Transaction: NewTransactionService(r.Transaction),
		Balance:     NewBalanceService(r.Balance),
	}
}
