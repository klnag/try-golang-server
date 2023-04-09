package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func getTodos(c *gin.Context) {
	db, err := sql.Open("mysql", "root:g4VxIX11W5JADW9ACLpU@tcp(containers-us-west-5.railway.app:5974)/railway")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM todo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}

	c.JSON(200, todos)
}

func createTodo(c *gin.Context) {
	db, err := sql.Open("mysql", "root:g4VxIX11W5JADW9ACLpU@tcp(containers-us-west-5.railway.app:5974)/railway")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO todo (title, completed) VALUES (?, ?)", todo.Title, todo.Completed)
	if err != nil {
		log.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	todo.ID = int(id)
	c.JSON(200, todo)
}

func updateTodo(c *gin.Context) {
	db, err := sql.Open("mysql", "root:g4VxIX11W5JADW9ACLpU@tcp(containers-us-west-5.railway.app:5974)/railway")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	id := c.Param("id")

	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err = db.Exec("UPDATE todo SET title=?, completed=? WHERE id=?", todo.Title, todo.Completed, id)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, todo)
}

func deleteTodo(c *gin.Context) {
	db, err := sql.Open("mysql", "root:g4VxIX11W5JADW9ACLpU@tcp(containers-us-west-5.railway.app:5974)/railway")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	id := c.Param("id")

	_, err = db.Exec("DELETE FROM todo WHERE id=?", id)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("Todo with ID %s has been deleted", id)})
}

func main() {
	r := gin.Default()

	r.GET("/todos", getTodos)
	r.POST("/todos", createTodo)
	r.PUT("/todos/:id", updateTodo)
	r.DELETE("/todos/:id", deleteTodo)

	r.Run()
}
