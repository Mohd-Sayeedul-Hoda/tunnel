package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/encoding"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/request"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/models"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories/postgres"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/utils"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/shared/password"
)

func SendEmailVerficationOtp(cfg *config.Config, userRepo repositories.UserRepo, emailRepo repositories.EmailOtpRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		v := request.NewValidator()

		var req request.BaseEmail
		err := encoding.Validated(w, r, v, &req)
		if err != nil {
			switch {
			case errors.Is(err, encoding.ErrInvalidRequest):
				badRequestResponse(w, r, err)
			case !v.Valid():
				failedValidationResponse(w, r, v)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		user, err := userRepo.GetByEmail(req.Email)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				notFoundResponse(w, r)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		if user.EmailVerified {
			errorResponse(w, r, http.StatusBadRequest, "email already verified")
			return
		}

		todayMidnightUtc := time.Now().UTC().Truncate(24 * time.Hour)
		totalSend, err := emailRepo.CountOtpsAfterUtcTime(req.Email, models.EmailVerificationOtpType, todayMidnightUtc)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		if totalSend >= 3 {
			errorResponse(w, r, http.StatusTooManyRequests, "rate limit exceeded: maximum emails per day reached")
			return
		}

		emailOtp := utils.GenerateToken(6)
		err = emailRepo.CreateOtp(req.Email,
			utils.HashOtp(cfg.EmailOtpSalt, emailOtp),
			models.EmailVerificationOtpType,
			time.Now().Add(cfg.EmailOtpExpiredIn),
		)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		//WARN: remove and email smtp client to send email
		slog.Info("Email Otp", slog.String("email", req.Email), slog.String("otp", emailOtp), slog.String("otp-type", string(models.EmailVerificationOtpType)))

		respondWithJSON(w, r, http.StatusCreated, envelope{
			"status": "success",
		})
	})
}

