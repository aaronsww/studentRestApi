package main

import (
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
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

func studentById(c *gin.Context) {
	id := c.Param("id")
	student, err := getStudentById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.h{"message": "Book not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, student)
}

func getStudentById(id string) (*student, error) {
	for i, s := range students {
		if s.ID == id {
			return &students[i], nil
		}
	}

	return nil, errors.New("student not found")
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
	router.GET("/students/:id", studentById)
	router.POST("/students", createStudent)
	router.Run("localhost:8080")
}
