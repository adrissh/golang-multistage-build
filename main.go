package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
)

type Response struct {
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       map[string]any
}

func main() {

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	serverInfo := map[string]string{
		"hostname":         hostname,
		"go_version":       runtime.Version(),
		"operation_system": runtime.GOOS,
		"architecture":     runtime.GOARCH,
	}

	port := 8000

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Status:     "success",
			StatusCode: 200,
			Message:    "Successfully fetched server information",
			Data: map[string]any{
				"server_info": serverInfo,
				"request_info": map[string]string{
					"source_ip":  r.RemoteAddr,
					"user_agent": r.Header.Get("User-Agent"),
				},
			},
		}
		json.NewEncoder(w).Encode(response)
		log.Printf("Called endpoint '/' by %v", r.RemoteAddr)

	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		res := map[string]string{
			"status": "UP",
		}
		json.NewEncoder(w).Encode(res)
		log.Printf("Called endpoint '/health' by %v", r.RemoteAddr)
	})

	log.Printf("Server running in :%v", port)

	errServe := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if errServe != nil {
		log.Fatal(errServe)
	}
}
