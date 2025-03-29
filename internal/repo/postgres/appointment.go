package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Abdulazizxoshimov/Hospital/entity"
	"github.com/Abdulazizxoshimov/Hospital/internal/repo/interfaces"
	postgres "github.com/Abdulazizxoshimov/Hospital/pkg/storage"
	"github.com/k0kubun/pp"
)

const (
	tableNameAppointment  = "appointments"
	tableNameAvailability = "doctor_availability"
)

type appointmentRepo struct {
	db                    postgres.PostgresDB
	tableNameAppointment  string
	tableNameAvailability string
}

func NewAppointmentRepo(db *postgres.PostgresDB) interfaces.Appointment {
	return &appointmentRepo{
		db:                    *db,
		tableNameAppointment:  tableNameAppointment,
		tableNameAvailability: tableNameAvailability,
	}
}

func (p *appointmentRepo) CreateAppointment(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	startTimeStr, ok := appointment.Appointment_time["start_time"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid appointment time format, expected string")
	}
	pp.Println(startTimeStr)

	startTime, err := time.Parse("2006-01-02T15:04:05Z", startTimeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid time format: %v", err)
	}

	appointmentEnd := startTime.Add(time.Hour)
	pp.Println("hellolololol")
	var count int
	query, args, err := p.db.Sq.Builder.
		Select("COUNT(*)").
		From(p.tableNameAvailability).
		Where("doctor_id = ?", appointment.DoctorID).
		Where("available_date = ?", startTime.Format("2006-01-02")).       // YYYY-MM-DD
		Where("start_time <= ?", startTime.Format("2006-01-02 15:04:05")). // YYYY-MM-DD HH:MM:SS
		Where("end_time > ?", startTime.Format("2006-01-02 15:04:05")).
		Where("is_booked = ?", false).
		ToSql()

	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return nil, err
	}

	if count == 1 {
		return nil, fmt.Errorf("doctor is not available at this time")
	}
	pp.Println("hellolololol")

	appointmentTimes := map[string]interface{}{
		"available_date": startTime.Format("2006-01-02"),
		"start_time":     startTime.Format("15:04:05"),
		"end_time":       appointmentEnd.Format("15:04:05"),
		"is_booked":      true,
	}

	appointmentTimesJSON, err := json.Marshal(appointmentTimes)
	if err != nil {
		return nil, err
	}

	data := map[string]any{
		"doctor_id":        appointment.DoctorID,
		"patient_id":       appointment.UserID,
		"appointment_time": json.RawMessage(appointmentTimesJSON),
		"start_time":       startTime.Format("2006-01-02 15:04:05"),
		"status":           "scheduled",
	}

	query, args, err = p.db.Sq.Builder.Insert(p.tableNameAppointment).SetMap(data).ToSql()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	updateQuery, updateArgs, err := p.db.Sq.Builder.
		Update(p.tableNameAvailability).
		Set("is_booked", true).
		Where("doctor_id = ?", appointment.DoctorID).
		Where("available_date = ?", startTime.Format("2006-01-02")).
		Where("start_time = ?", startTime.Format("2006-01-02 15:04:05")).
		ToSql()
	if err != nil {
		return nil, err
	}


	_, err = tx.Exec(ctx, updateQuery, updateArgs...)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return appointment, nil
}

