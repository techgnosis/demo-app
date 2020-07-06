package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/buger/jsonparser"
    "fmt"
    "log"
    "net/http"
    "os"
)

func mysql(w http.ResponseWriter, r *http.Request) {
	vcapservices := []byte(os.Getenv("VCAP_SERVICES"))
	hostname, err := jsonparser.GetString(vcapservices, "p-mysql", "[0]", "credentials", "hostname")
	database, err := jsonparser.GetString(vcapservices, "p-mysql", "[0]", "credentials", "name")
	username, err := jsonparser.GetString(vcapservices, "p-mysql", "[0]", "credentials", "username")
	password, err := jsonparser.GetString(vcapservices, "p-mysql", "[0]", "credentials", "password")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", username, password, hostname, database)
	db, err := sql.Open("mysql", connectionString)
	err = db.Ping()
	if err != nil {
	  panic(err.Error()) // proper error handling instead of panic in your app
	}
	
	rows, err := db.Query("select * from pet")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var (
		name string
		owner string
		species string
	)
	for rows.Next() {
		err := rows.Scan(&name, &owner, &species)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, name + " " + owner + " " + species + "\n")

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	
}