package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/encoding"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/request"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/models"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories/postgres"
)

func GetUsers(userRepo repositories.UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("user_id")
		userIdstr, err := strconv.Atoi(userId)
		if err != nil {
			badRequestResponse(w, r, errors.New("user_id should be a valid id"))
			return
		}

		user, err := userRepo.GetById(userIdstr)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				notFoundResponse(w, r)
			default:
				serverErrorResponse(w, r, err)
			}
			return
		}
		err = encoding.EncodeJson(w, r, http.StatusOK, envelope{"data": user})
		if err != nil {
			serverErrorResponse(w, r, err)
		}
	}
}

func CreateUser(userRepo repositories.UserRepo) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		user, problem, err := encoding.Validated[*request.User](w, r)
		if err != nil {
			badRequestResponse(w, r, err)
			return
		}

		if len(problem) > 0 {
			failedValidationResponse(w, r, problem)
			return
		}

		newUser := &models.User{
			Email: user.Email,
			Name:  user.Name,
		}

		err = userRepo.Create(newUser)
		if err != nil {
			serverErrorResponse(w, r, err)
			return
		}

		err = encoding.EncodeJson(w, r, http.StatusCreated, newUser)
		if err != nil {
			serverErrorResponse(w, r, err)
		}
	}
}
