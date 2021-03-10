package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v4"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// DB clients
var mysql_client *sql.DB
var redis_client *redis.Client
var postgres_client *pgx.Conn

// Prometheus metrics
var db_writes = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "demoapp_writes",
		Help: "The total number of /write invocations",
	})

var db_write_times = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "demoapp_writes_responsetime",
		Help: "Latency for the /write endpoint",
	})

// Other globals
var curl_response string

func main() {
	log.Println("App launched!!!!")

	prometheus.MustRegister(db_writes)
	prometheus.MustRegister(db_write_times)

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

	timer := prometheus.NewTimer(prometheus.ObserverFunc(db_write_times.Set))
	defer timer.ObserveDuration()

	// 100 + rand(50)
	random_number := rand.Intn(50)
	time.Sleep(time.Duration(100+random_number) * time.Millisecond)

	petname := "puppyface"
	sql_statement := fmt.Sprintf("insert into pet (name) values ('%s')", petname)
	_, err := mysql_client.Exec(sql_statement)
	if err != nil {
		log.Printf("error inserting row: %v", err)
	}
	db_writes.Inc()
	message := "mysql write success\n"
	log.Println(message)
	io.WriteString(w, message)

}

func writeRedis(w http.ResponseWriter, r *http.Request) {

	petname := "puppyface"
	err := redis_client.Set("favorite-pet", petname, 0).Err()
	if err != nil {
		log.Printf("error inserting row: %v", err)
	}
	db_writes.Inc()
	log.Println("wrote to redis")
}

func writePostgres(w http.ResponseWriter, r *http.Request) {

	var greeting string
	err := postgres_client.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		log.Printf("error inserting row: %v", err)
	}
	db_writes.Inc()
	message := "wrote to postgres kinda\n"
	log.Println(message)
	io.WriteString(w, message)
}

func writeTestmode(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(db_write_times.Set))
	defer timer.ObserveDuration()

	db_writes.Inc()
	message := "test write success\n"
	log.Println(message)

	// 100 + rand(50)
	random_number := rand.Intn(50)
	time.Sleep(time.Duration(100+random_number) * time.Millisecond)

	io.WriteString(w, message)
}
