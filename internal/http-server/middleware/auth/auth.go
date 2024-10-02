package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/denisushakov/todo-rest/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := config.Password
		if pass == "" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		jwtToken, err := jwt.Parse(cookie.Value, func(jwtToken *jwt.Token) (interface{}, error) {
			return config.SecretKeyBytes, nil
		})
		if err != nil || !jwtToken.Valid {
			http.Error(w, `{"error": "invalid token"}`, http.StatusUnauthorized)
			return
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error": "invalid token"}`, http.StatusUnauthorized)
			return
		}

		hashRow, ok := claims["password_hash"]
		if !ok {
			http.Error(w, `{"error": "haven't hash password's"}`, http.StatusUnauthorized)
			return
		}

		hash, ok := hashRow.(string)
		if !ok {
			http.Error(w, `{"error": "haven't password"}`, http.StatusUnauthorized)
			return
		}

		newHashString := GetHashString(hash)

		if hash != newHashString {
			http.Error(w, `{"error": "authentification required"}`, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetHashString(val string) string {
	hash := sha256.Sum256([]byte(val))
	return hex.EncodeToString(hash[:])
}
