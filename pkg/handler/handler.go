package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"testovoe/pkg/service"
)

type Handler struct {
	services *service.Service
}

func (h *Handler) CreateTransaction(c *gin.Context, operation string) (int, map[string]any) {
	var operationDescription = map[string]string{
		"additional": "Пополнение баланса",
		"withdrawal": "Списания средств с баланса",
	}

	description, isExists := operationDescription[operation]

	if !isExists {
		return http.StatusBadRequest, gin.H{
			"message":     "Неизвестная операция",
			"operation":   operation,
			"description": description,
		}
	}

	userId, err := strconv.Atoi(c.PostForm("userId"))

	if err != nil {
		return http.StatusBadRequest, gin.H{
			"message":     "Неверный идентификатор пользователя",
			"operation":   operation,
			"description": description,
		}
	}

	amount, err := strconv.ParseFloat(c.PostForm("amount"), 64)

	if amount <= 0 {
		return http.StatusBadRequest, gin.H{
			"message":     "Неверная сумма транзакции",
			"operation":   operation,
			"description": description,
		}
	}

	return h.services.Create(userId, amount, operation, description)
}

func (h *Handler) Transfer(c *gin.Context) (int, map[string]any) {
	senderId, err := strconv.Atoi(c.PostForm("sender"))

	if err != nil {
		return http.StatusBadRequest, gin.H{
			"message":   "Неверный идентификатор пользователя",
			"operation": "Перевод средств",
		}
	}

	recipientId, err := strconv.Atoi(c.PostForm("recipient"))

	if err != nil {
		return http.StatusBadRequest, gin.H{
			"message":   "Неверный идентификатор пользователя",
			"operation": "Перевод средств",
		}
	}

	amount, err := strconv.ParseFloat(c.PostForm("amount"), 64)

	if amount <= 0 {
		return http.StatusBadRequest, gin.H{
			"message":   "Неверная сумма транзакции",
			"operation": "Перевод средств",
		}
	}

	if err != nil {
		return http.StatusBadRequest, gin.H{
			"message":   "Неверная сумма транзакции",
			"operation": "Перевод средств",
		}
	}

	return h.services.Transfer(senderId, recipientId, amount)
}

func (h *Handler) GetBalance(c *gin.Context) (int, map[string]any) {
	userId, err := strconv.Atoi(c.Query("userId"))

	if err != nil {
		return http.StatusBadRequest, gin.H{
			"message":   "Неверный идентификатор пользователя",
			"operation": "Получение текущего баланса пользователя",
		}
	}

	currency := strings.ToUpper(c.Query("currency"))

	balance, err := h.services.GetBalanceByUserId(userId, currency)

	if err != nil {
		return http.StatusBadRequest, gin.H{
			"message":   "Ошибка при получении текущего баланса пользователя",
			"operation": "Получение текущего баланса пользователя",
		}
	}

	return http.StatusBadRequest, gin.H{
		"balance":   balance.Balance,
		"currency":  balance.Currency,
		"operation": "Получение текущего баланса пользователя",
	}
}

func (h *Handler) GetItems(c *gin.Context) (int, map[string]any) {
	userId, err := strconv.Atoi(c.Query("userId"))

	if err != nil {
		return http.StatusBadRequest, gin.H{
			"message":   "Неверный идентификатор пользователя",
			"operation": "Получение списка транзакций",
		}
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 2
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 2
	}

	order := c.Query("order")

	if order != "" && order != "created_at" && order != "amount" {
		return http.StatusBadRequest, gin.H{
			"message":   "Неверный идентификатор пользователя",
			"operation": "Получение списка транзакций",
		}
	}

	return h.services.GetItems(userId, limit, offset, order)
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
