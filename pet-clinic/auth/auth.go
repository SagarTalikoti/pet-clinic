package auth

import (
	"net/http"
	"os"
	"pet-clinic/utils"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtKey []byte

func init() {
	godotenv.Load()
	jwtKey = []byte(os.Getenv("JWT_SECRET"))
}

// GenerateJWT creates a token for a username
func GenerateJWT(username string) (string, error) {
	utils.Log.WithField("username", username).Info("Generating JWT token")

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &jwt.MapClaims{
		"username": username,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtKey)

	if err != nil {
		utils.Log.WithError(err).Error("Failed signing JWT token")
		return "", err
	}

	utils.Log.WithField("username", username).Info("JWT token generated successfully")
	return signed, nil
}

// JWTMiddleware protects routes
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.Log.Warn("Missing Authorization header")
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims := &jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			utils.Log.WithError(err).Error("Invalid or expired JWT token")
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		utils.Log.WithField("path", r.URL.Path).Info("JWT token validated successfully")
		next.ServeHTTP(w, r)
	})
}
