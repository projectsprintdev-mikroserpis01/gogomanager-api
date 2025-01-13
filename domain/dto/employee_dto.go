package dto

type EmployeeCreateReq struct {
	IdentityNumber   string `json:"identity_number" validate:"min=5,max=33,required"`
	Name             string `json:"name" validate:"min=4,max=33,required"`
	EmployeeImageURI string `json:"employee_image_uri" validate:"required,url"`
	Gender           string `json:"gender" validate:"oneof=male female,required"`
	DepartmentID     string `json:"department_id" validate:"required"`
}

type EmployeeDataRes struct {
	IdentityNumber   string `json:"identity_number"`
	Name             string `json:"name"`
	EmployeeImageURI string `json:"employee_image_uri"`
	Gender           string `json:"gender"`
	DepartmentID     string `json:"department_id"`
}

type EmployeeUpdateReq struct {
	IdentityNumber   string `json:"identity_number,omitempty" validate:"min=5,max=33"`
	Name             string `json:"name,omitempty" validate:"min=4,max=33"`
	EmployeeImageURI string `json:"employee_image_uri,omitempty" validate:"url"`
	Gender           string `json:"gender,omitempty" validate:"oneof=male female"`
	DepartmentID     string `json:"department_id,omitempty"`
}
