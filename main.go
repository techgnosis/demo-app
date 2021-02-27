package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v4"
)

var db *sql.DB

var redis_client *redis.Client

var postgres_client *pgx.Conn

func main() {
	fmt.Println("App launched")

	db_type := os.Getenv("DEMO_APP_DB_TYPE")

	if db_type == "mysql" {
		fmt.Println("mysql being configured")
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
			panic(err.Error())
		}
		err = db.Ping()
		if err != nil {
			fmt.Println("db.Ping error")
			panic(err.Error())
		}

		_, err = db.Exec("CREATE TABLE IF NOT EXISTS pet ( id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY, name VARCHAR(30) NOT NULL)")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("mysql configured")
		http.HandleFunc("/write", writeMysql)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

	if db_type == "redis" {
		fmt.Println("redis being configured")
		hostname := os.Getenv("DEMO_APP_REDIS_HOSTNAME")
		port := os.Getenv("DEMO_APP_REDIS_PORT")
		password := os.Getenv("DEMO_APP_REDIS_PASSWORD")
		tls_config := &tls.Config{
			InsecureSkipVerify: true,
		}
		redis_client = redis.NewClient(&redis.Options{
			Addr:      fmt.Sprintf("%s:%s", hostname, port),
			Password:  password,
			DB:        0,
			TLSConfig: tls_config,
		})
		fmt.Println("redis configured")
		http.HandleFunc("/write", writeRedis)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

	if db_type == "postgres" {
		fmt.Println("postgres being configured")
		hostname := os.Getenv("DEMO_APP_POSTGRES_HOSTNAME")
		database := os.Getenv("DEMO_APP_POSTGRES_DATABASE")
		username := os.Getenv("DEMO_APP_POSTGRES_USERNAME")
		password := os.Getenv("DEMO_APP_POSTGRES_PASSWORD")
		connectionString := fmt.Sprintf("postgresql://%s:%s@%s/%s", username, password, hostname, database)
		fmt.Println("connectionString: " + connectionString)
		var err error

		postgres_client, err = pgx.Connect(context.Background(), connectionString)

		if err != nil {
			log.Fatal(err)
		}
		defer postgres_client.Close(context.Background())
		http.HandleFunc("/write", writePostgres)
		log.Fatal(http.ListenAndServe(":8080", nil))

	}

}

func writeMysql(w http.ResponseWriter, r *http.Request) {
	fmt.Println("mysql write entered")

	petname := "puppyface"

	fmt.Println("About to write a row")
	sql_statement := fmt.Sprintf("insert into pet (name) values ('%s')", petname)
	_, err := db.Exec(sql_statement)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote a row")

}

func writeRedis(w http.ResponseWriter, r *http.Request) {
	fmt.Println("redis write entered")
	petname := "puppyface"
	fmt.Println("about to write to redis")
	err := redis_client.Set("favorite-pet", petname, 0).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("wrote to redis")
}

func writePostgres(w http.ResponseWriter, r *http.Request) {
	fmt.Println("postgres write entered")
	var greeting string

	err := postgres_client.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		panic(err)
	}
	fmt.Println(greeting)
	fmt.Println("wrote to postgres kinda")
}
