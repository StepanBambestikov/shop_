package api

type LoginRequest struct {
	Login        *string `json:"login" validate:"omitempty,required_with=Password"`
	Password     *string `json:"password" validate:"omitempty,required_with=Login"`
	RefreshToken *string `json:"refreshToken" validate:"omitempty"`
}

type SignupRequest struct {
	Email     string `json:"email" validate:"email,required,max=80"`
	FirstName string `json:"first_name" validate:"omitempty,max=60"`
	LastName  string `json:"last_name" validate:"omitempty,max=60"`
	Username  string `json:"username" validate:"alphanumunicode,min=4,max=60,required"`
	Password  string `json:"password" validate:"required,min=8"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" validate:"jwt,required"`
}

type GiveRoleRequest struct {
	Role string `json:"role" validate:"required"`
}
