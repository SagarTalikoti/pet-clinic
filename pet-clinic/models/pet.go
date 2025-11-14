package models

type Pet struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Species        string `json:"species"`
	Breed          string `json:"breed"`
	OwnerID        int    `json:"owner_id"`
	MedicalHistory string `json:"medical_history"`
}
