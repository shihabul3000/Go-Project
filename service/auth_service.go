package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/shihabul3000/Go-Project/apperrors"
	"github.com/shihabul3000/Go-Project/auth"
	"github.com/shihabul3000/Go-Project/dto"
	"github.com/shihabul3000/Go-Project/models"
	"github.com/shihabul3000/Go-Project/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	users        repository.UserRepository
	jwtSecret    string
	jwtExpiresIn time.Duration
	bcryptCost   int
}

func NewAuthService(users repository.UserRepository, jwtSecret string, jwtExpiresIn time.Duration, bcryptCost int) AuthService {
	return &authService{
		users:        users,
		jwtSecret:    jwtSecret,
		jwtExpiresIn: jwtExpiresIn,
		bcryptCost:   bcryptCost,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	role := strings.ToLower(strings.TrimSpace(req.Role))
	if role == "" {
		role = models.RoleDriver
	}
	if role != models.RoleDriver && role != models.RoleAdmin {
		return nil, &apperrors.ValidationError{Fields: map[string]string{"role": "must be one of: driver admin"}}
	}

	if _, err := s.users.FindByEmail(ctx, email); err == nil {
		return nil, apperrors.ErrDuplicateEmail
	} else if !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.bcryptCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     strings.TrimSpace(req.Name),
		Email:    email,
		Password: string(hash),
		Role:     role,
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	response := toUserResponse(*user)
	return &response, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(user.ID, user.Role, s.jwtSecret, s.jwtExpiresIn)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.AuthUserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func toUserResponse(user models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
