package main

import (
	"fmt"
	"net/http"
	"time"

	ApiService "github.com/artziel/go-api-service"
)

func welcome(w http.ResponseWriter, r *http.Request) {
	ApiService.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"Message": "Welcome to API Service sample in Golang!",
	})
}

func sleep2seconds(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	ApiService.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"Message": "Wait for 2 seconds!",
	})
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Not Found Error")
}

func main() {
	srv := ApiService.NewService(
		"Sample Restful Service",
		"1.0.0b",
		ApiService.ServiceConfig{
			Port: 8080,
		},
	)

	r := srv.Router()
	r.HandleFunc("/welcome", welcome).Methods("GET")
	r.HandleFunc("/sleep-10-seconds", sleep2seconds).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(NotFound)

	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}
}
