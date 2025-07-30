package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/encoding"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/request"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/cache"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories/postgres"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/utils"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/shared/password"
)

func Authenticate(cfg *config.Config, cacheRepo cache.CacheRepo, userRepo repositories.UserRepo) http.HandlerFunc {

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

		jwt, err := utils.CreateToken(user.Id, cfg.Token.AccessTokenExpiredIn, cfg.Token.AccessTokenPrivateKey)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		refreshToken, err := utils.CreateToken(user.Id, cfg.Token.RefreshTokenExpiredIn, cfg.Token.RefreshTokenPrivateKey)
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

		respondWithJSON(w, r, http.StatusOK, envelope{
			"status": "authentication successfull",
			"data": envelope{
				"access_token": jwt.Token,
			}})
	}
}

func RefreshAccessToken(cfg *config.Config, cache cache.CacheRepo, userRepo repositories.UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		RefreshToken := r.CookiesNamed("refresh_token")
		if len(RefreshToken) == 0 {
			notPermittedResponse(w, r)
			return
		}

		tokenClaims, err := utils.ValidetToken(RefreshToken[len(RefreshToken)-1].Value, cfg.Token.RefreshTokenPublicKey)

		if err != nil {
			switch {
			case errors.Is(err, utils.ErrTokenExpired), errors.Is(err, utils.ErrInvalidClaims):
				notPermittedResponse(w, r)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		userIdStr, err := cache.Get(tokenClaims.TokenUuid)
		if err != nil {
			notPermittedResponse(w, r)
			return
		}

		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			notPermittedResponse(w, r)
			return
		}
		user, err := userRepo.GetById(userId)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				notPermittedResponse(w, r)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		jwt, err := utils.CreateToken(user.Id, cfg.Token.AccessTokenExpiredIn, cfg.Token.AccessTokenPrivateKey)
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

		respondWithJSON(w, r, http.StatusOK, envelope{
			"status": "token refreshed successfully",
			"data": envelope{
				"access_token": jwt.Token,
			}})

	}
}
