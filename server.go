package main

import (
	"github.com/mongprunk/finalexam/customer"
	"github.com/mongprunk/finalexam/database"
)

func main() {

	database.Conn()
	customer.CreateTb()
	r := customer.NewRouter()
	r.Run(":1234")
}
