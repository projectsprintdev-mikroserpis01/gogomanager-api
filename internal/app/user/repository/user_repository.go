package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/dto"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
)

type userRepository struct {
	conn *sqlx.DB
	tx   *sqlx.Tx
}

func NewUserRepository(conn *sqlx.DB) contracts.UserRepository {
	return &userRepository{
		conn: conn,
	}
}

func (r *userRepository) BeginTransaction(ctx context.Context) error {
	tx, err := r.conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	r.tx = tx

	return nil
}

func (r *userRepository) CommitTransaction() error {
	err := r.tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) RollbackTransaction() error {
	err := r.tx.Rollback()
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUsers(ctx context.Context, query dto.GetUsersQuery) ([]entity.User, error) {
	statement := `
		SELECT
			u.id AS "id",
			u.name AS "name",
			u.email AS "email",
			u.password AS "password",
			u.role_id AS "role_id",
			u.created_at AS "created_at",
			u.updated_at AS "updated_at",
			u.deleted_at AS "deleted_at",
			r.id AS "role.id",
			r.name AS "role.name"
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE 1=1`
	args := map[string]interface{}{}

	if !query.IncludeDeleted {
		statement += ` AND deleted_at IS NULL`
	}

	if query.Search != "" {
		statement += ` AND (name LIKE :search OR email LIKE :search)`
		args["search"] = "%" + query.Search + "%"
	}

	statement += ` ORDER BY ` + query.SortBy + ` ` + query.Order
	statement += ` LIMIT :limit OFFSET :offset`

	args["limit"] = query.Limit
	args["offset"] = query.Limit * (query.Page - 1)

	finalQuery, finalArgs, err := sqlx.Named(statement, args)
	if err != nil {
		return nil, err
	}

	finalQuery = r.conn.Rebind(finalQuery)

	users := make([]entity.User, 0)
	err = r.conn.SelectContext(ctx, &users, finalQuery, finalArgs...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetUserByField(ctx context.Context, field, value string) (*entity.User, error) {
	var user entity.User

	statement := `SELECT
			u.id AS "id",
			u.name AS "name",
			u.email AS "email",
			u.password AS "password",
			u.role_id AS "role_id",
			u.created_at AS "created_at",
			u.updated_at AS "updated_at",
			u.deleted_at AS "deleted_at",
			r.id AS "role.id",
			r.name AS "role.name"
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.` + field + ` = $1
		AND deleted_at IS NULL
		`

	err := r.conn.GetContext(ctx, &user, statement, value)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) (uuid.UUID, error) {
	_, err := r.conn.NamedExecContext(
		ctx,
		`INSERT INTO users (id, name, password, email) VALUES (:id, :name, :password, :email)`,
		user,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *entity.User) (uuid.UUID, error) {
	_, err := r.conn.NamedExecContext(
		ctx,
		`UPDATE users SET name = :name, password = :password, email = :email, updated_at = NOW() WHERE id = :id`,
		user,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}

func (r *userRepository) SoftDeleteUser(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	_, err := r.conn.ExecContext(ctx, `UPDATE users SET deleted_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	_, err := r.conn.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *userRepository) RestoreUser(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	_, err := r.conn.ExecContext(ctx, `UPDATE users SET deleted_at = NULL WHERE id = $1`, id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *userRepository) CountUsers(ctx context.Context, query dto.GetUsersStatsQuery) (int64, error) {
	statement := `SELECT COUNT(*) FROM users WHERE 1=1`
	args := map[string]interface{}{}

	if !query.IncludeDeleted {
		statement += ` AND deleted_at IS NULL`
	}

	if query.Search != "" {
		statement += ` AND (name LIKE :search OR email LIKE :search)`
		args["search"] = "%" + query.Search + "%"
	}

	finalQuery, finalArgs, err := sqlx.Named(statement, args)
	if err != nil {
		return 0, err
	}

	finalQuery = r.conn.Rebind(finalQuery)

	var count int64
	err = r.conn.GetContext(ctx, &count, finalQuery, finalArgs...)
	if err != nil {
		return 0, err
	}

	return count, nil
}
