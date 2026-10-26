package auth

import "github.com/felipe/dev-test-api/internal/modules/users"

type StatusResponse struct {
	Initialized        bool `json:"initialized"`
	SetupRequiresToken bool `json:"setupRequiresToken"`
}

type AuthResponse struct {
	Token string             `json:"token"`
	User  users.UserResponse `json:"user"`
}
