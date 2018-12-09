package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	flagPort = flag.String("port", "9000", "Port to listen on")
)

//OpenPortalRequest - Struct for requests OpenPortalRequest
// {
//     "user_id": "B5AB66DB-A91C-43AD-8DFE-732520343F12",
//     "portal_alias": "3CDE4582-D55A-4520-A42F-91D66FABD7FF",
//     "timestamp": 1544314490429
// }
type OpenPortalRequest struct {
	UserID      string `json:"user_id"`
	PortalAlias string `json:"portal_alias"`
	Timestamp   int    `json:"timestamp"`
}

var results []OpenPortalRequest

// HealthHandler - Default handler
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

// OpenPortalRequestHandler - mock portal request
func OpenPortalRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var request OpenPortalRequest
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		fmt.Println(request.PortalAlias)

		results = append(results, request)

		fmt.Fprintln(w, "POST done")

	} else if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	flag.Parse()
}

//curl -d '{"user_id":"value1", "portal_alias":"value2", "timestamp": 1544314081995}' -H "Content-Type: application/json" -X POST http://localhost:9000/portal
func main() {
	// results = append(results, time.Now().Format(time.RFC3339))

	mux := http.NewServeMux()
	mux.HandleFunc("/", HealthHandler)
	mux.HandleFunc("/open", OpenPortalRequestHandler)

	log.Printf("listening on port %s", *flagPort)
	log.Fatal(http.ListenAndServe(":"+*flagPort, mux))
}
