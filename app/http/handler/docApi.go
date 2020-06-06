package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/holive/doc/app/docApi"
)

func (h *Handler) CreateDocApi(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	var s docApi.DocApi
	if err := json.Unmarshal(payload, &s); err != nil {
		respondWithJSONError(w, http.StatusBadRequest, err)
		return
	}

	newDoc, err := h.Services.DocApi.Create(r.Context(), &s)
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, *newDoc)
}

func (h *Handler) GetDocApi(w http.ResponseWriter, r *http.Request) {
	squad := chi.URLParam(r, "squad")
	projeto := chi.URLParam(r, "projeto")
	versao := chi.URLParam(r, "versao")

	if squad == "" || projeto == "" || versao == "" {
		respondWithJSONError(w, http.StatusInternalServerError, errors.New("missing url param"))
		return
	}

	f, err := h.Services.DocApi.Find(r.Context(), squad, projeto, versao)
	if err != nil {
		respondWithJSONError(w, http.StatusNotFound, err)
		return
	}

	respondWithJSON(w, http.StatusOK, f)
}

func (h *Handler) GetAllDocs(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	results, err := h.Services.DocApi.FindAll(r.Context(), limit, offset)
	if err != nil {
		respondWithJSONError(w, http.StatusNotFound, err)
		return
	}

	respondWithJSON(w, http.StatusOK, results)
}

func (h *Handler) DeleteDocApi(w http.ResponseWriter, r *http.Request) {
	squad := chi.URLParam(r, "squad")
	projeto := chi.URLParam(r, "projeto")
	versao := chi.URLParam(r, "versao")

	if squad == "" || projeto == "" || versao == "" {
		respondWithJSONError(w, http.StatusInternalServerError, errors.New("missing url param"))
		return
	}

	if err := h.Services.DocApi.Delete(r.Context(), squad, projeto, versao); err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{})
}
