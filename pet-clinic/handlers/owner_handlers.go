package handlers

import (
	"encoding/json"
	"net/http"
	"pet-clinic/db"
	"pet-clinic/models"
	"pet-clinic/utils"
	"strings"

	"github.com/gorilla/mux"
)

// CreateOwner - Adds a new owner with validation & logging
func CreateOwner(w http.ResponseWriter, r *http.Request) {
	utils.Log.Debug("POST /owners called")

	var o models.Owner
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		ErrorResponse(w, "Invalid JSON input", http.StatusBadRequest, err)
		return
	}

	o.Name = strings.TrimSpace(o.Name)
	o.Contact = strings.TrimSpace(o.Contact)
	o.Email = strings.TrimSpace(o.Email)

	if o.Name == "" || o.Email == "" {
		utils.Log.Warn("Owner creation failed: missing name or email")
		http.Error(w, "Name and Email are required fields", http.StatusBadRequest)
		return
	}

	err := db.DB.QueryRow(
		`INSERT INTO owners (name, contact, email)
		 VALUES ($1, $2, $3) RETURNING id`,
		o.Name, o.Contact, o.Email).Scan(&o.ID)

	if err != nil {
		ErrorResponse(w, "Failed to create owner", http.StatusInternalServerError, err)
		return
	}

	utils.Log.WithFields(map[string]interface{}{
		"id":    o.ID,
		"name":  o.Name,
		"email": o.Email,
	}).Info("Owner created successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(o)
}

// GetOwners - Fetch all owners
func GetOwners(w http.ResponseWriter, r *http.Request) {
	utils.Log.Debug("GET /owners called")

	rows, err := db.DB.Query("SELECT id, name, contact, email FROM owners")
	if err != nil {
		ErrorResponse(w, "Failed to fetch owners", http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	var owners []models.Owner
	for rows.Next() {
		var o models.Owner
		if err := rows.Scan(&o.ID, &o.Name, &o.Contact, &o.Email); err != nil {
			ErrorResponse(w, "Error scanning owner data", http.StatusInternalServerError, err)
			return
		}
		owners = append(owners, o)
	}

	utils.Log.WithField("count", len(owners)).Info("Owners fetched successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(owners)
}

// UpdateOwner - Update an owner by ID
func UpdateOwner(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	utils.Log.WithField("id", id).Debug("PUT /owners/{id} called")

	var o models.Owner
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		ErrorResponse(w, "Invalid JSON input", http.StatusBadRequest, err)
		return
	}

	o.Name = strings.TrimSpace(o.Name)
	o.Email = strings.TrimSpace(o.Email)
	if o.Name == "" || o.Email == "" {
		utils.Log.Warn("Owner update failed: missing name or email")
		http.Error(w, "Name and Email are required fields", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec(`UPDATE owners SET name=$1, contact=$2, email=$3 WHERE id=$4`,
		o.Name, o.Contact, o.Email, id)

	if err != nil {
		ErrorResponse(w, "Failed to update owner", http.StatusInternalServerError, err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.Log.WithField("id", id).Warn("No owner found to update")
		http.Error(w, "Owner not found", http.StatusNotFound)
		return
	}

	utils.Log.WithField("id", id).Info("Owner updated successfully")
	w.Write([]byte("Owner updated successfully"))
}

// DeleteOwner - Delete an owner by ID
func DeleteOwner(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	utils.Log.WithField("id", id).Debug("DELETE /owners/{id} called")

	result, err := db.DB.Exec("DELETE FROM owners WHERE id=$1", id)
	if err != nil {
		ErrorResponse(w, "Failed to delete owner", http.StatusInternalServerError, err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.Log.WithField("id", id).Warn("No owner found to delete")
		http.Error(w, "Owner not found", http.StatusNotFound)
		return
	}

	utils.Log.WithField("id", id).Warn("Owner deleted successfully") // 🟡 log deletion as WARN
	w.Write([]byte("Owner deleted successfully"))
}
