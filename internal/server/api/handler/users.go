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
		var err error

		v := request.NewValidator()
		filters := request.Filters{}

		filters.Page = request.ReadInt(r, "page", 1, v)
		filters.Limit = request.ReadInt(r, "limit", 20, v)

		if problems := filters.Valid(r.Context(), v); !problems.Valid() {
			failedValidationResponse(w, r, problems)
			return
		}

		users, err := userRepo.ListUsers(int32(filters.Limit), int32((filters.Page-1)*filters.Limit))
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

		v := request.NewValidator()

		var user request.User
		err := encoding.Validated(w, r, v, &user)
		if err != nil {
			switch {
			case errors.Is(err, encoding.ErrInvalidData):
				failedValidationResponse(w, r, v)
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
				v.AddError("email", "a user with this email address already exists")
				failedValidationResponse(w, r, v)
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
