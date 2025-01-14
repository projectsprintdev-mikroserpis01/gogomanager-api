package dto

type EmployeeCreateReq struct {
	IdentityNumber   string `json:"identityNumber" validate:"min=5,max=33,required"`
	Name             string `json:"name" validate:"min=4,max=33,required"`
	EmployeeImageURI string `json:"employeeImageUri" validate:"required,url"`
	Gender           string `json:"gender" validate:"oneof=male female,required"`
	DepartmentID     string `json:"departmentId" validate:"required"`
}

type EmployeeDataRes struct {
	IdentityNumber   string `json:"identityNumber"`
	Name             string `json:"name"`
	EmployeeImageURI string `json:"employeeImageUri"`
	Gender           string `json:"gender"`
	DepartmentID     string `json:"departmentId"`
}

type EmployeeUpdateReq struct {
	IdentityNumber   string `json:"identityNumber,omitempty" validate:"min=5,max=33"`
	Name             string `json:"name,omitempty" validate:"min=4,max=33"`
	EmployeeImageURI string `json:"employeeImageUri,omitempty" validate:"url"`
	Gender           string `json:"gender,omitempty" validate:"oneof=male female"`
	DepartmentID     string `json:"departmentId" validate:"required"`
}
