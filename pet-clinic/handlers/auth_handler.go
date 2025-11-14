package handlers

import (
	"encoding/json"
	"net/http"
	"pet-clinic/auth"
	"pet-clinic/utils"
)

// User struct represents login credentials
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Demo user (for testing)
var demoUser = User{
	Username: "admin",
	Password: "1234",
}

// Login authenticates a user and issues a JWT token
func Login(w http.ResponseWriter, r *http.Request) {
	utils.Log.Debug("Received login request")

	var username, password string

	// 🔹 Check Basic Auth first
	u, p, ok := r.BasicAuth()
	if ok {
		username = u
		password = p
		utils.Log.WithField("username", username).Debug("Attempting login via Basic Auth")
	} else {
		// 🔹 Fallback to JSON body
		var creds User
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			utils.Log.WithError(err).Warn("Invalid JSON in login request")
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}
		username = creds.Username
		password = creds.Password
		utils.Log.WithField("username", username).Debug("Attempting login via JSON body")
	}

	// 🔹 Validate credentials
	if username != demoUser.Username || password != demoUser.Password {
		utils.Log.WithField("username", username).Warn("Invalid login attempt (wrong username or password)")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 🔹 Generate JWT token
	token, err := auth.GenerateJWT(username)
	if err != nil {
		utils.Log.WithError(err).Error("Failed to generate JWT token")
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	utils.Log.WithField("username", username).Info("User logged in successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
