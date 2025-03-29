package entity

import "time"

type Appointment struct {
	ID               int64
	DoctorID         string
	UserID           string
	Appointment_time map[string]interface{}
	Status           string
}
type Availability struct {
	ID            int64
	DoctorID      string
	AvailableDate time.Time
	StartTime     time.Time
	EndTime       time.Time
	IsBooked      bool
}
type ListAppointments struct {
	Appointments []*Appointment
	TotalCount   int64
}
type ListAvailabilities struct {
	Availabilities []*Availability
	TotalCount     int64
}
