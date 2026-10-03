// Command hexclock serves the current UTC time as JSON.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// timeResponse is the body returned by GET /.
type timeResponse struct {
	UTC  string `json:"utc"`
	Unix int64  `json:"unix"`
	Hex  string `json:"hex"`
}

// newHandler builds the service's routes. now is injected so tests can pin the clock.
func newHandler(now func() time.Time) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		t := now().UTC()
		writeJSON(w, timeResponse{
			UTC:  t.Format(time.RFC3339Nano),
			Unix: t.Unix(),
			Hex:  "0x" + strconv.FormatInt(t.Unix(), 16),
		})
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"status": "ok"})
	})
	return mux
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func main() {
	addr := ":" + envOr("PORT", "8080")
	srv := &http.Server{
		Addr:              addr,
		Handler:           newHandler(time.Now),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("hexclock listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
