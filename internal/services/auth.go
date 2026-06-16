package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"rest-api-disbursement-system/internal/models"
	"rest-api-disbursement-system/internal/repository"
)

var ErrInvalidCredentials = errors.New("username atau password salah")

type AuthService interface {
	Login(input LoginInput) (*AuthResponse, error)
}

type LoginInput struct {
	Username string
	Password string
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type authService struct {
	users             repository.UserRepository
	jwtSecret         string
	jwtExpiresInHours int
}

func NewAuthService(users repository.UserRepository, jwtSecret string, jwtExpiresInHours int) AuthService {
	return &authService{
		users:             users,
		jwtSecret:         jwtSecret,
		jwtExpiresInHours: jwtExpiresInHours,
	}
}

func (s *authService) Login(input LoginInput) (*AuthResponse, error) {
	user, err := s.users.FindByUsername(input.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.buildAuthResponse(user)
}

func (s *authService) buildAuthResponse(user *models.User) (*AuthResponse, error) {
	token, err := GenerateJWT(s.jwtSecret, s.jwtExpiresInHours, user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  toUserResponse(user),
	}, nil
}

func GenerateJWT(secret string, expiresInHours int, user *models.User) (string, error) {
	if expiresInHours <= 0 {
		expiresInHours = 24
	}

	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresInHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func toUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}
