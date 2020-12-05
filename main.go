package main

import (
	"fmt"
	"log"
	"net/http"
)


func main() {
	fmt.Println("App launched")
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/env", getEnv)
	http.HandleFunc("/writeMysql", writeMysql)
	http.HandleFunc("/readMysql", readMysql)
	http.HandleFunc("/writeRedis", writeRedis)
	http.HandleFunc("/readRedis", readRedis)
	log.Fatal(http.ListenAndServe(":8080", nil))
}