package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var task Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	if err := DB.Create(&task).Error; err != nil {
		http.Error(w, "Ошибка при сохранении в БД", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	var tasks []Task

	if err := DB.Find(&tasks).Error; err != nil {
		http.Error(w, "Ошибка при сохранении в БД", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func main() {
	InitDB()

	DB.AutoMigrate(&Task{})

	router := mux.NewRouter()
	router.HandleFunc("/api/messages", GetMessage).Methods("GET")
	router.HandleFunc("/api/messages", CreateMessage).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
