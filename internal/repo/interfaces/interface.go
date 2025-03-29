package interfaces

import (
	"context"
	"time"

	"github.com/Abdulazizxoshimov/Hospital/entity"
)

type User interface {
	Create(ctx context.Context, req *entity.User) (*entity.User, error)
	Get(ctx context.Context, params map[string]string) (*entity.User, error)
	Update(ctx context.Context, req *entity.User) (*entity.User, error)
	List(ctx context.Context, req *entity.ListRequest) (*entity.ListUserRes, error)
	Delete(ctx context.Context, Filter *entity.DeleteRequest) error
	CheckUnique(ctx context.Context, filter *entity.GetRequest) (bool, error)
	UpdateRefresh(ctx context.Context, request *entity.UpdateRefresh) (*entity.Response, error)
	UpdatePassword(ctx context.Context, request *entity.UpdatePassword) (*entity.Response, error)
}

type Doctor interface {
	Create(context.Context, *entity.Doctor) (*entity.Doctor, error)
	Get(context.Context, string) (*entity.Doctor, error)
	Update(context.Context, *entity.Doctor) (*entity.Doctor, error)
	Delete(context.Context, string) error
	List(context.Context, *entity.ListRequest) (*entity.ListDoctorRes, error)
}

type Appointment interface {
	CreateAppointment(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error)
	UpdateAppointment(ctx context.Context, appointmentID int, newTime time.Time) error
	DeleteAppointment(ctx context.Context, appointmentID int) error
	GetAppointment(ctx context.Context, appointmentID int) (*entity.Appointment, error)
	ListAppointments(ctx context.Context, page, limit int) ([]*entity.Appointment, int, error)
	GetAvailability(ctx context.Context, availabilityID int) (*entity.Availability, error)
	ListAvailabilities(ctx context.Context, page, limit int) ([]*entity.Availability, int, error)
}
