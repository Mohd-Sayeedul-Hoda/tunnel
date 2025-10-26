package api

import (
	"net/http"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/handler"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/cache"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/shared/config"
)

func AddRoute(mux *http.ServeMux, cfg *config.Config, cacheRepo cache.CacheRepo, userRepo repositories.UserRepo, apiKeyRepo repositories.APIRepo, emailOtpRepo repositories.EmailOtpRepo) {

	// general
	mux.HandleFunc("/", handler.HandleRoot())
	mux.HandleFunc("GET /api/v1/healthcheck", handler.HealthCheck(cfg))

	// users
	requireVerified := newAuthenticateAndVerifyMiddleware(cfg)
	mux.Handle("GET /api/v1/users/me", requireVerified(handler.GetUsers(userRepo)))
	mux.Handle("DELETE /api/v1/users/{id}", requireVerified(adminOnly(handler.DeleteUser(userRepo))))
	mux.Handle("GET /api/v1/users", requireVerified(adminOnly(handler.ListUsers(userRepo))))
	mux.Handle("POST /api/v1/users/email/send-verfication", handler.SendEmailVerficationOtp(cfg, userRepo, emailOtpRepo))
	mux.Handle("POST /api/v1/users/email/verify-verfication", handler.VerifyEmailVerficationOtp(cfg, userRepo, emailOtpRepo))
	mux.Handle("POST /api/v1/users/passsword/forgot/send-otp", handler.SendForgotPasswordLink(cfg, userRepo, emailOtpRepo))
	mux.Handle("POST /api/v1/users/password/forgot/verify-otp", handler.VerifyForgotPasswordLink(cfg, userRepo, emailOtpRepo))

	mux.Handle("POST /api/v1/auth/signup", handler.SignupUser(userRepo))
	mux.Handle("POST /api/v1/auth/login", handler.AuthenticateUser(cfg, cacheRepo, userRepo))
	mux.Handle("POST /api/v1/auth/refresh-token", handler.RefreshUserAccessToken(cfg, cacheRepo, userRepo))
	mux.Handle("POST /api/v1/auth/logout", authenticate(cfg, handler.LogoutUser(cfg, cacheRepo, userRepo)))

	mux.Handle("GET /api/v1/api-key", requireVerified(handler.ListAPIKey(apiKeyRepo)))
	mux.Handle("POST /api/v1/api-key", requireVerified(handler.CreateAPIKey(apiKeyRepo)))
	mux.Handle("DELETE /api/v1/api-key/{id}", requireVerified(handler.DeleteAPIKey(apiKeyRepo)))
	mux.Handle("POST /api/v1/api-key/valid", handler.VerifyAPIKey(apiKeyRepo))

}
