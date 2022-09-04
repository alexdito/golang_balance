package service

import (
	"fmt"
	model2 "testovoe/pkg/model"
	"testovoe/pkg/repository"
)

type TransactionService struct {
	repository repository.Transaction
}

func (t TransactionService) Create(userId int, amount float64, operation string, description string) (int, map[string]any) {
	return t.repository.Create(&model2.Transaction{
		UserID:      userId,
		Amount:      amount,
		Operation:   operation,
		Description: description,
	})
}

func (t TransactionService) Transfer(senderId int, recipientId int, amount float64) (int, map[string]any) {
	senderTransaction := model2.Transaction{
		UserID:      senderId,
		Amount:      amount,
		Operation:   "withdrawal",
		Description: fmt.Sprintf("Перевод средств для %d", recipientId),
	}

	recipientTransaction := model2.Transaction{
		UserID:      recipientId,
		Amount:      amount,
		Operation:   "additional",
		Description: fmt.Sprintf("Перевод средств от %d", senderId),
	}

	result, arr := t.repository.Transfer(&senderTransaction, &recipientTransaction)

	return result, arr
}

func (t TransactionService) GetItems(userId int, limit int, offset int, order string) (int, map[string]any) {
	return t.repository.GetItems(&model2.TransactionFilter{
		UserId: userId,
		Limit:  limit,
		Offset: offset,
		Order:  order,
	})
}

func NewTransactionService(repository repository.Transaction) *TransactionService {
	return &TransactionService{repository: repository}
}
