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

func InsertTodo(title, status string) *sql.Row {
	return Conn().QueryRow("INSERT INTO customer (name,email,status) values ($1,$2,$3) RETURNING id", title, status)
}

func CreateTbCustomer() {
	createTb := `
	CREATE TABLE IF NOT EXISTS customer (
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
