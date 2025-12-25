package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Skorgum/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type polkaWebhook struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	// Verify Polka key
	polkaKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	if polkaKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	var webhook polkaWebhook
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if webhook.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeToChirpyRed(r.Context(), webhook.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to upgrade user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
