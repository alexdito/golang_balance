package operations

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"testovoe/app"
	"testovoe/model"
)

func GetBalance(c *gin.Context, app *app.App) (int, map[string]any) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	userId, err := strconv.ParseInt(c.Query("userId"), 0, 64)

	if err != nil {
		return http.StatusBadRequest, gin.H{
			"message":   "Неверный идентификатор пользователя",
			"operation": "Получение текущего баланса пользователя.",
		}
	}

	balance := model.Balance{UserId: userId}

	result := balance.GetBalance(app.DataBase)

	if !result {
		return http.StatusBadRequest, gin.H{
			"message":   "Пользователя с таким ID не существует",
			"operation": "Получение текущего баланса пользователя.",
		}
	}

	return http.StatusOK, gin.H{
		"balance":   balance.Balance,
		"operation": "Получение текущего баланса пользователя.",
	}
}

func ChangeBalance(c *gin.Context, app *app.App, operation string) (int, map[string]any) {
	app.Mutex.Lock()
	defer app.Mutex.Unlock()

	switch operation {
	case "additional":
		return additionalBalance(c, app)
	case "withdrawal":
		return withdrawalBalance(c, app)
	case "transfer":
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
			"message":   "Неверный идентификатор пользователя",
			"operation": "Списания средств с баланса.",
		}
	}

	amount, errAmount := strconv.ParseFloat(c.PostForm("amount"), 64)

	if errAmount != nil {
		return http.StatusBadRequest, gin.H{
			"message":   "Неверная сумма списания.",
			"operation": "Списания средств с баланса.",
		}
	}

	balanceModel := model.Balance{UserId: userId}

	result := balanceModel.GetBalance(app.DataBase)

	if !result {
		return http.StatusBadRequest, gin.H{
			"message":   "Пользователя с таким ID не существует",
			"operation": "Списания средств с баланса.",
		}
	}

	if balanceModel.Balance-amount < 0 {
		return http.StatusBadRequest, gin.H{
			"message":   "Недостаточно средств для списания",
			"operation": "Списания средств с баланса.",
		}
	}

	transaction := model.Transaction{
		UserID:      userId,
		Amount:      amount,
		Operation:   "withdrawal",
		Description: "Списания средств с баланса.",
	}

	balance, err := transaction.AddTransaction(app.DataBase)

	if err != nil {
		return http.StatusBadRequest, gin.H{
			"message":   err.Error(),
			"operation": "Списания средств с баланса.",
		}
	}

	return http.StatusOK, gin.H{
		"balance":   balance,
		"operation": "Списания средств с баланса.",
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
		Operation:   "withdrawal-transfer",
		Description: fmt.Sprintf("Перевод средств для %d", recipientId),
	}

	transactionRecipient := model.Transaction{
		UserID:      recipientId,
		Amount:      amount,
		Operation:   "additional-transfer",
		Description: fmt.Sprintf("Перевод средств от %d", senderId),
	}

	_, errTransactionSender := transactionSender.AddTransaction(app.DataBase)

	fmt.Println(errTransactionSender)
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
			"message":   "Неверный идентификатор пользователя",
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
		Operation:   "additional",
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
