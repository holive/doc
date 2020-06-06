package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/holive/feedado/app/user"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	var u user.User
	if err := json.Unmarshal(payload, &u); err != nil {
		respondWithJSONError(w, http.StatusBadRequest, err)
		return
	}

	newUser, err := h.Services.User.Create(r.Context(), &u)
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, *newUser)
	return
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	var u user.User
	if err := json.Unmarshal(payload, &u); err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.Services.User.Update(r.Context(), &u)
	if err != nil {
		respondWithJSONError(w, http.StatusNotFound, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{})
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")

	if err := h.Services.User.Delete(r.Context(), email); err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{})
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")

	u, err := h.Services.User.Find(r.Context(), email)
	if err != nil {
		respondWithJSONError(w, http.StatusNotFound, err)
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	results, err := h.Services.User.FindAll(r.Context(), limit, offset)
	if err != nil {
		respondWithJSONError(w, http.StatusNotFound, err)
		return
	}

	respondWithJSON(w, http.StatusOK, results)
}
