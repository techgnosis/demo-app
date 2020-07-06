package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/go-redis/redis/v7"
	"net/http"
	"os"
)

func useRedis(w http.ResponseWriter, r *http.Request) {
	vcapservices := []byte(os.Getenv("VCAP_SERVICES"))
	host, err := jsonparser.GetString(vcapservices, "p-redis", "[0]", "credentials", "host")
	port, err := jsonparser.GetString(vcapservices, "p-redis", "[0]", "credentials", "port")
	password, err := jsonparser.GetString(vcapservices, "p-redis", "[0]", "credentials", "password")


    client := redis.NewClient(&redis.Options{
	  Addr:    host +  ":" + port,
	  Password: password,
      DB:       0,
	})

	err = client.Set("pas-test-key", "5", 0).Err()
	if err != nil {
		panic(err)
	}

	_, err = client.Get("pas-test-key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("derp")
	fmt.Fprintf(w, "Howdy y'all!")
}
