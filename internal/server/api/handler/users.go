package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/encoding"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/request"
	tools "github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/utils"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/cache"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/models"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories/postgres"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/utils"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/shared/password"
)

func ListUsers(userRepo repositories.UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		v := request.NewValidator()
		filters := request.Pagination{}

		filters.Page = request.ReadInt(r, v, "page", 1)
		filters.Limit = request.ReadInt(r, v, "limit", 20)
		v = filters.Valid(r.Context(), v)
		if !v.Valid() {
			failedValidationResponse(w, r, v)
			return
		}

		users, err := userRepo.ListUsers(filters.Limit, (filters.Page-1)*filters.Limit)
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

		token := tools.ContextGetToken(r)

		user, err := userRepo.GetById(token.UserID)
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

func SignupUser(userRepo repositories.UserRepo) http.HandlerFunc {

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
			"status": "success",
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

func AuthenticateUser(cfg *config.Config, cacheRepo cache.CacheRepo, userRepo repositories.UserRepo) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		v := request.NewValidator()

		var req request.Login
		err := encoding.Validated(w, r, v, &req)
		if err != nil {
			switch {
			case errors.Is(err, encoding.ErrInvalidData):
				failedValidationResponse(w, r, v)
			default:
				badRequestResponse(w, r, err)
			}
			return
		}

		user, err := userRepo.GetByEmail(req.Email)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				InvalidCredentialsResponse(w, r)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		matched, err := password.MatchPassword(user.PasswordHash, req.Password)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}
		if !matched {
			InvalidCredentialsResponse(w, r)
			return
		}

		jwt, err := utils.CreateToken(user, cfg.Token.AccessTokenExpiredIn, cfg.Token.AccessTokenPrivateKey)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		refreshToken, err := utils.CreateToken(user, cfg.Token.RefreshTokenExpiredIn, cfg.Token.RefreshTokenPrivateKey)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		now := time.Now()
		err = cacheRepo.Set(refreshToken.TokenUuid, user.Id, time.Unix(jwt.ExpiresIn, 0).Sub(now))
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		accessCookie := http.Cookie{
			Name:     "jwt",
			Value:    jwt.Token,
			Path:     "/",
			MaxAge:   int(cfg.Token.AccessTokenMaxAge) * 60,
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &accessCookie)

		refreshCookie := http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken.Token,
			Path:     "/",
			MaxAge:   int(cfg.Token.RefreshTokenMaxAge) * 60,
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &refreshCookie)

		loginCookie := http.Cookie{
			Name:     "logged_in",
			Value:    "true",
			Path:     "/",
			MaxAge:   int(cfg.Token.AccessTokenMaxAge) * 60,
			Secure:   false,
			HttpOnly: false,
		}
		http.SetCookie(w, &loginCookie)

		response := envelope{
			"status": "success",
		}

		if cfg.AppEnv != "prod" {
			response["data"] = envelope{
				"access_token":  jwt.Token,
				"refresh_token": refreshToken.Token,
			}
		}

		respondWithJSON(w, r, http.StatusOK, response)
	}
}

func LogoutUser(cfg *config.Config, cacheRepo cache.CacheRepo, userRepo repositories.UserRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		refreshToken := r.CookiesNamed("refresh_token")
		if len(refreshToken) == 0 {
			NotPermittedResponse(w, r)
			return
		}

		tokenClaims, err := utils.ValidateToken(refreshToken[len(refreshToken)-1].Value, cfg.Token.RefreshTokenPublicKey)
		if err != nil {
			switch {
			case errors.Is(err, utils.ErrTokenExpired), errors.Is(err, utils.ErrInvalidClaims):
				NotPermittedResponse(w, r)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		_, err = cacheRepo.Delete(tokenClaims.TokenUuid)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		expired := time.Now().Add(-time.Hour * 24)
		jwtCookie := http.Cookie{
			Name:    "jwt",
			Value:   "",
			Expires: expired,
		}
		http.SetCookie(w, &jwtCookie)

		refreshCookie := http.Cookie{
			Name:    "refresh_token",
			Value:   "",
			Expires: expired,
		}
		http.SetCookie(w, &refreshCookie)

		loginCookie := http.Cookie{
			Name:    "logged_in",
			Value:   "",
			Expires: expired,
		}
		http.SetCookie(w, &loginCookie)

		respondWithJSON(w, r, http.StatusOK, envelope{
			"status": "success",
		})

	})
}

func RefreshUserAccessToken(cfg *config.Config, cache cache.CacheRepo, userRepo repositories.UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		RefreshToken, err := r.Cookie("refresh_token")
		if err != nil {
			NotPermittedResponse(w, r)
			return
		}

		tokenClaims, err := utils.ValidateToken(RefreshToken.Value, cfg.Token.RefreshTokenPublicKey)
		if err != nil {
			switch {
			case errors.Is(err, utils.ErrTokenExpired), errors.Is(err, utils.ErrInvalidClaims):
				NotPermittedResponse(w, r)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		userIdStr, err := cache.Get(tokenClaims.TokenUuid)
		if err != nil {
			NotPermittedResponse(w, r)
			return
		}

		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			NotPermittedResponse(w, r)
			return
		}
		user, err := userRepo.GetById(userId)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				NotPermittedResponse(w, r)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		jwt, err := utils.CreateToken(user, cfg.Token.AccessTokenExpiredIn, cfg.Token.AccessTokenPrivateKey)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		accessCookie := http.Cookie{
			Name:     "jwt",
			Value:    jwt.Token,
			Path:     "/",
			MaxAge:   int(cfg.Token.AccessTokenMaxAge) * 60,
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &accessCookie)

		loginCookie := http.Cookie{
			Name:     "logged_in",
			Value:    "true",
			Path:     "/",
			MaxAge:   int(cfg.Token.AccessTokenMaxAge) * 60,
			Secure:   false,
			HttpOnly: false,
		}
		http.SetCookie(w, &loginCookie)

		response := envelope{
			"status": "success",
		}

		if cfg.AppEnv != "prod" {
			response["data"] = envelope{
				"access_token": jwt.Token,
			}
		}

		respondWithJSON(w, r, http.StatusOK, response)

	}
}
