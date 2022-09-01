package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"testovoe/app"
	"testovoe/migrations"
	"testovoe/operations"
)

const (
	Additional = "additional"
	Withdrawal = "withdrawal"
	Transfer   = "transfer"
)

func main() {
	application, err := app.GetApp()

	if err != nil {
		fmt.Println("Ошибка подключения к БД!")
	}

	migrations.CreateTablesIfNotExist(application.DataBase)

	r := gin.Default()

	//Пополнение
	r.POST("/additional", func(c *gin.Context) {
		httpStatus, H := operations.ExecuteTransaction(c, application, Additional)
		c.JSON(httpStatus, H)
	})

	//Списание
	r.POST("/withdrawal", func(c *gin.Context) {
		httpStatus, H := operations.ExecuteTransaction(c, application, Withdrawal)
		c.JSON(httpStatus, H)
	})

	//Перевод
	r.POST("/transfer", func(c *gin.Context) {
		httpStatus, H := operations.ExecuteTransaction(c, application, Transfer)
		c.JSON(httpStatus, H)
	})

	//Баланс
	r.GET("/balance", func(c *gin.Context) {
		httpStatus, H := operations.GetBalance(c, application)
		c.JSON(httpStatus, H)
	})

	//Список транзакций
	r.GET("/transactions", func(c *gin.Context) {
		httpStatus, H := operations.GetTransactions(c, application)
		c.JSON(httpStatus, H)
	})

	r.Run()
}
