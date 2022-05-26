package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type employee struct {
	ID           int        `json:"id"`
	Email        string     `json:"email"`
	Name         string     `json:"name"`
	LastModified *time.Time `json:"lastModified,omitempty"`
}

// database of employees
var employees = []employee{
	{ID: 1, Email: "a@gmail.com", Name: "A G", LastModified: getTime("2022-01-01T13:00:00Z")},
	{ID: 2, Email: "b@gmail.com", Name: "B G"},
	{ID: 3, Email: "c@gmail.com", Name: "C G"},
}

func main() {
	router := gin.Default()
	router.GET("/employees", getEmployees)
	router.GET("/employees/:id", getEmployeeByID)
	router.POST("/employees", createEmployee)
	router.DELETE("/employees/:id", deleteEmployeeByID)
	router.PUT("/employees/:id", alterEmployeeByID)

	router.Run("localhost:8080")
}

func getEmployees(c *gin.Context) {
	c.JSON(http.StatusOK, employees)
}

// returns the employee identified by the path parameter 'id'
func getEmployeeByID(c *gin.Context) {
	var idParam string = c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id should be integer"})
		return
	}

	for _, a := range employees {
		if a.ID == id {
			c.JSON(http.StatusOK, a)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "employee not found"})
}

// adds an employee from the body
func createEmployee(c *gin.Context) {
	var newEmployee employee

	if err := c.BindJSON(&newEmployee); err != nil {
		log.Fatal(err)
		c.Status(http.StatusBadRequest)
		return
	}

	var now time.Time = time.Now()
	newEmployee.LastModified = &now
	employees = append(employees, newEmployee)
	c.JSON(http.StatusCreated, newEmployee)
}

// deletes the employee identified by the path parameter 'id'
func deleteEmployeeByID(c *gin.Context) {
	var idParam string = c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id should be integer"})
		return
	}

	for idx, a := range employees {
		if a.ID == id {
			employees = append(employees[:idx], employees[idx+1:]...)
			c.Status(http.StatusOK)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "employee not found"})
}

// modify the employee identified by the path parameter 'id'
func alterEmployeeByID(c *gin.Context) {
	var idParam string = c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id should be integer"})
		return
	}
	var newEmployee employee

	if err := c.BindJSON(&newEmployee); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	for idx, a := range employees {
		if a.ID == id {
			var now time.Time = time.Now()
			newEmployee.LastModified = &now
			newEmployee.ID = id
			employees[idx] = newEmployee
			c.Status(http.StatusOK)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "employee not found"})
}

// utility to get time from a well-formated string
func getTime(rfcTime string) *time.Time {
	result, _ := time.Parse(time.RFC3339, rfcTime)
	return &result
}
