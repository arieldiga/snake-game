package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// getEnv mengambil nilai environment variable, atau pakai nilai default jika kosong.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Connect() {
	// Online: layanan seperti Neon memberi satu alamat lengkap lewat DATABASE_URL.
	// Lokal: kalau DATABASE_URL kosong, pakai variabel terpisah di bawah.
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASSWORD", ""),
			getEnv("DB_NAME", "snake_game"),
		)
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Gagal membuka koneksi database:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	// Buat tabel otomatis kalau belum ada.
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS leaderboard (
		id          SERIAL PRIMARY KEY,
		player_name VARCHAR(50) NOT NULL,
		score       INTEGER NOT NULL,
		created_at  TIMESTAMP NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		log.Fatal("Gagal membuat tabel:", err)
	}

	fmt.Println("✅ Terhubung ke PostgreSQL")
}