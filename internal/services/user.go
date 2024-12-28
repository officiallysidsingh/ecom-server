package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/officiallysidsingh/ecom-server/internal/config"
	"github.com/officiallysidsingh/ecom-server/internal/models"
	"github.com/officiallysidsingh/ecom-server/internal/store"
	"github.com/officiallysidsingh/ecom-server/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailExists     = errors.New("email already registered")
	ErrHashingPassword = errors.New("error hashing password")
	ErrCreatingUser    = errors.New("error creating user")
)

type UserService interface {
	Login(ctx context.Context, loginReq *models.LoginRequest) (*models.User, string, error)
	Signup(ctx context.Context, user *models.SignupRequest) (string, string, error)
}

type userService struct {
	store     store.UserStore
	jwtSecret string
}

func NewUserService(store store.UserStore, envConfig *config.EnvConfig) UserService {
	return &userService{
		store:     store,
		jwtSecret: envConfig.JWT_SECRET,
	}
}

func (s *userService) Login(ctx context.Context, loginReq *models.LoginRequest) (*models.User, string, error) {
	// Fetch user from DB by Email
	user, err := s.store.GetByEmailFromDB(ctx, loginReq.Email)
	if err != nil {
		return nil, "", err
	}

	// Verify the password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		return nil, "", err
	}

	// Create JWT Config
	tokenConfig := config.NewJWTConfig(s.jwtSecret)

	// Generate JWT token
	token, err := utils.GenerateJWT(user.UserID, user.Role, tokenConfig)
	if err != nil {
		return nil, "", err
	}

	// Return user data
	user.Password = ""
	return user, token, nil
}

func (s *userService) Signup(ctx context.Context, user *models.SignupRequest) (string, string, error) {
	// Check if email already exists
	existingUser, err := s.store.GetByEmailFromDB(ctx, user.Email)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return "", "", fmt.Errorf("error checking email: %w", err)
	}

	if existingUser != nil {
		return "", "", ErrEmailExists
	}

	// If email doesn't exist
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", ErrHashingPassword, err)
	}

	// Prepare user for creation
	user.Password = string(hashedPassword)
	user.Role = "user"

	// Create user in DB
	userID, err := s.store.CreateInDB(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", ErrCreatingUser, err)
	}

	// Create JWT Config
	tokenConfig := config.NewJWTConfig(s.jwtSecret)

	// Generate JWT token
	token, err := utils.GenerateJWT(userID, user.Role, tokenConfig)
	if err != nil {
		return "", "", fmt.Errorf("user created but error generating token: %w", err)
	}

	return userID, token, nil
}