func (p *appointmentRepo) UpdateAppointment(ctx context.Context, appointmentID int, newTime time.Time) error {
	if time.Until(newTime) < 24*time.Hour {
		return fmt.Errorf("cannot update appointment within 24 hours of the scheduled time")
	}

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var oldAppointmentTimeJSON string
	query, args, err := p.db.Sq.Builder.
		Select("appointment_time").
		From(tableNameAppointment).
		Where("id = ?", appointmentID).
		ToSql()
	if err != nil {
		return err
	}

	err = tx.QueryRow(ctx, query, args...).Scan(&oldAppointmentTimeJSON)
	if err != nil {
		return err
	}

	var oldAvailability entity.Availability
	err = json.Unmarshal([]byte(oldAppointmentTimeJSON), &oldAvailability)
	if err != nil {
		return err
	}

	updateQuery, updateArgs, err := p.db.Sq.Builder.
		Update(tableNameAvailability).
		Set("is_booked", false).
		Where("available_date = ?", oldAvailability.AvailableDate).
		Where("start_time = ?", oldAvailability.StartTime).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, updateQuery, updateArgs...)
	if err != nil {
		return err
	}

	newAvailability := entity.Availability{
		AvailableDate: newTime,
		StartTime:     newTime,
		EndTime:       newTime.Add(time.Hour),
		IsBooked:      true,
	}

	newAvailabilityJSON, err := json.Marshal(newAvailability)
	if err != nil {
		return err
	}

	updateQuery, updateArgs, err = p.db.Sq.Builder.
		Update(tableNameAppointment).
		Set("appointment_time", string(newAvailabilityJSON)).
		Where("id = ?", appointmentID).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, updateQuery, updateArgs...)
	if err != nil {
		return err
	}

	updateQuery, updateArgs, err = p.db.Sq.Builder.
		Update(tableNameAvailability).
		Set("is_booked", true).
		Where("available_date = ?", newAvailability.AvailableDate).
		Where("start_time = ?", newAvailability.StartTime).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, updateQuery, updateArgs...)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *appointmentRepo) DeleteAppointment(ctx context.Context, appointmentID int) error {
	var appointmentTimeJSON string
	query, args, err := p.db.Sq.Builder.
		Select("appointment_time").
		From(tableNameAppointment).
		Where("id = ?", appointmentID).
		ToSql()
	if err != nil {
		return err
	}

	err = p.db.QueryRow(ctx, query, args...).Scan(&appointmentTimeJSON)
	if err != nil {
		return err
	}

	var appointmentTime entity.Availability
	err = json.Unmarshal([]byte(appointmentTimeJSON), &appointmentTime)
	if err != nil {
		return err
	}
	if time.Until(appointmentTime.StartTime) < 24*time.Hour {
		return fmt.Errorf("cannot cancel an appointment within 24 hours of the scheduled time")
	}

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query, args, err = p.db.Sq.Builder.
		Delete(tableNameAppointment).
		Where("id = ?", appointmentID).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	updateQuery, updateArgs, err := p.db.Sq.Builder.
		Update(tableNameAvailability).
		Set("is_booked", false).
		Where("available_date = ?", appointmentTime.AvailableDate).
		Where("start_time = ?", appointmentTime.StartTime).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, updateQuery, updateArgs...)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
func (p *appointmentRepo) GetAppointment(ctx context.Context, appointmentID int) (*entity.Appointment, error) {
	query, args, err := p.db.Sq.Builder.
		Select("*").
		From(p.tableNameAppointment).
		Where("id = ?", appointmentID).
		ToSql()
	if err != nil {
		return nil, err
	}

	var appointment entity.Appointment
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&appointment.ID,
		&appointment.DoctorID,
		&appointment.UserID,
		&appointment.Appointment_time,
		&appointment.Status,
	)
	if err != nil {
		return nil, err
	}

	return &appointment, nil
}

func (p *appointmentRepo) ListAppointments(ctx context.Context, page, limit int) ([]*entity.Appointment, int, error) {
	offset := (page - 1) * limit

	query, args, err := p.db.Sq.Builder.
		Select("*").
		From(p.tableNameAppointment).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, err
	}

	totalQuery, _, err := p.db.Sq.Builder.
		Select("COUNT(*)").
		From(p.tableNameAppointment).
		ToSql()
	if err != nil {
		return nil, 0, err
	}

	var total int
	err = p.db.QueryRow(ctx, totalQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var appointments []*entity.Appointment
	for rows.Next() {
		var appointment entity.Appointment
		if err := rows.Scan(
			&appointment.ID,
			&appointment.DoctorID,
			&appointment.UserID,
			&appointment.Appointment_time,
			&appointment.Status,
		); err != nil {
			return nil, 0, err
		}
		appointments = append(appointments, &appointment)
	}

	return appointments, total, nil
}

func (p *appointmentRepo) GetAvailability(ctx context.Context, availabilityID int) (*entity.Availability, error) {
	query, args, err := p.db.Sq.Builder.
		Select("*").
		From(p.tableNameAvailability).
		Where("id = ?", availabilityID).
		ToSql()
	if err != nil {
		return nil, err
	}

	var availability entity.Availability
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&availability.AvailableDate,
		&availability.StartTime,
		&availability.EndTime,
		&availability.IsBooked,
	)
	if err != nil {
		return nil, err
	}

	return &availability, nil
}

func (p *appointmentRepo) ListAvailabilities(ctx context.Context, page, limit int) ([]*entity.Availability, int, error) {
	offset := (page - 1) * limit

	query, args, err := p.db.Sq.Builder.
		Select("*").
		From(p.tableNameAvailability).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, err
	}

	totalQuery, _, err := p.db.Sq.Builder.
		Select("COUNT(*)").
		From(p.tableNameAvailability).
		ToSql()
	if err != nil {
		return nil, 0, err
	}

	var total int
	err = p.db.QueryRow(ctx, totalQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var availabilities []*entity.Availability
	for rows.Next() {
		var availability entity.Availability
		if err := rows.Scan(
			&availability.AvailableDate,
			&availability.StartTime,
			&availability.EndTime,
			&availability.IsBooked,
		); err != nil {
			return nil, 0, err
		}
		availabilities = append(availabilities, &availability)
	}

	return availabilities, total, nil
}
