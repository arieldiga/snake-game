package handlers

import (
	"encoding/json"
	"net/http"

	"snake-game/database"
	"snake-game/models"
)

// POST /api/score -> simpan skor baru
func SaveScoreHandler(w http.ResponseWriter, r *http.Request) {
	var s models.Score

	// Baca body request (JSON dari JavaScript) dan ubah jadi struct Go
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Data tidak valid", http.StatusBadRequest)
		return
	}

	if s.PlayerName == "" {
		s.PlayerName = "Anonim"
	}

	_, err := database.DB.Exec(
		"INSERT INTO leaderboard (player_name, score) VALUES ($1, $2)",
		s.PlayerName, s.Score,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GET /api/leaderboard -> ambil 10 skor tertinggi
func GetLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(
		"SELECT id, player_name, score, created_at FROM leaderboard ORDER BY score DESC LIMIT 10",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var scores []models.Score
	for rows.Next() {
		var s models.Score
		if err := rows.Scan(&s.ID, &s.PlayerName, &s.Score, &s.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		scores = append(scores, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scores)
}
