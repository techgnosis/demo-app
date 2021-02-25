package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	fmt.Println("App launched")

	hostname := os.Getenv("DEMO_APP_MYSQL_HOSTNAME")
	database := os.Getenv("DEMO_APP_MYSQL_DATABASE")
	username := os.Getenv("DEMO_APP_MYSQL_USERNAME")
	password := os.Getenv("DEMO_APP_MYSQL_PASSWORD")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", username, password, hostname, database)
	fmt.Println("connectionString: " + connectionString)
	var err error
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("sql.Open error")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	http.HandleFunc("/write", write)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func write(w http.ResponseWriter, r *http.Request) {
	fmt.Println("write entered")

	petname := "puppyface"

	err := db.Ping()
	if err != nil {
		fmt.Println("db.Ping error")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS pet ( id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY, name VARCHAR(30) NOT NULL)")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("About to write a row")
	sql_statement := fmt.Sprintf("insert into pet (name) values ('%s')", petname)
	_, err = db.Exec(sql_statement)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote a row")

}
