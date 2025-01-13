package entity

import "time"

type Employee struct {
	ID               int       `db:"id" json:"id"`
	IdentityNumber   string    `db:"identity_number" json:"identity_number"`
	Name             string    `db:"name" json:"name"`
	EmployeeImageURI string    `db:"employee_image_uri" json:"employee_image_uri"`
	Gender           string    `db:"gender" json:"gender"`
	DepartmentID     int       `db:"department_id" json:"department_id"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}
