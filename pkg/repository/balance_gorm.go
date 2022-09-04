package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"testovoe/pkg/model"
)

type BalanceGorm struct {
	db    *gorm.DB
	mutex sync.Mutex
}

func NewBalanceGorm(db *gorm.DB) *BalanceGorm {
	return &BalanceGorm{db: db}
}

func (b *BalanceGorm) GetBalanceByUserId(balanceModel *model.Balance) (*model.Balance, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if !balanceModel.FillBalance(b.db) {
		return balanceModel, errors.New("ошибка при получении баланса пользователя")
	}

	fmt.Println(balanceModel)

	return balanceModel, nil
}

func (b *BalanceGorm) SetBalance(balanceModel *model.Balance) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	balanceModel.SetBalance(b.db)
	return nil
}
