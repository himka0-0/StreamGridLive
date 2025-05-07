package models

type EmailMessage struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
