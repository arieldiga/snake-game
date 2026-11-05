package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	connStr := "host=localhost port=5432 user=postgres password=aadd515 dbname=snake_game sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Gagal membuka koneksi database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	fmt.Println("✅ Terhubung ke PostgreSQL")
}
