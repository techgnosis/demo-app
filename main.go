package main

import (
	"net/http"
	"log"
)


func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/env", getEnv)
	http.HandleFunc("/writeMysql", writeMysql)
	http.HandleFunc("/readMysql", readMysql)
	http.HandleFunc("/redis", useRedis)
	log.Fatal(http.ListenAndServe(":8080", nil))
}