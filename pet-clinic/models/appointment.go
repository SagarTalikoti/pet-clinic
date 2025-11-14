package models

type Appointment struct {
	ID     int    `json:"id"`
	Date   string `json:"date"`
	Time   string `json:"time"`
	PetID  int    `json:"pet_id"`
	Reason string `json:"reason"`
}
