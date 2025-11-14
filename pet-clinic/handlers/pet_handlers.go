package handlers

import (
	"encoding/json"
	"net/http"
	"pet-clinic/db"
	"pet-clinic/models"
	"pet-clinic/utils"

	"github.com/gorilla/mux"
)

// Add Pet
func AddPet(w http.ResponseWriter, r *http.Request) {
	utils.Log.Info("POST /pets called")
	var p models.Pet

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		ErrorResponse(w, "Invalid pet input", http.StatusBadRequest, err)
		return
	}

	_, err := db.DB.Exec(`INSERT INTO pets (name, species, breed, owner_id, medical_history)
        VALUES ($1, $2, $3, $4, $5)`,
		p.Name, p.Species, p.Breed, p.OwnerID, p.MedicalHistory)

	if err != nil {
		ErrorResponse(w, "Failed to create pet", http.StatusInternalServerError, err)
		return
	}

	utils.Log.WithField("name", p.Name).Info("Pet added successfully")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Pet created"))
}

// View All Pets
func GetPets(w http.ResponseWriter, r *http.Request) {
	utils.Log.Info("GET /pets called")
	rows, err := db.DB.Query("SELECT id, name, species, breed, owner_id, medical_history FROM pets")
	if err != nil {
		ErrorResponse(w, "Failed to fetch pets", http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	var pets []models.Pet
	for rows.Next() {
		var p models.Pet
		rows.Scan(&p.ID, &p.Name, &p.Species, &p.Breed, &p.OwnerID, &p.MedicalHistory)
		pets = append(pets, p)
	}

	json.NewEncoder(w).Encode(pets)
}

// Update Pet
func UpdatePet(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var p models.Pet

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		ErrorResponse(w, "Invalid update body", http.StatusBadRequest, err)
		return
	}

	_, err := db.DB.Exec(`UPDATE pets SET name=$1, species=$2, breed=$3, owner_id=$4, medical_history=$5 WHERE id=$6`,
		p.Name, p.Species, p.Breed, p.OwnerID, p.MedicalHistory, id)

	if err != nil {
		ErrorResponse(w, "Failed to update pet", http.StatusInternalServerError, err)
		return
	}

	utils.Log.WithField("id", id).Info("Pet updated successfully")
	w.Write([]byte("Pet updated successfully"))
}

// Delete Pet
func DeletePet(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	_, err := db.DB.Exec("DELETE FROM pets WHERE id=$1", id)
	if err != nil {
		ErrorResponse(w, "Failed to delete pet", http.StatusInternalServerError, err)
		return
	}

	utils.Log.WithField("id", id).Warn("Pet deleted")
	w.Write([]byte("Pet deleted"))
}
