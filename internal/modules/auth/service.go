package auth

import (
	"strconv"
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/internal/modules/users"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Setup(email, password string) (*AuthResponse, error)
	Login(email, password string) (*AuthResponse, error)
}

type authService struct {
	store     users.Store
	jwtSecret []byte
	expiryHrs int
}

func NewService(store users.Store, jwtSecret []byte, expiryHrsStr string) Service {
	hrs, _ := strconv.Atoi(expiryHrsStr)
	if hrs == 0 {
		hrs = 24
	}
	return &authService{
		store:     store,
		jwtSecret: jwtSecret,
		expiryHrs: hrs,
	}
}

func (s *authService) Setup(email, password string) (*AuthResponse, error) {
	count, err := s.store.Count()
	if err != nil {
		return nil, apierr.ErrInternal("Error al verificar los usuarios", "")
	}

	if count > 0 {
		return nil, apierr.ErrConflict(
			"System Already Initialized",
			"Ya existe un usuario administrador. Use /auth/login",
			"",
		)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apierr.ErrInternal("Error al generar el hash de la contraseña", "")
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hash),
		IsAdmin:      true,
	}

	if err := s.store.Create(user); err != nil {
		return nil, apierr.ErrInternal("Error al crear el usuario", "")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, apierr.ErrInternal("Error al generar el token", "")
	}

	return &AuthResponse{
		Token: token,
		User:  users.ToUserResponse(*user),
	}, nil
}

func (s *authService) Login(email, password string) (*AuthResponse, error) {
	user, err := s.store.FindByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrUnauthorized("Email o contraseña inválidos", "")
		}
		return nil, apierr.ErrInternal("Error al buscar el usuario", "")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, apierr.ErrUnauthorized("Email o contraseña inválidos", "")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, apierr.ErrInternal("Error al generar el token", "")
	}

	return &AuthResponse{
		Token: token,
		User:  users.ToUserResponse(*user),
	}, nil
}

func (s *authService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":      user.ID.String(),
		"email":    user.Email,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(time.Duration(s.expiryHrs) * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func GetUserID(c *jwt.MapClaims) (uuid.UUID, bool) {
	sub, ok := (*c)["sub"].(string)
	if !ok {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(sub)
	return id, err == nil
}

func IsAdmin(c *jwt.MapClaims) bool {
	isAdmin, ok := (*c)["is_admin"].(bool)
	return ok && isAdmin
}
