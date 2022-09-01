package operations

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"net/http"
	"strconv"
	"strings"
	"testovoe/app"
	"testovoe/model"
)

const (
	GettingCurrentBalanceMessage    = "Получение текущего баланса пользователя"
	InvalidUserIdMessage            = "Неверный идентификатор пользователя"
	WriteOffCashMessage             = "Списания средств с баланса"
	AdditionalCashMessage           = "Пополнение баланса"
	WrongTransactionAmountMessage   = "Неверная сумма транзакции"
	TransactionListOperationMessage = "Получение списка транзакций"
	TransferOperationMessage        = "Перевод средств"
	TransferOperationToMessage      = "Перевод средств для %d"
	TransferOperationFromMessage    = "Перевод средств от %d"
	WrongOperationMessage           = "Неизвестная операция"
)

var operationDescription = map[string]string{
	Additional: AdditionalCashMessage,
	Withdrawal: WriteOffCashMessage,
	Transfer:   TransferOperationMessage,
}

const (
	Additional = "additional"
	Withdrawal = "withdrawal"
	Transfer   = "transfer"
)

func GetTransactions(c *gin.Context, app *app.App) (int, map[string]any) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	userId, err := strconv.ParseInt(c.Query("userId"), 0, 64)

	if err != nil {
		return http.StatusBadRequest, gin.H{
			"message":   InvalidUserIdMessage,
			"operation": TransactionListOperationMessage,
		}
	}

	limit, err := strconv.ParseInt(c.Query("limit"), 0, 64)
	if err != nil {
		limit = 2
	}

	offset, err := strconv.ParseInt(c.Query("offset"), 0, 64)
	if err != nil {
		offset = 2
	}

	order := c.Query("order")

	if order == "" || order == "created_at" || order == "amount" {
		transactions := model.Transactions{}
		err = transactions.GetTransactionList(userId, int(limit), int(offset), order, app.DataBase)

		if err != nil {
			return http.StatusBadRequest, gin.H{
				"message":   InvalidUserIdMessage,
				"operation": TransactionListOperationMessage,
			}
		}

		return http.StatusOK, gin.H{
			"balance":   transactions.Transaction,
			"operation": TransactionListOperationMessage,
		}
	}

	return http.StatusBadRequest, gin.H{
		"message":   "Неизвестное поле для сортировки",
		"operation": TransactionListOperationMessage,
	}
}

func GetBalance(c *gin.Context, app *app.App) (int, map[string]any) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	userId, err := strconv.ParseInt(c.Query("userId"), 0, 64)

	if err != nil {
		return http.StatusBadRequest, gin.H{
			"message":   InvalidUserIdMessage,
			"operation": GettingCurrentBalanceMessage,
		}
	}

	balance, err := model.GetBalanceByUserId(userId, app.DataBase)

	if err != nil {
		return http.StatusBadRequest, gin.H{
			"message":   InvalidUserIdMessage,
			"operation": GettingCurrentBalanceMessage,
		}
	}

	currency := strings.ToUpper(c.Query("currency"))

	balance, currency = convertBalanceByCurrency(balance, currency)

	return http.StatusOK, gin.H{
		"balance":   balance,
		"currency":  currency,
		"operation": GettingCurrentBalanceMessage,
	}
}

func convertBalanceByCurrency(balance float64, currency string) (float64, string) {
	switch currency {
	case "USD":
		return math.Floor((balance/60.18)*100) / 100, currency
	case "EUR":
		return math.Floor((balance/60.33)*100) / 100, currency
	case "CYN":
		return math.Floor((balance/8.66)*100) / 100, currency
	default:
		return balance, "RUB"
	}
}

func ExecuteTransaction(c *gin.Context, app *app.App, operation string) (int, map[string]any) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	db := app.DataBase

	description, isExists := operationDescription[operation]

	if !isExists {
		return http.StatusBadRequest, gin.H{
			"message": WrongOperationMessage,
		}
	}

	tx := db.Begin()

	var httpStatus int
	var response map[string]any
	var err error

	switch operation {
	case Transfer:
		senderId, err := strconv.ParseInt(c.PostForm("sender"), 0, 64)

		if err != nil {
			return http.StatusBadRequest, gin.H{
				"message":   InvalidUserIdMessage,
				"operation": TransferOperationMessage,
			}
		}

		recipientId, err := strconv.ParseInt(c.PostForm("recipient"), 0, 64)

		if err != nil {
			return http.StatusBadRequest, gin.H{
				"message":   InvalidUserIdMessage,
				"operation": TransferOperationMessage,
			}
		}

		amount, err := strconv.ParseFloat(c.PostForm("amount"), 64)

		if amount <= 0 {
			return http.StatusBadRequest, gin.H{
				"message":   WrongTransactionAmountMessage,
				"operation": TransferOperationMessage,
			}
		}

		if err != nil {
			return http.StatusBadRequest, gin.H{
				"message":   WrongTransactionAmountMessage,
				"operation": TransferOperationMessage,
			}
		}

		httpStatus, response, err = addTransaction(senderId, amount, Withdrawal, fmt.Sprintf(TransferOperationToMessage, recipientId), tx)

		if err != nil {
			tx.Rollback()

			return http.StatusBadRequest, gin.H{
				"message":   WrongTransactionAmountMessage,
				"operation": TransferOperationMessage,
			}
		}

		httpStatus, response, err = addTransaction(recipientId, amount, Additional, fmt.Sprintf(TransferOperationFromMessage, senderId), tx)

		if err != nil {
			tx.Rollback()

			return http.StatusBadRequest, gin.H{
				"message":   WrongTransactionAmountMessage,
				"operation": TransferOperationMessage,
			}
		}

		httpStatus = http.StatusOK
		response = gin.H{
			"message":   fmt.Sprintf("Перевод средств от %d для %d успешно выполнен", senderId, recipientId),
			"operation": TransferOperationMessage,
		}
	case Additional, Withdrawal:
		userId, err := strconv.ParseInt(c.PostForm("userId"), 0, 64)

		if err != nil {
			return http.StatusBadRequest, gin.H{
				"message":   InvalidUserIdMessage,
				"operation": description,
			}
		}

		amount, err := strconv.ParseFloat(c.PostForm("amount"), 64)

		if amount <= 0 {
			return http.StatusBadRequest, gin.H{
				"message":   WrongTransactionAmountMessage,
				"operation": TransferOperationMessage,
			}
		}

		if err != nil {
			return http.StatusBadRequest, gin.H{
				"message":   WrongTransactionAmountMessage,
				"operation": description,
			}
		}

		httpStatus, response, err = addTransaction(userId, amount, operation, description, tx)
	default:
		return http.StatusBadRequest, gin.H{
			"message": WrongOperationMessage,
		}
	}

	if err != nil {
		tx.Rollback()
		return http.StatusBadRequest, gin.H{
			"message":   err,
			"operation": operation,
		}
	}

	tx.Commit()

	return httpStatus, response
}

func addTransaction(userId int64, amount float64, operation string, description string, db *gorm.DB) (int, map[string]any, error) {
	transaction := model.Transaction{
		UserID:      userId,
		Amount:      amount,
		Operation:   operation,
		Description: description,
	}

	balance, err := transaction.AddTransaction(db)

	if err != nil {
		return http.StatusBadRequest, gin.H{
			"message":   err.Error(),
			"operation": description,
		}, err
	}

	return http.StatusOK, gin.H{
		"balance":   balance,
		"operation": description,
	}, nil
}
