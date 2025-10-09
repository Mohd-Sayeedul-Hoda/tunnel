package handler

import (
	"errors"
	"net/http"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/encoding"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/request"
	tools "github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/utils"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/models"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories/postgres"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/utils"
)

func CreateAPIKey(apiKeyRepo repositories.APIRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		v := request.NewValidator()
		var req request.APIKeys
		err := encoding.Validated(w, r, v, &req)
		if err != nil {
			switch {
			case !v.Valid():
				failedValidationResponse(w, r, v)
			case errors.Is(err, encoding.ErrInvalidData):
				badRequestResponse(w, r, err)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		generatedKeyDetails, err := utils.GenerateAPIKeyToken(32)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		userDetails := tools.ContextGetToken(r)
		apikey := models.APIKey{
			Name:        req.Name,
			Prefix:      generatedKeyDetails.Prefix,
			APIkeyToken: generatedKeyDetails.FullKey,
			APIKeyHash:  generatedKeyDetails.KeyHash,
			ExpireAt:    req.ExpiresAt,
			UserId:      userDetails.UserID,
		}

		err = apiKeyRepo.CreateAPIKey(&apikey)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		respondWithJSON(w, r, http.StatusCreated, envelope{
			"status": "success",
			"data": envelope{
				"api_key": apikey,
			},
		})

	})
}

func ListAPIKey(apiKeyRepo repositories.APIRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		v := request.NewValidator()
		page := request.Pagination{}

		page.Page = request.ReadInt(r, v, "page", 1)
		page.Limit = request.ReadInt(r, v, "limit", 20)
		v = page.Valid(r.Context(), v)
		if !v.Valid() {
			failedValidationResponse(w, r, v)
			return
		}

		userDetails := tools.ContextGetToken(r)
		keys, err := apiKeyRepo.ListAPIKeys(userDetails.UserID, page.Page, (page.Page-1)*page.Limit)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		respondWithJSON(w, r, http.StatusOK, envelope{
			"status": "success",
			"data": envelope{
				"api_keys": keys,
			},
		})
	})
}

func DeleteAPIKey(apiKeyRepo repositories.APIRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := request.ReadIDParam(r)
		if err != nil {
			notFoundResponse(w, r)
			return
		}

		token := tools.ContextGetToken(r)

		err = apiKeyRepo.DeleteAPIKey(token.UserID, id)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				notFoundResponse(w, r)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		respondWithJSON(w, r, http.StatusOK, envelope{
			"status": "success",
		})
	})
}
