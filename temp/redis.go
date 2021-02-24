package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/go-redis/redis/v7"
	"net/http"
	"os"
)

func getRedisClient() *redis.Client {
	vcapservices := []byte(os.Getenv("VCAP_SERVICES"))
	host, _ := jsonparser.GetString(vcapservices, "p.redis", "[0]", "credentials", "host")
	port, _ := jsonparser.GetInt(vcapservices, "p.redis", "[0]", "credentials", "port")
	password, _ := jsonparser.GetString(vcapservices, "p.redis", "[0]", "credentials", "password")
	fmt.Println("host " + host)
	fmt.Printf("port %d\n", port)
	fmt.Println("password " + password)
	return redis.NewClient(&redis.Options{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       0,
	})
}

func writeRedis(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if err := r.ParseForm(); err != nil {
            fmt.Fprintf(w, "ParseForm() err: %v", err)
            return
        }
        petname := r.FormValue("favoritepet")
		client := getRedisClient()
		err := client.Set("favorite-pet", petname, 0).Err()
		if err != nil {
			panic(err)
		}
	}
}

func readRedis(w http.ResponseWriter, r *http.Request) {
	client := getRedisClient()
	value, err := client.Get("favorite-pet").Result()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, value)
}
