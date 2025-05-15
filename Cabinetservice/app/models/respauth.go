package models

type Respauth struct {
	Email string `json:"email"`
	Valid bool   `json:"valid"`
	Name  string `json:"name"`
}

type CreateCabinet struct {
	Email  string `json:"email"`
	Secret string `json:"secret"`
}
