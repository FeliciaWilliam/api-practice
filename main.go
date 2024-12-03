package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	host     = "localhost"
	port     = 8080
	user     = "postgres"
	password = ""
	db_name  = "api_practice"
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

var db *pgxpool.Pool
var tasks = []task{}
var nextID = 1

func readTasks(context *gin.Context, db *sql.DB) {
	/*context.IndentedJSON(http.StatusOK, tasks)
	query:=*/
	rows, err := db.Query("SELECT id, title, description, status FROM tasks")
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	defer rows.Close()

	var tasks []task
	for rows.Next() {
		var task task
		if err := rows.Scan(&task.ID, &task.title, &task.description, &task.status); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse task"})
			return
		}
		tasks = append(tasks, task)
	}

	context.JSON(http.StatusOK, tasks)

}

func readTasksbyID(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	/*if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for _, task := range tasks {
		if task.ID == id {
			context.JSON(http.StatusOK, task)
			return
		}
	}*/
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var task task
	err = db.QueryRow(context.Background(),
		"SELECT id, title, description, status FROM tasks WHERE id=$1", id).Scan(&task.ID, &task.title, &task.description, &task.status)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	context.JSON(http.StatusOK, task)

	context.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func createTasks(context *gin.Context) {
	var task task
	if err := context.ShouldBindJSON(&task); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	/*Task.ID = nextID
	nextID++
	tasks = append(tasks, Task)*/

	err := db.QueryRow(context.Background(),
		"INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id",
		task.title, task.description, task.status).Scan(&task.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	context.JSON(http.StatusCreated, task)
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

	/*for i, task := range tasks {
		if task.ID == id {
			updatedTask.ID = id
			tasks[i] = updatedTask
			context.JSON(http.StatusOK, updatedTask)
			return
		}
	}

	context.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})*/

	result, err := db.Exec(context.Background(),
		"UPDATE tasks SET title=$1, description=$2, status=$3 WHERE id=$4",
		updatedTask.title, updatedTask.description, updatedTask.status, id)
	if err != nil || result.RowsAffected() == 0 {
		context.JSON(http.StatusNotFound, gin.H{"error": "Task not found or update failed"})
		return
	}

	updatedTask.ID = id
	context.JSON(http.StatusOK, updatedTask)
}

func deleteTasks(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	/*for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			context.Status(http.StatusNoContent)
			return
		}
	}

	context.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})*/

	result, err := db.Exec(context.Background(), "DELETE FROM tasks WHERE id=$1", id)
	if err != nil || result.RowsAffected() == 0 {
		context.JSON(http.StatusNotFound, gin.H{"error": "Task not found or delete failed"})
		return
	}

	context.Status(http.StatusNoContent)
}

func main() {
	//db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/mydb?sslmode=disable")
	var err error
	pgxpool.New(context.Background(), "postgres://postgres:reimuhakurei@localhost:5432/postgres")
	if err != nil {
		fmt.Println("Something went wrong.")
		return
	}

	defer db.Close()

	router := gin.Default()
	router.GET("/tasks", func(c *gin.Context) {
		readTasks(c, db)
	})
	router.POST("/tasks", createTasks)
	router.PUT("/tasks", updateTasks)
	router.DELETE("/tasks", deleteTasks)

	router.Run("localhost:5432")
}
