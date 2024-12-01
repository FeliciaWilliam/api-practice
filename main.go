package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	host     = "localhost"
	port     = 8080
	user     = ""
	password = ""
	db_name  = ""
)

type task struct {
	ID          int
	title       string
	description string
	status      string
	priority    string
	due_date    string
	created_at  string
	updated_at  string
}

var tasks = []task{}
var nextID = 1

func readTasks(context *gin.Context, db *sql.DB) {
	context.IndentedJSON(http.StatusOK, tasks)
	query:=``
}

func readTasksbyID(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for _, task := range tasks {
		if task.ID == id {
			context.JSON(http.StatusOK, task)
			return
		}
	}

	context.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func createTasks(context *gin.Context) {
	var Task task
	if err := context.ShouldBindJSON(&Task); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	Task.ID = nextID
	nextID++
	tasks = append(tasks, Task)

	context.JSON(http.StatusCreated, Task)
}

func updateTasks(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedTask task
	if err := context.ShouldBindJSON(&updatedTask); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			updatedTask.ID = id
			tasks[i] = updatedTask
			context.JSON(http.StatusOK, updatedTask)
			return
		}
	}

	context.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func deleteTasks(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			context.Status(http.StatusNoContent)
			return
		}
	}

	context.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/mydb?sslmode=disable")
	if err != nil {
		fmt.Println("Something went wrong.")
		return
	}

	defer db.Close()

	router := gin.Default()
	router.GET("/tasks", func(c *gin.Context){
		readTasks(c,db)
	 })
	router.POST("/tasks", createTasks)
	router.PUT("/tasks", updateTasks)
	router.DELETE("/tasks", deleteTasks)

	router.Run("localhost:8080")
}
