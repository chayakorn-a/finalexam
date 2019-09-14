package main

import (
	"https://github.com/chayakorn-a/finalexam/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID     string `json:id`
	Title  string `json:title`
	Status string `json:status`
}

// On git Bash
// DATABASE_URL=postgres://ednlnmoi:huVNUQvl1bR9puqxPOrGkJ-pPiVQ23hY@otto.db.elephantsql.com:5432/ednlnmoi go run main.go
func main() {
	/*url := os.Getenv("DATABASE_URL")
	fmt.Println(url)
	db, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()*/

	//CreateTodoTable(db)
	//InsertTodoTable(db)
	//QueryTodoTable(db)
	//UpdateTodoTable(db)
	// QueryTodoTable(db)
	//QueryAllTodoTable(db)

	r := gin.Default()

	r.GET("/todos", GetTodoHandler)       // QueryAllTodoTable
	r.GET("/todos/:id", getIdTodoHandler) // QueryTodoTable
	r.POST("/todos", CreateTodoHandler)   // InsertTodoTable
	r.DELETE("/todos/:id", DeleteTodoHandler)
	r.PUT("/todos/:id", updateTodoHandler) // UpdateTodoTable

	r.Run(":1234")
}

var todos []Todo
var ids int = 0

// r.GET("/todos", GetTodoHandler)
// Retrieve all records to show only
func GetTodoHandler(c *gin.Context) {
	c.JSON(http.StatusOK, database.QueryAllTodoTable())
}

// r.POST("/todos",CreateTodoHandler)
// Get a record then insert to table with the new ID
func CreateTodoHandler(c *gin.Context) {
	var todo Todo

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ids++
	todo.ID = strconv.Itoa(ids)
	todos = append(todos, todo)
	//c.JSON(200, gin.H{"message": "pong post"})
	c.JSON(200, map[string]string{"message": "pong post"})
}

// r.GET("/todos/:id", getIdTodoHandler)
// Lookup in table to find with id
func getIdTodoHandler(c *gin.Context) {
	idd := c.Param("id")

	for _, n := range todos {
		if n.ID == idd {
			c.JSON(http.StatusOK, n)
			return
		}
	}
}

// r.DELETE("/todos/:id",DeleteTodoHandler)
// delete spcific ID from table
func DeleteTodoHandler(c *gin.Context) {
	idd := c.Param("id")
	num := Find(todos, idd)
	remove(todos, num)
	todos = todos[:(len(todos) - 1)]
	c.JSON(http.StatusOK, todos)
}

// r.PUT("/todos/:id",updateTodoHandler)
// Input id with JSON object
// Update existing record
func updateTodoHandler(c *gin.Context) {
	var todo Todo
	idd := c.Param("id")

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	num := Find(todos, idd)
	todos[num].Title = todo.Title
	todos[num].Status = todo.Status

	c.JSON(http.StatusOK, todos[num])
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//
func Find(a []Todo, id string) int {
	for i, n := range a {
		if id == n.ID {
			return i
		}
	}
	return len(a)
}

func remove(slice []Todo, s int) []Todo {
	return append(slice[:s], slice[s+1:]...)
}
