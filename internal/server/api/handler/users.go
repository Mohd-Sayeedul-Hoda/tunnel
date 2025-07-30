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

func ListUsers(userRepo repositories.UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		page, err := request.ReadInt(r, "page", 1)
		if err != nil {
			badRequestResponse(w, r, err)
			return
		}

		limit, err := request.ReadInt(r, "limit", 20)
		if err != nil {
			badRequestResponse(w, r, err)
			return
		}

		users, err := userRepo.ListUsers(int32(limit), int32((page-1)*limit))
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		respondWithJSON(w, r, http.StatusOK, envelope{
			"status": "success",
			"data": envelope{
				"users": users,
			},
		})

	}
}

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

		respondWithJSON(w, r, http.StatusOK, envelope{
			"status": "success",
			"data": envelope{
				"users": user,
			},
		})
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

		respondWithJSON(w, r, http.StatusCreated, envelope{
			"status": "sucess",
			"data": envelope{
				"users": newUser,
			},
		})
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
		respondWithJSON(w, r, http.StatusOK, envelope{"status": "sucess"})
	}
}
