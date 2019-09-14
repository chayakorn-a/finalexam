package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

var db *sql.DB

func authMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "token2019" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}

func main() {
	db = GetSqlDB()

	CreateCustomerTable()

	r := gin.Default()
	r.Use(authMiddleware)

	r.POST("/customers", CreateCustomerHandler)    // InsertCustomerTable
	r.GET("/customers/:id", getIdCustomerHandler)  // QueryCustomerTable
	r.GET("/customers", GetCustomerHandler)        // QueryAllCustomerTable
	r.PUT("/customers/:id", updateCustomerHandler) // UpdateCustomerTable
	r.DELETE("/customers/:id", DeleteCustomerHandler)

	r.Run(":2019")
	defer db.Close()
}

func GetCustomerHandler(c *gin.Context) {
	c.JSON(http.StatusOK, QueryAllCustomerTable())
}

func QueryAllCustomerTable() []Customer {
	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers")
	if err != nil {
		log.Fatal("can't prepare query all todos statement", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("can't query all todos", err)
	}
	var custs []Customer
	for rows.Next() {
		var cust Customer
		err = rows.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
		if err != nil {
			log.Fatal("can't scan row into variable", err)
		}
		custs = append(custs, cust)
	}
	fmt.Println("query all todos success")
	return custs
}

func CreateCustomerHandler(c *gin.Context) {
	var cust Customer

	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	InsertCustomerTable(&cust)
	c.JSON(http.StatusCreated, cust)
}

func InsertCustomerTable(cust *Customer) {
	row := db.QueryRow("INSERT INTO customers (name, email, status) values ($1, $2, $3) RETURNING id", cust.Name, cust.Email, cust.Status)
	var id int
	err := row.Scan(&id)
	if err != nil {
		fmt.Println("can't scan id", err)
		return
	}
	cust.ID = id
	fmt.Println("insert todo scccess id : ", id)
}

func getIdCustomerHandler(c *gin.Context) {
	idd := c.Param("id")
	id, _ := strconv.Atoi(idd)
	cust := QueryCustomerTable(id)
	c.JSON(http.StatusOK, cust)
}

func QueryCustomerTable(rowId int) Customer {
	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers where id=$1")
	if err != nil {
		log.Fatal("can't prepare query one row statement", err)
	}
	row := stmt.QueryRow(rowId)
	var cust Customer
	err = row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
	if err != nil {
		log.Fatal("can't Scan row into variables", err)
	}
	return cust
}

func updateCustomerHandler(c *gin.Context) {
	var cust Customer
	idd := c.Param("id")

	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, _ := strconv.Atoi(idd)
	UpdateCustomerTable(id, cust)
	c.JSON(http.StatusOK, cust)
}

func UpdateCustomerTable(id int, cust Customer) {
	stmt, err := db.Prepare("UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1;")
	if err != nil {
		log.Fatal("can't prepare statement update", err)
	}

	if _, err := stmt.Exec(id, cust.Name, cust.Email, cust.Status); err != nil {
		log.Fatal("error execute update", err)
	}
	fmt.Println("update success")
}

func DeleteCustomerHandler(c *gin.Context) {
	idd := c.Param("id")

	id, _ := strconv.Atoi(idd)
	DeleteCustomerTable(id)
	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}

func DeleteCustomerTable(id int) {
	stmt, err := db.Prepare("DELETE FROM customers WHERE id=$1;")

	if err != nil {
		log.Fatal("can't prepare statement update", err)
	}

	if _, err := stmt.Exec(id); err != nil {
		log.Fatal("error execute update", err)
	}
	fmt.Println("update success")
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func CreateCustomerTable() {
	createTb := `
	CREATE TABLE IF NOT EXISTS customers (
	id SERIAL PRIMARY KEY,
	name TEXT,
	email TEXT,
	status TEXT
	);
	`
	_, err := db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table ", err)
	} else {
		fmt.Println("create table success")
	}
}

func GetSqlDB() *sql.DB {
	url := os.Getenv("DATABASE_URL")
	fmt.Println(url)
	db1, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	return db1
}
