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

func getStudents(c *gin.Context, db *sql.DB) {
	// c.IndentedJSON(http.StatusOK, students)
	// Query all students from the database
	rows, err := db.Query("SELECT id, name FROM students")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var students []student
	for rows.Next() {
		var s student
		err := rows.Scan(&s.ID, &s.Name)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		students = append(students, s)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, students)
}

func studentById(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	student, err := getStudentById(id, db)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Student not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, student)
}

func getStudentById(id string, db *sql.DB) (*student, error) {
	row := db.QueryRow("SELECT id, name FROM students WHERE id = $1", id)

	var s student
	err := row.Scan(&s.ID, &s.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("student not found")
		}
		return nil, err
	}

	return &s, nil
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
	router.GET("/students", func(c *gin.Context) {
		getStudents(c, db)
	})
	router.GET("/students/:id", func(c *gin.Context) {
		studentById(c, db)
	})
	router.POST("/addStudent", func(c *gin.Context) {
		createStudent(c, db)
	})
	router.Run("localhost:8080")
}
