package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Abdulazizxoshimov/Hospital/entity"
	"github.com/Abdulazizxoshimov/Hospital/internal/repo/interfaces"
	postgres "github.com/Abdulazizxoshimov/Hospital/pkg/storage"

	"github.com/Masterminds/squirrel"
)

const (
	userServiceTableName = "users"
)

type userRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewUserRepo(db *postgres.PostgresDB) interfaces.User {
	return &userRepo{
		tableName: userServiceTableName,
		db:        db,
	}
}

func (p *userRepo) usersSelectQueryPrefix() squirrel.SelectBuilder {
	return p.db.Sq.Builder.
		Select(
			"id",
			"full_name",
			"username",
			"email",
			"phone_number",
			"password",
			"role",
			"refresh_token",
			"created_at",
		).From(p.tableName)
}

func (p *userRepo) Create(ctx context.Context, user *entity.User) (*entity.User, error) {

	data := map[string]any{
		"id":            user.ID,
		"full_name":     user.FullName,
		"username":      user.UserName,
		"email":         user.Email,
		"phone_number":  user.PhoneNumber,
		"password":      user.Password,
		"role":          user.Role,
		"refresh_token": user.RefreshToken,
		"created_at":    user.CreatedAt,
	}
	query, args, err := p.db.Sq.Builder.Insert(p.tableName).SetMap(data).ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}

	return user, nil
}

func (p *userRepo) Update(ctx context.Context, user *entity.User) (*entity.User, error) {

	clauses := map[string]any{
		"username":     user.UserName,
		"email":        user.Email,
		"phone_number": user.PhoneNumber,
	}
	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", user.ID)).
		ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, p.tableName+" update")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return nil, p.db.Error(fmt.Errorf("no sql rows"))
	}

	return user, nil
}

func (p *userRepo) Delete(ctx context.Context, req *entity.DeleteRequest) error {
	sqlStr, args, err := p.db.Sq.Builder.
		Delete(p.tableName).
		Where(p.db.Sq.Equal("id", req.Id)).
		ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, p.tableName+" hard delete")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no rows deleted"))
	}

	return nil
}

func (p *userRepo) Get(ctx context.Context, params map[string]string) (*entity.User, error) {

	var (
		user entity.User
	)

	queryBuilder := p.usersSelectQueryPrefix()

	for key, value := range params {
		if key == "id" || key == "email" || key == "refresh_token" || key == "username" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
		}
	}
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "get"))
	}
	var (
		nullPhoneNumber sql.NullString
		nullFullName    sql.NullString
		nullRefresh     sql.NullString
	)

	if err = p.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&nullFullName,
		&user.UserName,
		&user.Email,
		&nullPhoneNumber,
		&user.Password,
		&user.Role,
		&nullRefresh,
		&user.CreatedAt,
	); err != nil {
		return nil, p.db.Error(err)
	}
	if nullPhoneNumber.Valid {
		user.PhoneNumber = nullPhoneNumber.String
	}
	if nullFullName.Valid {
		user.FullName = nullFullName.String
	}

	if nullRefresh.Valid {
		user.RefreshToken = nullRefresh.String
	}

	return &user, nil
}

func (p *userRepo) List(ctx context.Context, req *entity.ListRequest) (*entity.ListUserRes, error) {
	var users entity.ListUserRes

	queryBuilder := p.usersSelectQueryPrefix().PlaceholderFormat(squirrel.Dollar).OrderBy("created_at")

	if role, exists := req.Filter["role"]; exists && role != "" {
		queryBuilder = queryBuilder.Where(p.db.Sq.Equal("role", role))
	}

	if req.Limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(req.Offset))
	}

	countQuery := p.db.Sq.Builder.Select("*, COUNT(*) OVER() AS total_count").
		FromSelect(queryBuilder, "subquery")

	query, args, err := countQuery.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s list", p.tableName))
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()

	var totalCount int64

	for rows.Next() {
		var user entity.User
		var fullName, phoneNumber, refreshToken sql.NullString

		if err = rows.Scan(
			&user.ID,
			&fullName,
			&user.UserName,
			&user.Email,
			&phoneNumber,
			&user.Password,
			&user.Role,
			&refreshToken,
			&user.CreatedAt,
			&totalCount, 
		); err != nil {
			return nil, p.db.Error(err)
		}

		if fullName.Valid {
			user.FullName = fullName.String
		}
		if phoneNumber.Valid {
			user.PhoneNumber = phoneNumber.String
		}
		if refreshToken.Valid {
			user.RefreshToken = refreshToken.String
		}

		users.User = append(users.User, &user)
	}

	users.TotalCount = totalCount

	return &users, nil
}


func (p *userRepo) CheckUnique(ctx context.Context, filter *entity.GetRequest) (bool, error) {

	queryBuilder := p.db.Sq.Builder.Select("COUNT(1)").
		From(p.tableName)

		allowedFilters := []string{"email", "username", "phone_number"}
		for _, key := range allowedFilters {
			if value, exists := filter.Filter[key]; exists {
				queryBuilder = queryBuilder.Where(squirrel.Eq{key: value})
			}
		}


	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return false, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "isUnique"))
	}

	var count int
	err = p.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return true, p.db.Error(err)
	}
	if count != 0 {
		return true, nil
	}

	return false, nil
}

func (p *userRepo) UpdateRefresh(ctx context.Context, request *entity.UpdateRefresh) (*entity.Response, error) {

	clauses := map[string]any{
		"refresh_token": request.RefreshToken,
	}
	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", request.UserID)).
		ToSql()
	if err != nil {
		return &entity.Response{Status: false}, p.db.ErrSQLBuild(err, p.tableName+" update")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return &entity.Response{Status: false}, p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return &entity.Response{Status: false}, p.db.Error(fmt.Errorf("no sql rows"))
	}

	return &entity.Response{Status: true}, nil
}

func (p *userRepo) UpdatePassword(ctx context.Context, request *entity.UpdatePassword) (*entity.Response, error) {

	clauses := map[string]any{
		"password": request.NewPassword,
	}
	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", request.UserID)).
		ToSql()
	if err != nil {
		return &entity.Response{Status: false}, p.db.ErrSQLBuild(err, p.tableName+" update")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return &entity.Response{Status: false}, p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return &entity.Response{Status: false}, p.db.Error(fmt.Errorf("no sql rows"))
	}

	return &entity.Response{Status: true}, nil
}
