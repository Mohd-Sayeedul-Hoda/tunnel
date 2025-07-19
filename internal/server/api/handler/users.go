package handler

import (
	"errors"
	"net/http"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/encoding"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/request"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/models"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories/postgres"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/shared/password"
)

func GetUsers(userRepo repositories.UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId, err := request.ReadIDParam(r)
		if err != nil {
			notFoundResponse(w, r)
			return
		}

		user, err := userRepo.GetById(userId)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				notFoundResponse(w, r)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		err = encoding.EncodeJson(w, r, http.StatusOK, envelope{"data": user})
		if err != nil {
			ServerErrorResponse(w, r, err)
		}
	}
}

func CreateUser(userRepo repositories.UserRepo) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var user request.User
		problem, err := encoding.Validated(w, r, &user)
		if err != nil {
			switch {
			case errors.Is(err, encoding.ErrInvalidRequest):
				badRequestResponse(w, r, err)
			case errors.Is(err, encoding.ErrInvalidData):
				failedValidationResponse(w, r, problem)
			default:
				badRequestResponse(w, r, err)
			}
			return
		}

		newUser := &models.User{
			Email: user.Email,
			Name:  user.Name,
		}

		hash, err := password.SetPassword(user.Password)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}
		newUser.PasswordHash = hash

		err = userRepo.Create(newUser)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrUniqueViolation):
				problem.AddError("email", "a user with this email address already exists")
				failedValidationResponse(w, r, problem)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		err = encoding.EncodeJson(w, r, http.StatusCreated, envelope{"data": newUser})
		if err != nil {
			ServerErrorResponse(w, r, err)
		}
	}
}

func DeleteUser(userRepo repositories.UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := request.ReadIDParam(r)
		if err != nil {
			notFoundResponse(w, r)
			return
		}

		err = userRepo.Delete(id)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				notFoundResponse(w, r)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}
		respondWithJSON(w, r, http.StatusOK, envelope{"status": "user deleted"})
	}
}

func Authenticate(userRepo repositories.UserRepo) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var req request.Login
		problems, err := encoding.Validated(w, r, &req)
		if err != nil {
			switch {
			case errors.Is(err, encoding.ErrInvalidRequest):
				badRequestResponse(w, r, err)
			case errors.Is(err, encoding.ErrInvalidData):
				failedValidationResponse(w, r, problems)
			default:
				badRequestResponse(w, r, err)
			}
			return
		}

		user, err := userRepo.GetByEmail(req.Email)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				invalidCredentialsResponse(w, r)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		matched, err := password.MatchPassword(user.Password, req.Password)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}
		if !matched {
			invalidCredentialsResponse(w, r)
			return
		}

	}
}
