package handler

import (
	
	"encoding/json"
	"net/http"
	"time"

	"http-server-projeto-korp/internal/model"

)

func ProjetoKorpHandler(w http.ResponseWriter, r *http.Request) {
	
	response := model.Response{
		Nome:    "Projeto Korp",
		Horario: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)

}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}