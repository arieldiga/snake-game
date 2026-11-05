package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"snake-game/database"
	"snake-game/models"
)

const (
	maxNameLength = 20
	maxScore      = 100000 // batas wajar, tolak skor yang mustahil
)

// POST /api/score -> simpan skor baru
func SaveScoreHandler(w http.ResponseWriter, r *http.Request) {
	// Batasi ukuran body maksimal 1 KB supaya tidak bisa dibanjiri data besar.
	r.Body = http.MaxBytesReader(w, r.Body, 1024)

	var s models.Score
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Data tidak valid", http.StatusBadRequest)
		return
	}

	// Validasi di server: jangan percaya data dari browser.
	s.PlayerName = strings.TrimSpace(s.PlayerName)
	if s.PlayerName == "" {
		s.PlayerName = "Anonim"
	}
	if utf8.RuneCountInString(s.PlayerName) > maxNameLength {
		http.Error(w, "Nama terlalu panjang (maks 20 karakter)", http.StatusBadRequest)
		return
	}
	if s.Score < 0 || s.Score > maxScore {
		http.Error(w, "Skor tidak valid", http.StatusBadRequest)
		return
	}

	// Context dengan timeout: query dibatalkan kalau database macet > 5 detik.
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := database.DB.ExecContext(ctx,
		"INSERT INTO leaderboard (player_name, score) VALUES ($1, $2)",
		s.PlayerName, s.Score,
	)
	if err != nil {
		// Detail error dicatat di log server saja, bukan dikirim ke pengguna.
		log.Println("Gagal simpan skor:", err)
		http.Error(w, "Terjadi kesalahan server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GET /api/leaderboard -> ambil 10 skor tertinggi
func GetLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := database.DB.QueryContext(ctx,
		"SELECT id, player_name, score, created_at FROM leaderboard ORDER BY score DESC, created_at ASC LIMIT 10",
	)
	if err != nil {
		log.Println("Gagal ambil leaderboard:", err)
		http.Error(w, "Terjadi kesalahan server", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// make(...) memastikan hasil kosong jadi [] bukan null di JSON.
	scores := make([]models.Score, 0, 10)
	for rows.Next() {
		var s models.Score
		if err := rows.Scan(&s.ID, &s.PlayerName, &s.Score, &s.CreatedAt); err != nil {
			log.Println("Gagal baca baris:", err)
			http.Error(w, "Terjadi kesalahan server", http.StatusInternalServerError)
			return
		}
		scores = append(scores, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scores)
}