package database

import (
	"database/sql"
	"log"
	"os"
)

var db *sql.DB

func Conn() *sql.DB {
	if db != nil {
		return db
	}
	var err error
	url := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatal("can't connect database : ", err)
	}
	return db
}

func InsertCustomer(name, email, status string) (int, error) {

	row := Conn().QueryRow("INSERT INTO customers (name,email,status) values ($1,$2,$3) RETURNING id", name, email, status)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func SelectAllCustomers() (*sql.Rows, error) {

	sql := "SELECT id, name, email, status FROM customers"

	stmt, err := Conn().Prepare(sql)
	if err != nil {
		return nil, err
	}

	return stmt.Query()
}

func SelectCustomersById(pId int) (*sql.Row, error) {
	stmt, err := Conn().Prepare("SELECT id, name, email, status FROM customers WHERE id=$1")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(pId)

	return row, nil
}

func UpdateCustomer(pId int, pName, pEmail, pStatus string) (int, error) {
	stmt, err := Conn().Prepare("UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1;")
	if err != nil {
		return 0, err
	}

	_, err = stmt.Exec(pId, pName, pEmail, pStatus)
	if err != nil {
		return 0, err
	}

	return pId, nil
}

func DeleteCustomer(pId int) error {
	stmt, err := Conn().Prepare("DELETE FROM customers WHERE id=$1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(pId)
	if err != nil {
		return err
	}

	return nil
}

func CreateTbCustomers() {
	createTb := `
	CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);
	`
	_, err := Conn().Exec(createTb)
	if err != nil {
		log.Fatal("can't create table : ", err)
	}
}
