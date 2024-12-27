package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/officiallysidsingh/ecom-server/internal/models"
	"github.com/officiallysidsingh/ecom-server/internal/utils"
)

type UserStore interface {
	GetByEmailFromDB(ctx context.Context, email string) (*models.User, error)
	GetByIdFromDB(ctx context.Context, userID string) (*models.User, error)
	CreateInDB(ctx context.Context, user *models.User) (string, error)
}

type userStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) UserStore {
	return &userStore{
		db: db,
	}
}

func (s *userStore) GetByEmailFromDB(ctx context.Context, userEmail string) (*models.User, error) {
	var user models.User

	// SQL query to get user by email
	query := `
		SELECT user_id, name, email, password, role
		FROM users
		WHERE email = $1
	`

	fields := []interface{}{
		userEmail,
	}

	if err := utils.ExecGetQuery(
		s.db,
		query,
		fields,
		&user,
	); err != nil {
		// If no rows found
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("User with Email %s not found", userEmail)
			return nil, fmt.Errorf("user with Email %s not found", userEmail)
		}
		log.Printf("Error fetching user with Email %s from DB: %v", userEmail, err)
		return nil, err
	}

	return &user, nil
}

func (s *userStore) GetByIdFromDB(ctx context.Context, userID string) (*models.User, error) {
	var user models.User

	// SQL query to get user by email
	query := `
		SELECT user_id, name, email, role
		FROM users
		WHERE user_id = $1
	`

	fields := []interface{}{
		userID,
	}

	if err := utils.ExecGetQuery(
		s.db,
		query,
		fields,
		&user,
	); err != nil {
		// If no rows found
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("User with ID %s not found", userID)
			return nil, fmt.Errorf("user with ID %s not found", userID)
		}
		log.Printf("Error fetching user with Email %s from DB: %v", userID, err)
		return nil, err
	}

	return &user, nil
}

func (s *userStore) CreateInDB(ctx context.Context, user *models.User) (string, error) {
	// Begin a transaction to ensure atomicity
	tx, err := s.db.Beginx()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return "", fmt.Errorf("failed to start database transaction: %w", err)
	}

	// Ensure transaction is properly rolled back in case of failure
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Error rolling back transaction: %v", rollbackErr)
			}
		}
	}()

	// SQL query to insert a new user
	query := `
		INSERT INTO users (user_id, name, email, password, role, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING email
	`

	fields := []interface{}{
		user.Name,
		user.Email,
		user.Password,
		user.Role,
	}

	// Execute the query and return the added email
	var userEmail string
	if err = utils.ExecGetTransactionQuery(
		s.db,
		tx,
		query,
		fields,
		&userEmail,
	); err != nil {
		log.Printf("Error adding user with Email %s to DB: %v", user.Email, err)
		return "", err
	}

	// Commit the transaction if update was successful
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction for user with Email %s: %v", user.Email, err)
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Log the success and return the added email
	log.Printf("User with Email %s added successfully", user.Email)
	return userEmail, nil
}
