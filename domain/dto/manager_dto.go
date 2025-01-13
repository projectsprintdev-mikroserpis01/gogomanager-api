package dto

type GetCurrentManagerRequest struct {
}

type GetCurrentManagerResponse struct {
	Email           string `db:"email" json:"email"`
	Name            string `db:"name" json:"name"`
	UserImageUri    string `db:"user_image_uri" json:"userImageUri"`
	CompanyName     string `db:"company_name" json:"companyName"`
	CompanyImageUri string `db:"company_image_uri" json:"companyImageUri"`
}
type UpdateManagerRequest struct {
	Email           *string `json:"email" validate:"email"`
	Name            *string `json:"name" validate:"min=4,max=100,ascii"`
	UserImageUri    *string `json:"userImageUri" validate:"url"`
	CompanyName     *string `json:"companyName" validate:"min=4,max=52,ascii"`
	CompanyImageUri *string `json:"companyImageUri" validate:"url"`
}

type UpdateManagerResponse struct {
	Email           *string `json:"email"`
	Name            *string `json:"name"`
	UserImageUri    *string `json:"userImageUri"`
	CompanyName     *string `json:"companyName"`
	CompanyImageUri *string `json:"companyImageUri"`
}
