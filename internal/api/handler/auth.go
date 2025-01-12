package handler

import (
	"encoding/json"
	"net/http"

	"auth/internal/models"
	"auth/internal/service"
	"auth/internal/utils"
)

func Register(ctx utils.MyContext, service service.AuthorizationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			utils.NewErrorResponse(ctx, w, "invalid JSON payload", http.StatusBadRequest)
			return
		}

		id, err := service.CreateUser(ctx, user)
		if err != nil {
			utils.NewErrorResponse(ctx, w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := utils.WriteResponse(w, http.StatusOK, map[string]interface{}{"id": id}); err != nil {
			utils.NewErrorResponse(ctx, w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func Login(ctx utils.MyContext, service service.AuthorizationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.LoginInput

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.NewErrorResponse(ctx, w, "invalid JSON payload", http.StatusBadRequest)
			return
		}

		userIP := utils.GetClientIP(r)
		if userIP == "" {
			utils.NewErrorResponse(ctx, w, "unable to determine user IP", http.StatusBadRequest)
			return
		}

		tokenPair, err := service.Login(ctx, input, userIP)
		if err != nil {
			utils.NewErrorResponse(ctx, w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = utils.WriteResponse(w, http.StatusOK, tokenPair); err != nil {
			utils.NewErrorResponse(ctx, w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func Refresh(ctx utils.MyContext, service service.AuthorizationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokenPair models.TokenPair
		if err := json.NewDecoder(r.Body).Decode(&tokenPair); err != nil {
			utils.NewErrorResponse(ctx, w, "invalid JSON payload", http.StatusBadRequest)
			return
		}

		if tokenPair.AccessToken == "" || tokenPair.RefreshToken == "" {
			utils.NewErrorResponse(ctx, w, "both access and refresh tokens are required", http.StatusBadRequest)
			return
		}

		userIP := utils.GetClientIP(r)
		if userIP == "" {
			utils.NewErrorResponse(ctx, w, "unable to determine user IP", http.StatusBadRequest)
			return
		}

		newTokenPair, err := service.RefreshTokens(ctx, tokenPair.AccessToken, tokenPair.RefreshToken, userIP)
		if err != nil {
			if err.Error() == "hash and token do not match" || err.Error() == "invalid access token" {
				utils.NewErrorResponse(ctx, w, err.Error(), http.StatusUnauthorized)
			} else {
				utils.NewErrorResponse(ctx, w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if err = utils.WriteResponse(w, http.StatusOK, newTokenPair); err != nil {
			utils.NewErrorResponse(ctx, w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}
