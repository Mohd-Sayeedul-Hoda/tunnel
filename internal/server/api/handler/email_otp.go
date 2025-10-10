package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/encoding"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/request"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/models"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories/postgres"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/utils"
)

func SendEmailOtp(cfg *config.Config, userRepo repositories.UserRepo, emailOtpRepo repositories.EmailOtpRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		v := request.NewValidator()

		var req request.SendOtp
		err := encoding.Validated(w, r, v, &req)
		if err != nil {
			switch {
			case errors.Is(err, encoding.ErrInvalidData):
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
		totalSend, err := emailOtpRepo.CountOtpsAfterUtcTime(req.Email, models.EmailVerificationOtpType, todayMidnightUtc)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		if totalSend >= cfg.TotalAllowEmailForType {
			errorResponse(w, r, http.StatusTooManyRequests, "rate limit exceeded: maximum emails per day reached")
			return
		}

		emailOtp := utils.GenerateEmailOtpToken(cfg.EmailOtpLenght)
		err = emailOtpRepo.CreateOtp(req.Email,
			utils.HashOtp(cfg.EmailOtpSalt, emailOtp),
			models.OtpType(req.OtpType),
			time.Now().Add(cfg.EmailOtpExpiredIn),
		)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		//WARN: remove and email smtp client to send email
		slog.Info("Email Otp", slog.String("email", req.Email), slog.String("otp", emailOtp), slog.String("otp-type", req.OtpType))

		respondWithJSON(w, r, http.StatusCreated, envelope{
			"status": "success",
		})
	})
}

func VerifyEmailOtp(cfg *config.Config, userRepo repositories.UserRepo, emailRepo repositories.EmailOtpRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		v := request.NewValidator()

		var req request.VerifyOtp
		err := encoding.Validated(w, r, v, &req)
		if err != nil {
			switch {
			case errors.Is(err, encoding.ErrInvalidData):
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

		emailOtp, err := emailRepo.GetOtp(req.Email, models.OtpType(req.OtpType))
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
