package repository

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"sync"
	model2 "testovoe/pkg/model"
)

type TransactionGorm struct {
	db    *gorm.DB
	mutex sync.Mutex
}

func NewTransactionGorm(db *gorm.DB) *TransactionGorm {
	return &TransactionGorm{db: db}
}

func (r *TransactionGorm) Create(transaction *model2.Transaction) (int, map[string]any) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	db := r.db.Begin()

	err := transaction.Create(db)

	if err != nil {
		db.Rollback()

		return http.StatusBadRequest, gin.H{
			"operation": transaction.Description,
			"error":     err.Error(),
		}
	}

	balance, err := changeBalance(transaction, db)

	if err != nil {
		db.Rollback()

		return http.StatusBadRequest, gin.H{
			"operation": transaction.Description,
			"error":     err.Error(),
		}
	}

	db.Commit()

	return http.StatusOK, gin.H{
		"balance":     balance,
		"operation":   transaction.Operation,
		"description": transaction.Description,
	}
}

func (r *TransactionGorm) Transfer(transactionSender *model2.Transaction, transactionRecipient *model2.Transaction) (int, map[string]any) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	db := r.db.Begin()

	err := transactionSender.Create(db)

	if err != nil {
		db.Rollback()

		return http.StatusBadRequest, gin.H{
			"operation": "Перевод средств",
			"error":     err.Error(),
		}
	}

	_, err = changeBalance(transactionSender, db)

	if err != nil {
		db.Rollback()

		return http.StatusBadRequest, gin.H{
			"operation": "Перевод средств",
			"error":     err.Error(),
		}
	}

	err = transactionRecipient.Create(db)

	if err != nil {
		db.Rollback()

		return http.StatusBadRequest, gin.H{
			"operation": "Перевод средств",
			"error":     err.Error(),
		}
	}

	_, err = changeBalance(transactionRecipient, db)

	if err != nil {
		db.Rollback()

		return http.StatusBadRequest, gin.H{
			"operation": "Перевод средств",
			"error":     err.Error(),
		}
	}

	db.Commit()

	return http.StatusOK, gin.H{
		"operation":   "Перевод средств",
		"description": fmt.Sprintf("Перевод средств от %d для %d успешно выполнен", transactionSender.UserID, transactionRecipient.UserID),
	}
}

func (r *TransactionGorm) GetItems(filter *model2.TransactionFilter) (int, map[string]any) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	transactions := model2.Transactions{}
	err := transactions.GetItems(filter, r.db)

	if err != nil {
		return http.StatusBadRequest, gin.H{
			"error":     err.Error(),
			"operation": "Получение списка транзакций",
		}
	}

	return http.StatusOK, gin.H{
		"transactions": transactions.TransactionList,
		"operation":    "Получение списка транзакций",
	}
}

func changeBalance(transaction *model2.Transaction, db *gorm.DB) (float64, error) {
	balanceModel := model2.Balance{UserId: transaction.UserID}

	if !balanceModel.FillBalance(db) {
		return 0, errors.New("ошибка при получении баланса")
	}

	switch transaction.Operation {
	case "additional":
		balanceModel.Balance += transaction.Amount
	case "withdrawal":
		balanceModel.Balance -= transaction.Amount
	default:
		return 0, errors.New("неизвестная операция")
	}

	if balanceModel.Balance < 0 {
		return 0, errors.New("недостаточно средств на балансе")
	}

	if balanceModel.SetBalance(db) {
		return balanceModel.Balance, nil
	}

	return 0, errors.New("ошибка при изменении баланса")
}
