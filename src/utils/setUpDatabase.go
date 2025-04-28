package utils

import (
    "database/sql"
    "fmt"
    "os"
    _ "github.com/lib/pq"
)

func SetupDatabase() (*sql.DB, error) {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        return nil, fmt.Errorf("DATABASE_URL is not set in the environment")
    }
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        return nil, fmt.Errorf("error opening database connection: %v", err)
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(10)
    db.SetConnMaxLifetime(300)
    
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error connecting to database: %v", err)
    }
    return db, nil
}
