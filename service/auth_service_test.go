package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shihabul3000/Go-Project/apperrors"
	"github.com/shihabul3000/Go-Project/dto"
	"github.com/shihabul3000/Go-Project/models"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	byEmail map[string]*models.User
	nextID  uint
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		byEmail: map[string]*models.User{},
		nextID:  1,
	}
}

func (r *fakeUserRepository) Create(ctx context.Context, user *models.User) error {
	if _, exists := r.byEmail[user.Email]; exists {
		return apperrors.ErrDuplicateEmail
	}
	user.ID = r.nextID
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt
	r.nextID++
	copied := *user
	r.byEmail[user.Email] = &copied
	return nil
}

func (r *fakeUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user, ok := r.byEmail[email]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	copied := *user
	return &copied, nil
}

func (r *fakeUserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	for _, user := range r.byEmail {
		if user.ID == id {
			copied := *user
			return &copied, nil
		}
	}
	return nil, apperrors.ErrNotFound
}

func TestAuthServiceRegisterHashesPasswordAndDefaultsRole(t *testing.T) {
	users := newFakeUserRepository()
	service := NewAuthService(users, "test-secret", time.Hour, bcrypt.MinCost)

	response, err := service.Register(context.Background(), dto.RegisterRequest{
		Name:     "John Doe",
		Email:    "JOHN.DOE@SPOTSYNC.COM",
		Password: "securePassword123",
	})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if response.Role != models.RoleDriver {
		t.Fatalf("expected default role %q, got %q", models.RoleDriver, response.Role)
	}
	if response.Email != "john.doe@spotsync.com" {
		t.Fatalf("expected normalized email, got %q", response.Email)
	}

	stored := users.byEmail["john.doe@spotsync.com"]
	if stored.Password == "securePassword123" {
		t.Fatal("password was stored in plain text")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte("securePassword123")); err != nil {
		t.Fatalf("stored password is not a bcrypt hash: %v", err)
	}
}

func TestAuthServiceLoginRejectsWrongPassword(t *testing.T) {
	users := newFakeUserRepository()
	service := NewAuthService(users, "test-secret", time.Hour, bcrypt.MinCost)
	_, err := service.Register(context.Background(), dto.RegisterRequest{
		Name:     "Jane Driver",
		Email:    "jane@spotsync.com",
		Password: "securePassword123",
		Role:     models.RoleDriver,
	})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	_, err = service.Login(context.Background(), dto.LoginRequest{
		Email:    "jane@spotsync.com",
		Password: "wrong-password",
	})
	if !errors.Is(err, apperrors.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
