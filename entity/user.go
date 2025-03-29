package entity

import "time"

type User struct{
	ID string
	FullName string
	UserName string
	Password string
	Role string
	Email string
	PhoneNumber string
	RefreshToken string
	CreatedAt time.Time
}

type Doctor struct{
	ID string	
	UserID string
	Specialization string
	Working_hour string
	ExtraInfo map[string]interface{}
}

type Response struct {
	Status bool
}
type UpdateRefresh struct {
	UserID       string
	RefreshToken string
}

type UpdatePassword struct {
	UserID      string
	NewPassword string
}

type IsUnique struct {
	Email string
}
type ListDoctorRes struct{
	Doctors []*Doctor
	TotalCount int64
}

type ListUserRes struct {
	User       []*User
	TotalCount int64
}
type ListRequest struct{
   Limit int
   Offset int
   Filter map[string]string
}
type DeleteRequest struct{
	Id string
	DeletedAt time.Time	
}

type GetRequest struct {
	Filter map[string]string
}

