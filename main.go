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

func main() {
	router := gin.Default()
	router.GET("/students", getStudents)
	router.Run("localhost:8080")
}
