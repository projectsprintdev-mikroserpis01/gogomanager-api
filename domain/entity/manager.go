package entity

import "time"

type Manager struct {
	ID              int       `db:"id"`
	Email           string    `db:"email"`
	Password        string    `db:"password"`
	Name            *string   `db:"name"`
	UserImageURI    *string   `db:"user_image_uri"`
	CompanyName     *string   `db:"company_name"`
	CompanyImageURI *string   `db:"company_image_uri"`
	CreatedAt       time.Time `db:"created_at"`
}
