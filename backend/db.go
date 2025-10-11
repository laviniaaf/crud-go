package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() *sql.DB {
	
// environment variables defined in docker-compose
	dbUser := getEnvOrDefault("DB_USER", "usuario")
	dbPass := getEnvOrDefault("DB_PASS", "senha123")
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "3306")
	dbName := getEnvOrDefault("DB_NAME", "database")

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
	CREATE TABLE IF NOT EXISTS items  (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price DECIMAL(10,2) NOT NULL
	)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
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
