package entity

import "time"

type Department struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ManagerID int    `json:"manager_id"`
	CreatedAt time.Time
}
