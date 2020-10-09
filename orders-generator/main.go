package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/stdlib"
)

const createQ = `CREATE TABLE orders_info (
	orderid SERIAL NOT NULL PRIMARY KEY,
	custid INTEGER NOT NULL,
	amount INTEGER NOT NULL,
	city VARCHAR(255) NOT NULL
  );`

const dropQ = `drop TABLE orders_info`

var db *sql.DB
var cities []string

var (
	pgHost     string
	pgUser     string
	pgPassword string
	pgDB       string
)

const connStringFormat = "host=%s port=5432 user=%s password=%s dbname=%s sslmode=require"

func init() {
	pgHost = os.Getenv("PG_HOST")
	pgUser = os.Getenv("PG_USER")
	pgPassword = os.Getenv("PG_PASSWORD")
	pgDB = os.Getenv("PG_DB")
}

func main() {
	var err error

	connString := fmt.Sprintf(connStringFormat, pgHost, pgUser, pgPassword, pgDB)
	log.Println("connecting to", connString)

	db, err = sql.Open("pgx", connString)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	cities = []string{"New Delhi", "Seattle", "New York", "Austin", "Chicago", "Cleveland"}
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-exit:
			log.Println("application stopped")
			return
		default:
			insert()
			time.Sleep(3 * time.Second)
		}
	}
}

const format = "2006-01-02 15:04:05"

func insert() {
	q := "insert into retail.orders_info (custid, amount, city, purchase_time) values ($1,$2,$3,$4)"

	custid := rand.Intn(1000) + 1
	amount := rand.Intn(100) + 100
	city := cities[rand.Intn(len(cities))]

	t := time.Now().UTC().Format(format)
	_, err := db.Exec(q, custid, amount, city, t)

	if err != nil {
		log.Println("failed to insert order", err)
		return
	}
}
