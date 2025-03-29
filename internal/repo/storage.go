package repo

import (
	db "github.com/Abdulazizxoshimov/Hospital/pkg/storage"

	"github.com/Abdulazizxoshimov/Hospital/internal/repo/interfaces"
	"github.com/Abdulazizxoshimov/Hospital/internal/repo/postgres"
)

type StorageI interface{
	User() interfaces.User
	Doctor() interfaces.Doctor
	Appointment() interfaces.Appointment
}
type storagePg struct{
	user interfaces.User
	doctor interfaces.Doctor
	appointment interfaces.Appointment
}


func NewStoragePg(db *db.PostgresDB)StorageI{
	return &storagePg{
		user: postgres.NewUserRepo(db),
		doctor : postgres.NewDoctorRepo(db),
		appointment: postgres.NewAppointmentRepo(db),
	}
}

func (s *storagePg)User()interfaces.User{
	return s.user
}
func (s *storagePg)Doctor()interfaces.Doctor{
	return s.doctor
}
func (s *storagePg)Appointment()interfaces.Appointment{
	return s.appointment
}