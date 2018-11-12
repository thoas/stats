package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/thoas/stats"
)

func StatsMiddleware(middleware *stats.Stats) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			beginning, recorder := middleware.Begin(w)
			next.ServeHTTP(w, r)
			middleware.End(beginning, stats.WithRecorder(recorder))
		}
		return http.HandlerFunc(fn)
	}
}

func main() {
	mux := chi.NewRouter()

	Stats := stats.New()
	mux.Use(StatsMiddleware(Stats))

	mux.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hello\": \"world\"}"))
	})

	mux.Get("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		stats := Stats.Data()

		b, _ := json.Marshal(stats)

		w.Write(b)
	})

	fmt.Println("Serving on :8080")
	http.ListenAndServe(":8080", mux)
}
