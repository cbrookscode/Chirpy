package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/cbrookscode/Chirpy/internal/database"
)

// in memory struct to keep track of apidata
type apiConfig struct {
	fileserverhits atomic.Int32
	dbqueries      *database.Queries
}

type expectedJSON struct {
	Body string `json:"body"`
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

// replace tsk tsk words with ****
func cleanstring(text string) string {
	bad_words := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	words := strings.Split(text, " ")
	var cleaned []string
	for _, word := range words {
		// val, ok := somemap[key]
		lower_word := strings.ToLower(word)
		if _, found := bad_words[lower_word]; !found {
			cleaned = append(cleaned, word)
		} else {
			cleaned = append(cleaned, "****")
		}
	}
	final := strings.Join(cleaned, " ")
	return final
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnErr struct {
		Errval string `json:"error"`
	}

	// Put msg into response struct for marshalling, marshal into json and provide to responsewriter
	respBody := returnErr{
		Errval: msg,
	}
	data, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling json response: %s", err)
		w.WriteHeader(500)
		w.Write([]byte("Sorry there was an error on the server marshalling your response"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload expectedJSON) {
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	// take out the tsk tsk words, put into response struct for marshalling, marshal into json and provide to responsewriter
	cleaned_txt := cleanstring(payload.Body)
	respBody := returnVals{
		CleanedBody: cleaned_txt,
	}
	data, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling json response: %s", err)
		w.WriteHeader(500)
		w.Write([]byte("Sorry there was an error on the server marshalling your response"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

// handler to validate if charlen of chirp is 140 characters or less
func valcharlen(w http.ResponseWriter, r *http.Request) {
	// use decoder instead of unmarshal so that we dont need to read the entire io.Reader into memory. Unmarshal would be fine if we already had the full byte slice.
	decoder := json.NewDecoder(r.Body)
	params := expectedJSON{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding expected json: %s", err)
		w.WriteHeader(500)
		w.Write([]byte("Sorry there was an error on the server marshalling your response"))
		return
	}

	// no issues with expected json from client, continue

	// Check length, then prepare marshalled json response.
	length := len(params.Body)

	// Invalid chirp
	if length > 140 {
		respondWithError(w, 400, "Chirp is too long")
	}

	// valid chirp
	respondWithJSON(w, 200, params)
}
