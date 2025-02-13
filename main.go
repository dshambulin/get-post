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

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var task Task
	if err := DB.First(&task, id).Error; err != nil {
		http.Error(w, "Файл не найден", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}

	if err := DB.Save(&task).Error; err != nil {
		http.Error(w, "Ошибка при сохранении в БД", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var task Task
	if err := DB.First(&task, id).Error; err != nil {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	if err := DB.Delete(&task).Error; err != nil {
		http.Error(w, "Ошибка при удалении", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Успешное удаление задачи"})
}

func main() {
	InitDB()

	DB.AutoMigrate(&Task{})

	router := mux.NewRouter()

	router.HandleFunc("/api/messages", GetMessage).Methods("GET")
	router.HandleFunc("/api/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/api/messages/{id}", UpdateMessage).Methods("PATCH")
	router.HandleFunc("/api/messages/{id}", DeleteMessage).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
