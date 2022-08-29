package operations

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"testovoe/app"
	"testovoe/model"
)

const (
	GettingCurrentBalanceMessage    = "Получение текущего баланса пользователя."
	InvalidUserIdMessage            = "Неверный идентификатор пользователя."
	WriteOffCashMessage             = "Списания средств с баланса."
	WrongWriteOffAmountMessage      = "Неверная сумма списания."
	DeficiencyCashOnBalanceMessage  = "Недостаточно средств для списания."
	TransactionListMessage          = "Недостаточно средств для списания."
	TransactionListOperationMessage = "Получение списка транзакций."
	//"Неверная сумма перевода."
	//"Перевод средств."
	//"Пользователя с таким ID не существует"
	//"Недостаточно средств для перевода"

)

const (
	Additional         = "additional"
	Withdrawal         = "withdrawal"
	Transfer           = "transfer"
	AdditionalTransfer = "additional-transfer"
	WithdrawalTransfer = "withdrawal-transfer"
)

func GetTransactions(c *gin.Context, app *app.App) (int, map[string]any) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	userId, errUserId := strconv.ParseInt(c.PostForm("userId"), 0, 64)

	if errUserId != nil {
		return http.StatusBadRequest, gin.H{
			"message":   InvalidUserIdMessage,
			"operation": TransactionListOperationMessage,
		}
	}

	limit, errLimit := strconv.ParseInt(c.PostForm("limit"), 0, 64)
	if errLimit != nil {
		limit = 2
	}

	offset, errOffset := strconv.ParseInt(c.PostForm("offset"), 0, 64)
	if errOffset != nil {
		offset = 2
	}

	order := c.PostForm("order")

	transactions := model.Transactions{}
	err := transactions.GetTransactionList(userId, int(limit), int(offset), order, app.DataBase)

	fmt.Println(transactions)

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

	return http.StatusOK, gin.H{
		"balance":   model.GetBalanceByUserId(userId, app.DataBase),
		"operation": GettingCurrentBalanceMessage,
	}
}

func ChangeBalance(c *gin.Context, app *app.App, operation string) (int, map[string]any) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	switch operation {
	case Additional:
		return additionalBalance(c, app)
	case Withdrawal:
		return withdrawalBalance(c, app)
	case Transfer:
		return transfer(c, app)
	default:
		return http.StatusBadRequest, gin.H{
			"message": "Неизвестная операция",
		}
	}
}

func withdrawalBalance(c *gin.Context, app *app.App) (int, map[string]any) {
	userId, errUserId := strconv.ParseInt(c.PostForm("userId"), 0, 64)

	if errUserId != nil {
		return http.StatusBadRequest, gin.H{
			"message":   InvalidUserIdMessage,
			"operation": WriteOffCashMessage,
		}
	}

	amount, errAmount := strconv.ParseFloat(c.PostForm("amount"), 64)

	if errAmount != nil {
		return http.StatusBadRequest, gin.H{
			"message":   WrongWriteOffAmountMessage,
			"operation": WriteOffCashMessage,
		}
	}

	balance := model.GetBalanceByUserId(userId, app.DataBase)

	if balance-amount < 0 {
		return http.StatusBadRequest, gin.H{
			"message":   DeficiencyCashOnBalanceMessage,
			"operation": WriteOffCashMessage,
		}
	}

	transaction := model.Transaction{
		UserID:      userId,
		Amount:      amount,
		Operation:   Withdrawal,
		Description: WriteOffCashMessage,
	}

	newBalance, err := transaction.AddTransaction(app.DataBase)

	if err != nil {
		return http.StatusBadRequest, gin.H{
			"message":   err.Error(),
			"operation": WriteOffCashMessage,
		}
	}

	return http.StatusOK, gin.H{
		"balance":   newBalance,
		"operation": WriteOffCashMessage,
	}
}

func transfer(c *gin.Context, app *app.App) (int, map[string]any) {
	senderId, errSender := strconv.ParseInt(c.PostForm("sender"), 0, 64)

	if errSender != nil {
		return http.StatusBadRequest, gin.H{
			"message":   "Неверный идентификатор отправителя.",
			"operation": "Перевод средств.",
		}
	}

	recipientId, errRecipient := strconv.ParseInt(c.PostForm("recipient"), 0, 64)

	if errRecipient != nil {
		return http.StatusBadRequest, gin.H{
			"message":   "Неверный идентификатор получателяю",
			"operation": "Перевод средств.",
		}
	}

	amount, errAmount := strconv.ParseFloat(c.PostForm("amount"), 64)

	if errAmount != nil {
		return http.StatusBadRequest, gin.H{
			"message":   "Неверная сумма перевода.",
			"operation": "Перевод средств.",
		}
	}

	balanceSender := model.Balance{UserId: senderId}
	result := balanceSender.GetBalance(app.DataBase)

	if !result {
		return http.StatusBadRequest, gin.H{
			"message":   "Пользователя с таким ID не существует",
			"operation": "Перевод средств.",
		}
	}

	if balanceSender.Balance-amount < 0 {
		return http.StatusBadRequest, gin.H{
			"message":   "Недостаточно средств для перевода",
			"operation": "Перевод средств.",
		}
	}

	transactionSender := model.Transaction{
		UserID:      senderId,
		Amount:      amount,
		Operation:   WithdrawalTransfer,
		Description: fmt.Sprintf("Перевод средств для %d", recipientId),
	}

	transactionRecipient := model.Transaction{
		UserID:      recipientId,
		Amount:      amount,
		Operation:   AdditionalTransfer,
		Description: fmt.Sprintf("Перевод средств от %d", senderId),
	}

	_, errTransactionSender := transactionSender.AddTransaction(app.DataBase)

	if errTransactionSender != nil {
		return http.StatusBadRequest, gin.H{
			"message":   "Ошибка при переводе средств",
			"operation": "Перевод средств.",
		}
	}
	_, errTransactionRecipient := transactionRecipient.AddTransaction(app.DataBase)

	if errTransactionRecipient != nil {
		return http.StatusBadRequest, gin.H{
			"message":   "Ошибка при переводе средств",
			"operation": "Перевод средств.",
		}
	}

	return http.StatusOK, gin.H{
		"operation": "Перевод средств.",
	}
}

func additionalBalance(c *gin.Context, app *app.App) (int, map[string]any) {
	userId, errUserId := strconv.ParseInt(c.PostForm("userId"), 0, 64)

	if errUserId != nil {
		return http.StatusBadRequest, gin.H{
			"message":   InvalidUserIdMessage,
			"operation": "Пополнение баланса.",
		}
	}

	amount, errAmount := strconv.ParseFloat(c.PostForm("amount"), 64)

	if errAmount != nil {
		return http.StatusBadRequest, gin.H{
			"message":   "Неверная сумма зачисления.",
			"operation": "Пополнение баланса.",
		}
	}

	transaction := model.Transaction{
		UserID:      userId,
		Amount:      amount,
		Operation:   Additional,
		Description: "Пополнение баланса.",
	}

	balance, err := transaction.AddTransaction(app.DataBase)

	if err != nil {
		return http.StatusBadRequest, gin.H{
			"message":   err.Error(),
			"operation": "Пополнение баланса.",
		}
	}

	return http.StatusOK, gin.H{
		"balance":   balance,
		"operation": "Пополнение баланса.",
	}
}
