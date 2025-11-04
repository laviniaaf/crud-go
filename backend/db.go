package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func InitDB() *sql.DB {
	// load .env
	_ = godotenv.Load()

	// environment variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPass == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatalf("Missing DB configuration: set DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME in environment or .env")
	}

	// connection
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	var err error

	db, err = sql.Open("mysql", dsn)

	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	for i := 0; i < 10; i++ {
		if err := db.Ping(); err != nil {
			log.Printf("Attempt %d: Error connecting to database: %v. Trying again in 5 seconds...", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("Connected to the bank:", dbName)
		break
	}

	query := `
	CREATE TABLE IF NOT EXISTS bills (
		id BINARY(16) PRIMARY KEY,
		embasa DECIMAL(10,2) NOT NULL,
		coelba DECIMAL(10,2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`
	_, err = db.Exec(query)

	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	log.Println("Table items created and verified successfully!!!!")

	return db
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
