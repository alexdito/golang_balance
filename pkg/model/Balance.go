package model

import (
	"gorm.io/gorm"
	"math"
)

type Balance struct {
	gorm.Model
	UserId   int `gorm:"uniqueIndex"`
	Balance  float64
	Currency string
}

func (b *Balance) FillBalance(db *gorm.DB) bool {
	b.Currency = "RUB"

	result := db.Where("user_id = ?", b.UserId).Find(&b)

	if result.Error != nil {
		return false
	}

	return true
}

func (b *Balance) SetBalance(db *gorm.DB) bool {
	result := db.Model(&b).Where("user_id =?", b.UserId).Update("balance", b.Balance)

	if result.RowsAffected == 0 {
		result = db.Create(&b)

		if result.Error != nil {
			return false
		}

		return true
	}

	return true
}

func (b *Balance) ConvertBalanceByCurrency() {
	switch b.Currency {
	case "USD":
		b.Balance = math.Floor((b.Balance/60.18)*100) / 100
	case "EUR":
		b.Balance = math.Floor((b.Balance/60.33)*100) / 100
	case "CYN":
		b.Balance = math.Floor((b.Balance/8.66)*100) / 100
	default:
		b.Currency = "RUB"
	}
}
