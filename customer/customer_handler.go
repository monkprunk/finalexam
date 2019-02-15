package customer

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/mongprunk/finalexam/database"
)

var id int = 0
var customers []Customer

func getCustomersHandler(c *gin.Context) {
	pStatus := c.Query("status")
	var temp []Todo
	sql := "SELECT id, name, email, status FROM customers"
	if pStatus != "" {
		sql += " WHERE status=$1"
	}

	stmt, err := database.Conn().Prepare(sql)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	rows, err := stmt.Query()
	if pStatus != "" {
		rows, err = stmt.Query(pStatus)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	for rows.Next() {
		t := Todo{}
		err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
		if err != nil {
			log.Fatal("can't scan : ", err)
		}
		temp = append(temp, t)
	}
	c.JSON(http.StatusOK, temp)
}

func getCustomersByIdHandler(c *gin.Context) {
	pId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusProcessing, err.Error())
	}
	stmt, err := database.Conn().Prepare("SELECT id, title, status FROM customers WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	row := stmt.QueryRow(pId)
	t := Todo{}
	err = row.Scan(&t.ID, &t.Title, &t.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, t)

}

func createCustomersHandler(c *gin.Context) {
	var item Todo
	err := c.ShouldBindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	row := database.InsertTodo(item.Title, item.Status)
	var id int
	err = row.Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	item.ID = id
	c.JSON(http.StatusCreated, item)
}

func updateCustomersHandler(c *gin.Context) {
	pId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusProcessing, err.Error())
	}
	var temp Todo
	err = c.ShouldBindJSON(&temp)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	stmt, err := database.Conn().Prepare("UPDATE customers SET title=$2, status=$3 WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = stmt.Exec(pId, &temp.Title, &temp.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	temp.ID = pId
	c.JSON(http.StatusOK, temp)
}

func deleteCustomersHandler(c *gin.Context) {
	pId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusProcessing, err.Error())
	}

	stmt, err := database.Conn().Prepare("DELETE FROM customers WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	_, err = stmt.Exec(pId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("")
	v1.GET("/customers", getCustomersHandler)
	v1.GET("/customers/:id", getCustomersByIdHandler)
	v1.POST("/customers", createCustomersHandler)
	v1.PUT("/customers/:id", updateCustomersHandler)
	v1.DELETE("/customers/:id", deleteCustomersHandler)
	return r
}
