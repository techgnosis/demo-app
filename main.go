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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var mysql_client *sql.DB

var redis_client *redis.Client

var postgres_client *pgx.Conn

var (
	test_writes = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "demo_app_test_writes",
			Help: "The total number of /write invocations while in test mode",
		},
	)
)

func main() {
	log.Println("App launched")

	db_type := os.Getenv("DEMO_APP_DB_TYPE")

	var write_func func(http.ResponseWriter, *http.Request)

	if db_type == "mysql" {
		log.Println("mysql mode")
		hostname := os.Getenv("DEMO_APP_MYSQL_HOSTNAME")
		database := os.Getenv("DEMO_APP_MYSQL_DATABASE")
		username := os.Getenv("DEMO_APP_MYSQL_USERNAME")
		password := os.Getenv("DEMO_APP_MYSQL_PASSWORD")
		connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", username, password, hostname, database)
		log.Println("connectionString: " + connectionString)
		var err error
		mysql_client, err = sql.Open("mysql", connectionString)
		if err != nil {
			log.Fatalf("sql.Open error: %v", err)
		}
		err = mysql_client.Ping()
		if err != nil {
			log.Fatalf("db.Ping error: %v", err)
		}

		_, err = mysql_client.Exec("CREATE TABLE IF NOT EXISTS pet ( id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY, name VARCHAR(30) NOT NULL)")
		if err != nil {
			log.Fatalf("failed to create table: %v", err)
		}
		log.Println("mysql configured")
		write_func = writeMysql

	}

	if db_type == "redis" {
		log.Println("redis mode")
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
		log.Println("redis configured")
		write_func = writeRedis
	}

	if db_type == "test" {
		log.Println("test mode")
		prometheus.MustRegister(test_writes)
		write_func = writeTestmode
	}

	if db_type == "postgres" {
		log.Println("postgres mode")
		hostname := os.Getenv("DEMO_APP_POSTGRES_HOSTNAME")
		database := os.Getenv("DEMO_APP_POSTGRES_DATABASE")
		username := os.Getenv("DEMO_APP_POSTGRES_USERNAME")
		password := os.Getenv("DEMO_APP_POSTGRES_PASSWORD")
		connectionString := fmt.Sprintf("postgresql://%s:%s@%s/%s", username, password, hostname, database)
		log.Println("connectionString: " + connectionString)
		var err error

		postgres_client, err = pgx.Connect(context.Background(), connectionString)
		if err != nil {
			log.Fatalf("failed to create postgres client: %v", err)
		}
		defer postgres_client.Close(context.Background())
		log.Println("postgres configured")
		write_func = writePostgres
	}

	http.HandleFunc("/write", write_func)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func writeMysql(w http.ResponseWriter, r *http.Request) {

	petname := "puppyface"
	sql_statement := fmt.Sprintf("insert into pet (name) values ('%s')", petname)
	_, err := mysql_client.Exec(sql_statement)
	if err != nil {
		log.Printf("error inserting row: %v", err)
	}
	log.Println("wrote to mysql")

}

func writeRedis(w http.ResponseWriter, r *http.Request) {

	petname := "puppyface"
	err := redis_client.Set("favorite-pet", petname, 0).Err()
	if err != nil {
		log.Printf("error inserting row: %v", err)
	}
	log.Println("wrote to redis")
}

func writePostgres(w http.ResponseWriter, r *http.Request) {

	var greeting string
	err := postgres_client.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		log.Printf("error inserting row: %v", err)
	}
	log.Println("wrote to postgres kinda")
}

func writeTestmode(w http.ResponseWriter, r *http.Request) {
	log.Println("test write")
	test_writes.Inc()
}
