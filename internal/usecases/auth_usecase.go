package usecases

import (
	"dungeons-dragon-service/internal/domain/model"
	"dungeons-dragon-service/internal/domain/repository"
	"dungeons-dragon-service/internal/dto"
	"dungeons-dragon-service/internal/helper"
	"dungeons-dragon-service/internal/http/custom"
	"dungeons-dragon-service/internal/infrastructure/jwt"
	"strings"
	"time"
)

type AuthUseCase interface {
	Register(username, email, password string) (*dto.LoginResponse, error)
	Login(username, password string) (*dto.LoginResponse, error)
}

type authUseCase struct {
	users     repository.UserRepository
	jwtSecret string
}

func NewAuthUsecase(users repository.UserRepository, jwtSecret string) AuthUseCase {
	return &authUseCase{users: users, jwtSecret: jwtSecret}
}

func (u *authUseCase) Register(username, email, password string) (*dto.LoginResponse, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || len(password) < 6 {
		return nil, custom.NewBadRequestError("invalid email or password")
	}
	// Check if user already exists
	existingUser, err := u.users.FindByEmail(email)
	if err != nil {
		return nil, custom.NewUnexpectedError("failed to check if email exists")
	}
	if existingUser != nil {
		return nil, custom.NewConflictError("email already exists")
	}
	// check username
	existingUser, err = u.users.FindByUsername(username)
	if err != nil {
		return nil, custom.NewUnexpectedError("failed to check if username exists")
	}
	if existingUser != nil {
		return nil, custom.NewConflictError("username already exists")
	}

	salt, err := helper.GenerateSalt(16)
	if err != nil {
		return nil, custom.NewUnexpectedError("failed to generate salt")
	}
	hash := helper.HashPasswordArgon2(password, salt)
	user, err := u.users.Create(&model.User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		Role:         model.RoleUser,
	})
	if err != nil {
		return nil, custom.NewUnexpectedError("failed to create user")
	}
	token, err := jwt.GenerateToken(u.jwtSecret, user.ID.String(), string(model.RoleUser), 24*time.Hour)
	if err != nil {
		return nil, custom.NewUnexpectedError("failed to generate token")
	}
	return &dto.LoginResponse{Token: token}, nil
}

func (u *authUseCase) Login(username, password string) (*dto.LoginResponse, error) {
	user, err := u.users.FindByUsername(username)
	if err != nil || user == nil {
		return nil, custom.NewNotFoundError("user not found")
	}
	// Use Argon2 password verification
	if !helper.VerifyPasswordArgon2(password, user.PasswordHash) {
		return nil, custom.NewUnauthorizedError("invalid credentials")
	}
	token, err := jwt.GenerateToken(u.jwtSecret, user.ID.String(), string(user.Role), 24*time.Hour)
	if err != nil {
		return nil, custom.NewUnexpectedError("failed to generate token")
	}
	return &dto.LoginResponse{Token: token}, nil
}
