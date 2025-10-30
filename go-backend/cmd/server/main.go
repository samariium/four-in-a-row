package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/yourname/fourinarow/internal/config"
	"github.com/yourname/fourinarow/internal/game"
	"github.com/yourname/fourinarow/internal/store"
)

func main() {
	cfg := config.Load()

	// Connect to Mongo (10s timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoStore, err := store.NewMongoStore(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("mongo connect error: %v", err)
	}

	mgr := game.NewManager(mongoStore, cfg.MatchBotAfterMs, cfg.RejoinGraceMs, cfg.BotMoveDelayMs)

	// HTTP mux
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	// Leaderboard (log real error so we can diagnose 500s)
	mux.HandleFunc("/leaderboard", func(w http.ResponseWriter, r *http.Request) {
		top, err := mongoStore.TopPlayers(r.Context(), 10)
		if err != nil {
			log.Printf("leaderboard error: %v", err) // <— view this in server console
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(top)
	})

	// WebSocket
	mux.HandleFunc("/ws", mgr.HandleWS)

	// Basic CORS wrapper so the React app can call the API
	handler := cors(mux)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	log.Printf("🚀 Go backend on http://localhost:%s", cfg.Port)
	log.Fatal(server.ListenAndServe())
}

// very small CORS middleware (adjust origin if you want)
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
