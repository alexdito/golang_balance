package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"testovoe/pkg/app"
	"testovoe/pkg/handler"
	"testovoe/pkg/migrations"
	"testovoe/pkg/repository"
	"testovoe/pkg/service"
)

const (
	Additional = "additional"
	Withdrawal = "withdrawal"
)

func main() {
	application, err := app.GetApp()

	if err != nil {
		fmt.Println("Ошибка подключения к БД!")
	}

	migrations.CreateTablesIfNotExist(application.DataBase)

	r := gin.Default()

	repositories := repository.NewRepository(application.DataBase)
	transactionService := service.NewService(repositories)
	handlers := handler.NewHandler(transactionService)

	//Баланс
	r.GET("/balance", func(c *gin.Context) {
		c.JSON(handlers.GetBalance(c))
	})

	//Пополнение
	r.POST("/additional", func(c *gin.Context) {
		c.JSON(handlers.CreateTransaction(c, Additional))
	})

	//Списание
	r.POST("/withdrawal", func(c *gin.Context) {
		c.JSON(handlers.CreateTransaction(c, Withdrawal))
	})

	//Перевод
	r.POST("/transfer", func(c *gin.Context) {
		c.JSON(handlers.Transfer(c))
	})

	//Список транзакций
	r.GET("/transactions", func(c *gin.Context) {
		c.JSON(handlers.GetItems(c))
	})

	r.Run()
}
