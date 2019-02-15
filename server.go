package main

import (
	"github.com/mongprunk/finalexam/customer"
	"github.com/mongprunk/finalexam/database"
)

func main() {

	database.CreateTbCustomers()
	r := customer.NewRouter()
	r.Run(":2019")
}
