package main

import (
	"os"

	"github.com/joho/godotenv"

	"net/http"

	"errors"

	"github.com/gin-gonic/gin"

	"database/sql"

	"fmt"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "testdbee"
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
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Student not found."})
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

func createStudent(c *gin.Context, db *sql.DB) {
	var newStudent student

	if err := c.BindJSON(&newStudent); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert data into the database
	result, err := db.Exec("INSERT INTO students (name) VALUES ($1)", newStudent.Name)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, result)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	dbPassword := os.Getenv("DB_PASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, dbPassword, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	router := gin.Default()
	router.GET("/students", getStudents)
	router.GET("/students/:id", studentById)
	router.POST("/students", func(c *gin.Context) {
		createStudent(c, db)
	})
	router.Run("localhost:8080")
}
