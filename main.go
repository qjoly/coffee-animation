package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		realIP := r.Header.Get("X-Real-IP")

		if realIP == "" {
			realIP = r.RemoteAddr
		}

		log.Printf("Received request from IP: %s", realIP)

		next.ServeHTTP(w, r)
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Transfer-Encoding", "chunked")

	if !strings.Contains(r.Header.Get("User-Agent"), "curl") {
		fmt.Fprintf(w, "Come on, do a curl !")
		flusher.Flush()
		return
	}

	clearScreen := "\033[2J\033[H"

	for i := 1; i < 10; i++ {
		for _, frame := range GetFrames() {
			fmt.Fprintf(w, clearScreen+"%s\n\n", frame)
			flusher.Flush()
			time.Sleep(500 * time.Millisecond)
		}

	}

	fmt.Fprint(w, clearScreen)
	flusher.Flush()
}

func main() {
	http.Handle("/", Logger(http.HandlerFunc(Handler)))

	fmt.Println("Starting server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
