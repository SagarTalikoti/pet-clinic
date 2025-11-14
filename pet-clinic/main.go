package main

import (
	"fmt"
	"log"
	"net/http"

	"pet-clinic/auth"
	"pet-clinic/db"
	"pet-clinic/handlers"
	"pet-clinic/utils"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize logger
	utils.InitLogger()
	utils.Log.Info("Pet Clinic API server starting...")

	// Connect to database
	db.Connect()

	// Router setup
	r := mux.NewRouter()

	// Public route
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	// Protected routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(auth.JWTMiddleware)

	// Owner routes
	api.HandleFunc("/owners", handlers.CreateOwner).Methods("POST")
	api.HandleFunc("/owners", handlers.GetOwners).Methods("GET")
	api.HandleFunc("/owners/{id}", handlers.UpdateOwner).Methods("PUT")
	api.HandleFunc("/owners/{id}", handlers.DeleteOwner).Methods("DELETE")

	// Pet routes
	api.HandleFunc("/pets", handlers.AddPet).Methods("POST")
	api.HandleFunc("/pets", handlers.GetPets).Methods("GET")
	api.HandleFunc("/pets/{id}", handlers.UpdatePet).Methods("PUT")
	api.HandleFunc("/pets/{id}", handlers.DeletePet).Methods("DELETE")

	// Appointment routes
	api.HandleFunc("/appointments", handlers.BookAppointment).Methods("POST")
	api.HandleFunc("/appointments", handlers.GetAppointments).Methods("GET")
	api.HandleFunc("/appointments/{id}", handlers.UpdateAppointment).Methods("PUT")
	api.HandleFunc("/appointments/{id}", handlers.DeleteAppointment).Methods("DELETE")

	// File management routes
	api.HandleFunc("/upload", handlers.UploadFile).Methods("POST")
	api.HandleFunc("/download/{filename}", handlers.DownloadFile).Methods("GET")

	fmt.Println(" Server running at http://localhost:9090")
	utils.Log.Info(" Server running at :9090")

	log.Fatal(http.ListenAndServe(":9090", r))
}
