package model

type UserRequest struct {
	FirstName       string `json:"first_name" validate:"required,gt=0,lte=50"`
	LastName        string `json:"last_name" validate:"required,gt=0,lte=50"`
	Email           string `json:"email" validate:"required,gt=0,lte=225"`
	Password        string `json:"password" validate:"required,gt=0"`
	ConfirmPassword string `json:"confirm_password" validate:"required,gt=0"`
}

type UserResponse struct {
	ID             uint64  `json:"id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	Email          string  `json:"email"`
	HashedPassword []byte  `json:"-"`
	RoleID         uint64  `json:"role_id"`
	CreatedUser    uint64  `json:"created_user"`
	CreatedAt      string  `json:"created_at"`
	UpdatedUser    *uint64 `json:"updated_user"`
	UpdatedAt      *string `json:"updated_at"`
	DeletedUser    *uint64 `json:"deleted_user"`
	DeletedAt      *string `json:"-"`
}
