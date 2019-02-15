package customer

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/mongprunk/finalexam/database"
	"github.com/mongprunk/finalexam/middleware"
)

var id int = 0
var customers []Customer

func getCustomersHandler(c *gin.Context) {
	var result []Customer
	rows, err := database.SelectAllCustomers()
	if err != nil {
		c.JSON(http.StatusProcessing, err.Error())
	}
	for rows.Next() {
		t := Customer{}
		err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		result = append(result, t)
	}
	c.JSON(http.StatusOK, result)
}

func getCustomersByIdHandler(c *gin.Context) {
	pId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusProcessing, err.Error())
	}
	var row *sql.Row
	row, err = database.SelectCustomersById(pId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	result := Customer{}
	err = row.Scan(&result.ID, &result.Name, &result.Email, &result.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)

}

func createCustomersHandler(c *gin.Context) {
	var item Customer
	err := c.ShouldBindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	item.ID, err = database.InsertCustomer(item.Name, item.Email, item.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, item)
}

func updateCustomersHandler(c *gin.Context) {
	pId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusProcessing, err.Error())
		return
	}
	var temp Customer
	err = c.ShouldBindJSON(&temp)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	temp.ID, err = database.UpdateCustomer(pId, temp.Name, temp.Email, temp.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, temp)
}

func deleteCustomersHandler(c *gin.Context) {
	pId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusProcessing, err.Error())
		return
	}

	err = database.DeleteCustomer(pId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.LoginMiddleware)
	r.GET("/customers", getCustomersHandler)
	r.GET("/customers/:id", getCustomersByIdHandler)
	r.POST("/customers", createCustomersHandler)
	r.PUT("/customers/:id", updateCustomersHandler)
	r.DELETE("/customers/:id", deleteCustomersHandler)
	return r
}
