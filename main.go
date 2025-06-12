package main

// underscore in before import url means its only being used for side effects. In this case it allows us to register the driver needed to open database connection using sql.open()
import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/cbrookscode/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// initialize in memory tracker
	state := &apiConfig{}

	// load .env file, get dburl from file, open connection based on dburl, setup new database queries pointer using database.New(), store in app state variable.
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}
	dbQueries := database.New(db)
	state.dbqueries = dbQueries

	// ServeMux is an HTTP request multiplexer. It matches the URL of each incoming request against a list of registered patterns and calls the handler for the pattern that most closely matches the URL.
	mux := http.NewServeMux()
	port := "8080"
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Registers route handling for /app, which will serve files for webpage. Note: mux.Handle requires an http.handler where mux.handlefunc converts the handle function to a http.Handler under the hood.
	apphandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", state.middlewareMetricsInc(apphandler))

	// Registers /healthz route for telling if webpage is ready
	// Using the Go standard library, you can specify a method like this: [METHOD ][HOST]/[PATH].  Note that all three parts are optional.
	mux.HandleFunc("GET /api/healthz", healthzhandler)

	mux.HandleFunc("POST /api/validate_chirp", valcharlen)

	// Registes /metrics route for getting data on webpage visit number
	mux.HandleFunc("GET /admin/metrics", state.metricshandler)

	// Registers /reset route for setting metrics data back to 0
	mux.HandleFunc("POST /admin/reset", state.resetmetricshandler)

	// call that starts your HTTP server and keeps it running, continuously listening for incoming HTTP requests.
	server.ListenAndServe()
}
