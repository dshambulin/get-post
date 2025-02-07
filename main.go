package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var task string

func GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Здарова, %s!", task)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Че то не так с JSON", http.StatusBadRequest)
		return
	}
	task = data["task"]
	fmt.Fprintf(w, "Обновлено: %s", task)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/hello", GetHandler).Methods("GET")
	router.HandleFunc("/api/task", PostHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
