package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() (*sql.DB, error) {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")

	// connect to the default database
	connectionString := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable",
		dbUser, "postgres", dbPassword, dbHost)
	pgDB, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Failed to open a DB connection: %v", err)
		return nil, fmt.Errorf("Failed to open a DB connection: %v", err)
	}
	defer pgDB.Close()

	// check of zota_payment database exists and create it if it doesn't
	err = pgDB.QueryRow("SELECT 1 FROM pg_database WHERE datname = $1", dbName).Scan(new(int))
	if err == sql.ErrNoRows {
		_, err = pgDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			log.Fatalf("Failed to create database %s: %v", dbName, err)
			return nil, fmt.Errorf("Failed to create database %s: %v", dbName, err)
		}
		log.Printf("Database %s created", dbName)
	} else if err != nil {
		log.Fatalf("Failed to check if %s database exists: %v", dbName, err)
		return nil, fmt.Errorf("Failed to check if %s database exists: %v", dbName, err)
	}

	// connect to the zota_payment database
	connectionString = fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable",
		dbUser, dbName, dbPassword, dbHost)
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Failed to open a DB connection: %v", err)
		return nil, fmt.Errorf("Failed to open a DB connection: %v", err)
	}

	// create the callback_notifications table if it doesn't exist
	err = RunSQLScript(db, "postgres/scripts/setup.pgsql")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RunSQLScript(db *sql.DB, filename string) error {
	script, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read SQL script file: %v", err)
		return fmt.Errorf("Failed to read SQL script file: %v", err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		log.Fatalf("Failed to run SQL script: %v", err)
		return fmt.Errorf("Failed to run SQL script: %v", err)
	}

	log.Printf("SQL script %s executed successfully", filename)
	return nil
}

func CloseDB() {
	// close the DB connection
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close DB connection: %v", err)
		} else {
			log.Println("DB connection closed")
		}
	}
}
