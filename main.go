package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"testovoe/app"
	"testovoe/migrations"
	"testovoe/operations"
)

func main() {
	application, err := app.GetApp()

	if err != nil {
		fmt.Println("Ошибка подключения к БД!")
	}

	migrations.CreateTablesIfNotExist(application.DataBase)

	r := gin.Default()

	// Пополнение
	r.POST("/additional", func(c *gin.Context) {
		httpStatus, H := operations.ChangeBalance(c, application, "additional")

		c.JSON(httpStatus, H)
	})

	//Списание
	r.POST("/withdrawal", func(c *gin.Context) {
		httpStatus, H := operations.ChangeBalance(c, application, "withdrawal")

		c.JSON(httpStatus, H)
	})

	//Перевод
	r.POST("/transfer", func(c *gin.Context) {
		httpStatus, H := operations.ChangeBalance(c, application, "transfer")

		c.JSON(httpStatus, H)
	})

	//Баланс
	r.GET("/balance", func(c *gin.Context) {
		httpStatus, H := operations.GetBalance(c, application)

		c.JSON(httpStatus, H)
	})

	r.Run()
}
