package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/holive/feedado/app/feedado"
)

type Handler struct {
	Services *feedado.Services
}

type WorkerHandler struct {
	Services *feedado.WorkerServices
}

type Message struct {
	Name string
	Body string
	Time int64
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	health(w, r)
}

func (h *WorkerHandler) Health(w http.ResponseWriter, r *http.Request) {
	health(w, r)
}

func health(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(map[string]string{"status": "ok"})
	if err != nil {
		fmt.Println("could not Marshal health json response")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func respondWithJSONError(w http.ResponseWriter, code int, err error) {
	payload := map[string]interface{}{
		"error": err.Error(),
	}
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("{ \"error\": \"could not marshal error\"}"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(response)
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
	}
}
