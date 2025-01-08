package enums

type RoleEnum string

const (
	SuperAdmin RoleEnum = "Superadmin"
	LeadAdmin  RoleEnum = "Lead Admin"
	Admin      RoleEnum = "Admin"
	User       RoleEnum = "User"
)

func (r RoleEnum) String() string {
	return string(r)
}
