package handlers

import (
	"encoding/json"
	"net/http"
	"pet-clinic/db"
	"pet-clinic/models"
	"pet-clinic/utils"

	"github.com/gorilla/mux"
)

// Book Appointment
func BookAppointment(w http.ResponseWriter, r *http.Request) {
	var a models.Appointment
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		ErrorResponse(w, "Invalid appointment input", http.StatusBadRequest, err)
		return
	}

	_, err := db.DB.Exec(`INSERT INTO appointments (date, time, pet_id, reason)
		VALUES ($1, $2, $3, $4)`,
		a.Date, a.Time, a.PetID, a.Reason)

	if err != nil {
		ErrorResponse(w, "Appointment booking failed", http.StatusInternalServerError, err)
		return
	}

	utils.Log.Info("Appointment booked successfully")
	w.Write([]byte("Appointment created"))
}

// Get Appointments
func GetAppointments(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, date, time, pet_id, reason FROM appointments")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var appts []models.Appointment
	for rows.Next() {
		var a models.Appointment
		rows.Scan(&a.ID, &a.Date, &a.Time, &a.PetID, &a.Reason)
		appts = append(appts, a)
	}
	json.NewEncoder(w).Encode(appts)
}

// Update Appointment
func UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var a models.Appointment
	json.NewDecoder(r.Body).Decode(&a)

	_, err := db.DB.Exec(`UPDATE appointments SET date=$1, time=$2, pet_id=$3, reason=$4 WHERE id=$5`,
		a.Date, a.Time, a.PetID, a.Reason, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Appointment updated"))
}

// Cancel Appointment
func DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := db.DB.Exec("DELETE FROM appointments WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(" Appointment cancelled"))
}
