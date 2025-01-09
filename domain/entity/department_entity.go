package entity

import "time"

type Department struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	ManagerID int       `db:"manager_id" json:"manager_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
