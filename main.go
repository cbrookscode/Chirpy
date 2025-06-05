package main

import (
	"net/http"
)

// Setting up handler to match expected signature for http.HandlerFunc
func myhandler(reswriter http.ResponseWriter, webrequest *http.Request) {
	reswriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	reswriter.WriteHeader(200)
	reswriter.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	port := "8080"
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Register hanlder for standing serving files for webpage
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	//Register /healthz handler. Handling /healthz requests to make sure page is website is ready
	mux.HandleFunc("/", myhandler)

	server.ListenAndServe()
}
