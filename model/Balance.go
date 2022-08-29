package model

import (
	"gorm.io/gorm"
)

type Balance struct {
	gorm.Model
	UserId  int64
	Balance float64
}

func GetBalanceByUserId(id int64, db *gorm.DB) float64 {
	balanceModel := Balance{UserId: id}

	result := db.Where("user_id = ?", balanceModel.UserId).First(&balanceModel)

	if result.Error != nil {
		return 0
	}

	return balanceModel.Balance
}

func (b *Balance) GetBalance(db *gorm.DB) bool {
	result := db.Where("user_id = ?", b.UserId).First(&b)

	if result.Error != nil {
		return b.createBalanceIfNotExists(db)
	}

	return true
}

func (b *Balance) SetBalance(db *gorm.DB) bool {
	result := db.Model(&b).Where("user_id =?", b.UserId).Update("balance", b.Balance)

	if result.Error != nil {
		return false
	}

	return true
}

func (b *Balance) createBalanceIfNotExists(db *gorm.DB) bool {
	b.Balance = 0
	result := db.Create(&b)

	if result.Error != nil {
		return false
	}

	return true
}
