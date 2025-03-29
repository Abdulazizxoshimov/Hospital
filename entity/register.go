package entity

import "time"

type UserRegister struct {
	Email    string
	Username string
	Password string
}

type UserResponse struct {
	ID           string
	FullName     string
	UserName     string
	Role         string
	Email        string
	PhoneNumber  string
	RefreshToken string
	AccesToken   string
	CreatedAt    time.Time
}
type UserCreateResponse struct {
	ID string
}
type UserUpdate struct {
	ID          string
	FullName    string
	UserName    string
	Role        string
	Email       string
	PhoneNumber string
}
