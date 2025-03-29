package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Abdulazizxoshimov/Hospital/entity"
	"github.com/Abdulazizxoshimov/Hospital/internal/repo/interfaces"
	postgres "github.com/Abdulazizxoshimov/Hospital/pkg/storage"
	"github.com/Masterminds/squirrel"
)

const (
	doctorTableName = "doctors"
)

type doctorRepo struct {
	db        *postgres.PostgresDB
	tableName string
}

func NewDoctorRepo(db *postgres.PostgresDB) interfaces.Doctor {
	return &doctorRepo{
		db:        db,
		tableName: doctorTableName,
	}
}

func (p *doctorRepo) Create(ctx context.Context, doctor *entity.Doctor) (*entity.Doctor, error) {
	extraInfoJSON, err := json.Marshal(doctor.ExtraInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal extra_info: %w", err)
	}

	doctorData := map[string]any{
		"id":             doctor.ID,
		"user_id":        doctor.UserID,
		"specialization": doctor.Specialization,
		"working_hours":  doctor.Working_hour,
		"extra_info":     extraInfoJSON,
	}

	query, args, err := p.db.Sq.Builder.Insert(p.tableName).SetMap(doctorData).ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, "doctor create")
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}

	return doctor, nil
}

func (p *doctorRepo) Get(ctx context.Context, doctorID string) (*entity.Doctor, error) {
	var doctor entity.Doctor
	var extraInfoJSON []byte

	query, args, err := p.db.Sq.Builder.
		Select("id", "user_id", "specialization", "working_hours", "extra_info").
		From(p.tableName).
		Where(p.db.Sq.Equal("id", doctorID)).
		ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, "doctor get")
	}

	err = p.db.QueryRow(ctx, query, args...).Scan(
		&doctor.ID,
		&doctor.UserID,
		&doctor.Specialization,
		&doctor.Working_hour,
		&extraInfoJSON,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("doctor not found")
		}
		return nil, p.db.Error(err)
	}

	if err := json.Unmarshal(extraInfoJSON, &doctor.ExtraInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal extra_info: %w", err)
	}

	return &doctor, nil
}

func (p *doctorRepo) Update(ctx context.Context, doctor *entity.Doctor) (*entity.Doctor, error) {
	extraInfoJSON, err := json.Marshal(doctor.ExtraInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal extra_info: %w", err)
	}

	clauses := map[string]any{
		"specialization": doctor.Specialization,
		"extra_info":     extraInfoJSON,
		"working_hours":  doctor.Working_hour,
	}

	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", doctor.ID)).
		ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, p.tableName+" update")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return nil, fmt.Errorf("no sql rows affected")
	}

	return doctor, nil
}

func (p *doctorRepo) Delete(ctx context.Context, doctorID string) error {
	sqlStr, args, err := p.db.Sq.Builder.
		Delete(p.tableName).
		Where(p.db.Sq.Equal("id", doctorID)).
		ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, p.tableName+" delete")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no sql rows affected")
	}

	return nil
}
func (p *doctorRepo) List(ctx context.Context, req *entity.ListRequest) (*entity.ListDoctorRes, error) {
	var doctors entity.ListDoctorRes

	queryBuilder := p.db.Sq.Builder.
		Select("id", "user_id", "specialization", "working_hours", "extra_info").
		From(p.tableName).
		PlaceholderFormat(squirrel.Dollar)

	if name, exists := req.Filter["name"]; exists {
		queryBuilder = queryBuilder.Where("LOWER(name) LIKE LOWER($1)", "%"+name+"%")
	}
	if specialization, exists := req.Filter["specialization"]; exists {
		queryBuilder = queryBuilder.Where("LOWER(specialization) LIKE LOWER($2)", "%"+specialization+"%")
	}

	if req.Limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(req.Offset))
	}

	queryBuilder = queryBuilder.OrderBy("created_at DESC")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s list", p.tableName))
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		var doctor entity.Doctor
		if err := rows.Scan(&doctor.ID, &doctor.UserID, &doctor.Specialization, &doctor.Working_hour, &doctor.ExtraInfo); err != nil {
			return nil, p.db.Error(err)
		}
		doctors.Doctors = append(doctors.Doctors, &doctor)
	}

	countQuery := p.db.Sq.Builder.
		Select("COUNT(*)").
		From(p.tableName)

	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s count", p.tableName))
	}

	if err := p.db.QueryRow(ctx, countSQL, countArgs...).Scan(&doctors.TotalCount); err != nil {
		doctors.TotalCount = 0
	}

	return &doctors, nil
}
