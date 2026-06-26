package users

import "github.com/felipe/dev-test-api/internal/models"

func ToUserResponse(u models.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToUserListResponse(u models.User) UserListResponse {
	return UserListResponse{
		ID:        u.ID,
		Email:     u.Email,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
	}
}
