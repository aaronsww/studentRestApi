package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// "errors"
)

type student struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var students = []student{
	{ID: "1", Name: "Jeevan"},
	{ID: "2", Name: "Aaron"},
	{ID: "3", Name: "Joseph"},
}

func getStudents(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, students)
}

func createStudent(c *gin.Context) {
	var newStudent student

	if err := c.BindJSON(&newStudent); err != nil {
		return
	}

	students = append(students, newStudent)
	c.IndentedJSON(http.StatusCreated, newStudent)
}

func main() {
	router := gin.Default()
	router.GET("/students", getStudents)
	router.POST("/students", createStudent)
	router.Run("localhost:8080")
}
