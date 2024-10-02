package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/denisushakov/todo-rest/internal/config"
	"github.com/denisushakov/todo-rest/internal/http-server/middleware/auth"
	"github.com/denisushakov/todo-rest/pkg/models"
	"github.com/golang-jwt/jwt/v5"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var auth models.Auth

	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if auth.Password != config.Password {
		http.Error(w, `{"error": "wrong password"}`, http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken(auth.Password)
	if err != nil {
		http.Error(w, `{"error": "token invalid"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"token": token})
}

func GenerateToken(password string) (string, error) {
	hashString := auth.GetHashString(password)

	claims := jwt.MapClaims{
		"password_hash": hashString,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := jwtToken.SignedString(config.SecretKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to sign jwt: %w", err)
	}
	return signedToken, nil
}
