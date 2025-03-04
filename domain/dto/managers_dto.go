package dto

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32,ascii"`
	Action   string `json:"action" validate:"required,oneof=create login"`
}

type AuthResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type ManagerProfile struct {
	Email           string `json:"email"`
	Name            string `json:"name"`
	UserImageUri    string `json:"userImageUri"`
	CompanyName     string `json:"companyName"`
	CompanyImageUri string `json:"companyImageUri"`
}
