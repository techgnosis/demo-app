package main

import (
	"net/http"
	"log"
)


func main() {
	http.HandleFunc("/env", getEnv)
	http.HandleFunc("/mysql", mysql)
	log.Fatal(http.ListenAndServe(":8080", nil))

}