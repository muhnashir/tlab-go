package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database successfully.")

	// Read migration file
	content, err := ioutil.ReadFile("migrations/20260211120000_create_initial_tables.up.sql")
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}

	requests := strings.Split(string(content), ";")
	for _, request := range requests {
		request = strings.TrimSpace(request)
		if request == "" {
			continue
		}
		_, err := db.Exec(request)
		if err != nil {
			log.Printf("Failed to execute statement: %s\nError: %v\n", request, err)
		} else {
			fmt.Println("Executed statement successfully.")
		}
	}

	fmt.Println("Database setup complete.")
}
