package entities

type UserRole string

const (
	UserRoleNotverified UserRole = "notverified"
	UserRoleVerified             = "verified"
	UserRoleSeller               = "seller"
)

func (ur UserRole) String() string {
	return string(ur)
}

func UserRoleFromString(role string) (UserRole, bool) {
	switch role {
	case "notverified":
		return UserRoleNotverified, true
	case "verified":
		return UserRoleVerified, true
	case "seller":
		return UserRoleSeller, true
	}
	return UserRole(""), false
}

type User struct {
	ID            *string  `json:"id" validate:"omitempty,uuid4"`
	FirstName     *string  `json:"first_name" validate:"omitempty"`
	LastName      *string  `json:"last_name" validate:"omitempty"`
	Username      string   `json:"username" validate:"required,alphanumunicode"`
	Email         string   `json:"email" validate:"required,email"`
	EmailVerified bool     `json:"verified" validate:"required"`
	MarketRole    UserRole `json:"market_role" validate:"required"`
	Password      *string  `json:"-" validate:"gt=8"`
}
