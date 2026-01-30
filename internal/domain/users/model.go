package users

import "time"

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	Role         Role
	Address      string
	Phone        string
	Bio          string
	CreatedAt    time.Time
}

type PublicUser struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateProfileInput struct {
	Name    *string
	Address *string
	Phone   *string
	Bio     *string
}

func (u User) Public() PublicUser {
	return PublicUser{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		Address:   u.Address,
		Phone:     u.Phone,
		Bio:       u.Bio,
		CreatedAt: u.CreatedAt,
	}
}
