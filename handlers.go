package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

// in memory struct to keep track of apidata
type apiConfig struct {
	fileserverhits atomic.Int32
}

// wrapper for app handler so that we can increment filserverhits accurately.
func (state *apiConfig) middlewareMetricsInc(og_handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state.fileserverhits.Add(1)
		og_handler.ServeHTTP(w, r)
	})
}

// handler method for metrics counter. needs to be a method for the apiconfig struct as a poitner so that we can access the in memory counter variable when the handler is used.
func (state *apiConfig) metricshandler(r http.ResponseWriter, w *http.Request) {
	val := state.fileserverhits.Load()
	r.Header().Set("Content-Type", "text/html; charset=utf-8")
	r.WriteHeader(200)
	html := fmt.Sprintf(`
	<html>
  		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
  		</body>
	</html>
	`, val)
	r.Write([]byte(html))
}

// handler method for reseting metric coutner back to 0. needs to be a method for the apiconfig struct as a poitner so that we can access the in memory counter variable when the handler is used.
func (state *apiConfig) resetmetricshandler(r http.ResponseWriter, w *http.Request) {
	state.fileserverhits.Store(0)
	r.Header().Set("Content-Type", "text/plain; charset=utf-8")
	r.WriteHeader(200)
	r.Write([]byte("Metrics counter has been reset to 0."))
}

// handler setup for healthz to determine if website is good.
func healthzhandler(r http.ResponseWriter, w *http.Request) {
	r.Header().Set("Content-Type", "text/plain; charset=utf-8")
	r.WriteHeader(200)
	r.Write([]byte("OK"))
}

// handler to validate if charlen of chirp is 140 characters or less
func valcharlen(w http.ResponseWriter, r *http.Request) {
	type expected struct {
		Body string `json:"body"`
	}

	// use decoder instead of unmarshal so that we dont need to read the entire io.Reader into memory. Unmarshal would be fine if we already had the full byte slice.
	decoder := json.NewDecoder(r.Body)
	params := expected{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding expected json: %s", err)
		w.WriteHeader(500)
		w.Write([]byte("Sorry there was an error on the server marshalling your response"))
		return
	}

	// no issues with expected json from client, continue

	// Setup structs for returning true or false for valid chirp, and for errors
	type returnVals struct {
		Body bool `json:"valid"`
	}

	type returnErr struct {
		Errval string `json:"error"`
	}

	// Check length, then prepare marshalled json response.
	length := len(params.Body)

	// Invalid chirp
	if length > 140 {
		respBody := returnErr{
			Errval: "Chirp is too long",
		}
		data, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling json response: %s", err)
			w.WriteHeader(500)
			w.Write([]byte("Sorry there was an error on the server marshalling your response"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(data)
		return
	}

	// valid chirp
	respBody := returnVals{
		Body: true,
	}
	data, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling json response: %s", err)
		w.WriteHeader(500)
		w.Write([]byte("Sorry there was an error on the server marshalling your response"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}
