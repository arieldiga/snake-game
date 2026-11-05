package main

import (
	"fmt"
	"log"
	"net/http"

	"snake-game/database"
	"snake-game/handlers"
)

func main() {
	database.Connect()
	defer database.DB.Close()

	mux := http.NewServeMux()

	// Sajikan file statis (HTML, CSS, JS) dari folder static/
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	// Endpoint API untuk leaderboard
	mux.HandleFunc("POST /api/score", handlers.SaveScoreHandler)
	mux.HandleFunc("GET /api/leaderboard", handlers.GetLeaderboardHandler)

	fmt.Println("🐍 Server jalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