func VerifyEmailVerficationOtp(cfg *config.Config, userRepo repositories.UserRepo, emailRepo repositories.EmailOtpRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		v := request.NewValidator()

		var req request.VerifyUserOTP
		err := encoding.Validated(w, r, v, &req)
		if err != nil {
			switch {
			case !v.Valid():
				failedValidationResponse(w, r, v)
			case errors.Is(err, encoding.ErrInvalidRequest):
				badRequestResponse(w, r, err)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		user, err := userRepo.GetByEmail(req.Email)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				notFoundResponse(w, r)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		if user.EmailVerified {
			errorResponse(w, r, http.StatusBadRequest, "email already verified")
			return
		}

		emailOtp, err := emailRepo.GetOtp(req.Email, models.EmailVerificationOtpType)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				notFoundResponse(w, r)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		if emailOtp.IsInvalidated || emailOtp.Used || emailOtp.Attempts > 3 || emailOtp.ExpiresAt.Before(time.Now()) {
			errorResponse(w, r, http.StatusUnauthorized, "invalid otp")
			return
		}

		hashOtp := utils.HashOtp(cfg.EmailOtpSalt, req.EmailOtp)
		if !(hashOtp == emailOtp.EmailOtp) {
			if emailOtp.Attempts < 3 {
				err = emailRepo.IncreaseOtpAttempt(emailOtp.Id)
				if err != nil {
					ServerErrorResponse(w, r, err)
					return
				}
			} else {
				err = emailRepo.IncreaseAttemptAndInvalidateOtp(emailOtp.Id)
				if err != nil {
					ServerErrorResponse(w, r, err)
					return
				}
			}
			errorResponse(w, r, http.StatusUnauthorized, "invalid otp")
			return
		}

		err = emailRepo.VerifyOtp(emailOtp.Id)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}
		err = userRepo.VerifyUserEmail(user.Id)
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

func SendForgotPasswordLink(cfg *config.Config, userRepo repositories.UserRepo, emailRepo repositories.EmailOtpRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		v := request.NewValidator()
		var req request.BaseEmail
		err := encoding.Validated(w, r, v, &req)
		if err != nil {
			switch {
			case errors.Is(err, encoding.ErrInvalidData):
				failedValidationResponse(w, r, v)
			case errors.Is(err, encoding.ErrInvalidRequest):
				badRequestResponse(w, r, err)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		_, err = userRepo.GetByEmail(req.Email)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				notFoundResponse(w, r)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		todayMidnightUtc := time.Now().UTC().Truncate(24 * time.Hour)
		totalSend, err := emailRepo.CountOtpsAfterUtcTime(req.Email, models.ForgotPasswordOtpType, todayMidnightUtc)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		if totalSend >= 3 {
			errorResponse(w, r, http.StatusTooManyRequests, "rate limit exceeded: maximum emails per day reached")
			return
		}

		otp := utils.GenerateToken(32)
		hashOtp := utils.HashOtp(cfg.EmailOtpSalt, otp)
		err = emailRepo.CreateOtp(
			req.Email,
			hashOtp,
			models.ForgotPasswordOtpType,
			time.Now().Add(cfg.EmailOtpExpiredIn),
		)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		encodedToken := base64.StdEncoding.EncodeToString(fmt.Appendf([]byte{}, "%s|%s", req.Email, otp))

		//TODO: change route here
		//WARN: remove and email smtp client to send email
		url := fmt.Sprintf("http://localhost:5173/forgot-password?token=%s", encodedToken)
		slog.Info("Email Otp", slog.String("email", req.Email), slog.String("url", url), slog.String("otp-type", string(models.ForgotPasswordOtpType)))

		respondWithJSON(w, r, http.StatusCreated, envelope{
			"status": "success",
		})
	})
}

func VerifyForgotPasswordLink(cfg *config.Config, userRepo repositories.UserRepo, emailRepo repositories.EmailOtpRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		v := request.NewValidator()
		var req request.ForgotPasswordVerify
		err := encoding.Validated(w, r, v, &req)
		if err != nil {
			switch {
			case errors.Is(err, encoding.ErrInvalidRequest):
				badRequestResponse(w, r, err)
			case errors.Is(err, encoding.ErrInvalidData):
				failedValidationResponse(w, r, v)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		decodeByte, err := base64.StdEncoding.DecodeString(req.EmailOtp)
		if err != nil {
			badRequestResponse(w, r, errors.New("invalid base64 encoding in otp"))
			return
		}

		tokenData := strings.SplitN(string(decodeByte), "|", 2)
		email := tokenData[0]
		token := tokenData[1]

		_, err = userRepo.GetByEmail(email)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				notFoundResponse(w, r)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		emailOtp, err := emailRepo.GetOtp(email, models.ForgotPasswordOtpType)
		if err != nil {
			switch {
			case errors.Is(err, postgres.ErrNotFound):
				notFoundResponse(w, r)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		if emailOtp.IsInvalidated || emailOtp.Used || emailOtp.Attempts > 3 || emailOtp.ExpiresAt.Before(time.Now()) {
			errorResponse(w, r, http.StatusUnauthorized, "invalid otp")
			return
		}

		hashOtp := utils.HashOtp(cfg.EmailOtpSalt, token)
		if !(hashOtp == emailOtp.EmailOtp) {
			if emailOtp.Attempts < 3 {
				err = emailRepo.IncreaseOtpAttempt(emailOtp.Id)
				if err != nil {
					ServerErrorResponse(w, r, err)
					return
				}
			} else {
				err = emailRepo.IncreaseAttemptAndInvalidateOtp(emailOtp.Id)
				if err != nil {
					ServerErrorResponse(w, r, err)
					return
				}
			}
			errorResponse(w, r, http.StatusUnauthorized, "invalid otp")
			return
		}

		err = emailRepo.VerifyOtp(emailOtp.Id)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		hash, err := password.SetPassword(req.Password)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		err = userRepo.UpdateUserPassword(email, hash)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		respondWithJSON(w, r, http.StatusOK, envelope{
			"status": "success",
		})

	})
}
