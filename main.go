package main

import (
	"fmt"
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
func (state *apiConfig) metricshandler(reswriter http.ResponseWriter, webrequest *http.Request) {
	val := state.fileserverhits.Load()
	reswriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	reswriter.WriteHeader(200)
	reswriter.Write([]byte(fmt.Sprintf("Hits: %d", val)))
}

// handler method for reseting metric coutner back to 0. needs to be a method for the apiconfig struct as a poitner so that we can access the in memory counter variable when the handler is used.
func (state *apiConfig) resetmetricshandler(reswriter http.ResponseWriter, webrequest *http.Request) {
	state.fileserverhits.Store(0)
	reswriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	reswriter.WriteHeader(200)
	reswriter.Write([]byte("Metrics counter has been reset to 0."))
}

// handler setup for healthz to determine if website is good.
func healthzhandler(reswriter http.ResponseWriter, webrequest *http.Request) {
	reswriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	reswriter.WriteHeader(200)
	reswriter.Write([]byte("OK"))
}

func main() {
	// initialize in memory tracker
	state := &apiConfig{}

	// ServeMux is an HTTP request multiplexer. It matches the URL of each incoming request against a list of registered patterns and calls the handler for the pattern that most closely matches the URL.
	mux := http.NewServeMux()
	port := "8080"
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Registers route handling for /app, which will serve files for webpage. Note: mux.Handle requires an http.handler where mux.handlefunc converts the handle function to a http.Handler under the hood.
	apphandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", state.middlewareMetricsInc(apphandler))

	// Registers /healthz route for telling if webpage is ready
	// Using the Go standard library, you can specify a method like this: [METHOD ][HOST]/[PATH].  Note that all three parts are optional.
	mux.HandleFunc("GET /api/healthz", healthzhandler)

	// Registes /metrics route for getting data on webpage visit number
	mux.HandleFunc("GET /api/metrics", state.metricshandler)

	// Registers /reset route for setting metrics data back to 0
	mux.HandleFunc("POST /api/reset", state.resetmetricshandler)

	// call that starts your HTTP server and keeps it running, continuously listening for incoming HTTP requests.
	server.ListenAndServe()
}
