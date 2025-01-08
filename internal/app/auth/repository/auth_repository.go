package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/contracts"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/domain/entity"
)

type authRepository struct {
	conn *sqlx.DB
	tx   *sqlx.Tx
}

func NewAuthRepository(conn *sqlx.DB) contracts.AuthRepository {
	return &authRepository{
		conn: conn,
	}
}

func (a *authRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := a.conn.GetContext(ctx, &user, `
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
		WHERE email = $1
		AND deleted_at IS NULL
	`, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *authRepository) RegisterUser(ctx context.Context, user entity.User) (uuid.UUID, error) {
	_, err := a.conn.ExecContext(
		ctx,
		"INSERT INTO users (id, email, password, name) VALUES ($1, $2, $3, $4)",
		user.ID,
		user.Email,
		user.Password,
		user.Name,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}
