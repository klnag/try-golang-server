package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

type Todos struct {
	ID     int    `json:"id"`
	Title  string `json:"title" validate:"required"`
	IsDone bool   `json:"isDone"`
}

var todoss []Todos = []Todos{
	{ID: 0, Title: "run", IsDone: false},
	{ID: 1, Title: "clean", IsDone: false},
}

func allTodoss(c *gin.Context) {
	c.JSON(200, gin.H{"message": todoss})
}

func createTodos(c *gin.Context) {
	var todos Todos
	validate := validator.New()

	if err := c.BindJSON(&todos); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if err := validate.Struct(todos); err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	todoss = append(todoss, todos)
	c.JSON(400, gin.H{"msg": todos})
}

func mains() {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/todos")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()

	r.GET("/api/todo", allTodoss)
	r.POST("/api/todo", createTodo)

	r.Run()
}
