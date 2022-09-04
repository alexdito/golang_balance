package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"testovoe/pkg/model"
	"testovoe/pkg/repository"
)

type BalanceService struct {
	repository repository.Balance
}

func (b BalanceService) GetBalanceByUserId(userId int, currency string) (*model.Balance, error) {
	balanceModel, err := b.repository.GetBalanceByUserId(&model.Balance{UserId: userId})

	if err != nil {
		return balanceModel, errors.New("ошибка")
	}

	balanceModel.Currency = currency
	balanceModel.ConvertBalanceByCurrency()

	return balanceModel, nil
}

func (b BalanceService) SetBalance(userId int, amount float64) (int, map[string]any) {
	return http.StatusBadRequest, gin.H{
		"message":   "Неверный идентификатор пользователя",
		"operation": "Ошибка при установке баланса",
	}
}

func NewBalanceService(repository repository.Balance) *BalanceService {
	return &BalanceService{repository: repository}
}
