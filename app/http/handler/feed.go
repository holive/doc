package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	url2 "net/url"

	"github.com/go-chi/chi"

	"github.com/holive/feedado/app/feed"
)

func (h *Handler) CreateFeed(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	var s feed.Feed
	if err := json.Unmarshal(payload, &s); err != nil {
		respondWithJSONError(w, http.StatusBadRequest, err)
		return
	}

	newFeed, err := h.Services.Feed.Create(r.Context(), &s)
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, *newFeed)
	return
}

func (h *Handler) UpdateFeed(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	var f feed.Feed
	if err := json.Unmarshal(payload, &f); err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.Services.Feed.Update(r.Context(), &f)
	if err != nil {
		respondWithJSONError(w, http.StatusNotFound, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{})
}

func (h *Handler) DeleteFeed(w http.ResponseWriter, r *http.Request) {
	source := chi.URLParam(r, "source")

	url, err := url2.QueryUnescape(source)
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.Services.Feed.Delete(r.Context(), url); err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{})
}

func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	source := chi.URLParam(r, "source")

	url, err := url2.QueryUnescape(source)
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	f, err := h.Services.Feed.FindBySource(r.Context(), url)
	if err != nil {
		respondWithJSONError(w, http.StatusNotFound, err)
		return
	}

	respondWithJSON(w, http.StatusOK, f)
}

func (h *Handler) GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	results, err := h.Services.Feed.FindAll(r.Context(), limit, offset)
	if err != nil {
		respondWithJSONError(w, http.StatusNotFound, err)
		return
	}

	respondWithJSON(w, http.StatusOK, results)
}
