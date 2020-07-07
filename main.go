package main

import (
	"net/http"
	"log"
)


func main() {
	http.HandleFunc("/env", getEnv)
	http.HandleFunc("/writeMysql", writeMysql)
	http.HandleFunc("/readMysql", readMysql)
	http.HandleFunc("/redis", useRedis)
	log.Fatal(http.ListenAndServe(":8080", nil))

}