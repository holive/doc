package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/holive/doc/app/squads"
)

func (h *Handler) CreateSquad(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	var s squads.Squad
	if err := json.Unmarshal(payload, &s); err != nil {
		respondWithJSONError(w, http.StatusBadRequest, err)
		return
	}

	newSquad, err := h.Services.Squads.Create(r.Context(), s.Name)
	if err != nil {
		if err.Error() == "squad already exists" {
			respondWithJSONError(w, http.StatusForbidden, err)
			return
		}

		respondWithJSONError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, newSquad)
}
