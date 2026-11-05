package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"snake-game/database"
	"snake-game/handlers"
)

// File di folder static/ ikut "dibungkus" ke dalam program saat build,
// jadi di hosting mana pun tidak ada risiko folder static hilang.
//go:embed static
var staticFiles embed.FS

func main() {
	database.Connect()
	defer database.DB.Close()

	mux := http.NewServeMux()

	// Sajikan file statis (HTML, CSS, JS) dari folder static/
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	// Endpoint API untuk leaderboard
	mux.HandleFunc("POST /api/score", handlers.SaveScoreHandler)
	mux.HandleFunc("GET /api/leaderboard", handlers.GetLeaderboardHandler)

	// Pakai http.Server agar bisa atur timeout & mati dengan rapi.
	// Layanan hosting menentukan port lewat env PORT; lokal tetap 8080.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Jalankan server di goroutine (proses paralel ringan milik Go).
	go func() {
		fmt.Println("🐍 Server jalan di port " + port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Tunggu Ctrl+C, lalu matikan server dengan rapi (request yang berjalan diselesaikan dulu).
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	fmt.Println("Server dimatikan.")
}