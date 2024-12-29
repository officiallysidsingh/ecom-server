package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/officiallysidsingh/ecom-server/internal/models"
	"github.com/officiallysidsingh/ecom-server/internal/store"
)

type contextKey string

const userContextKey contextKey = "user"

func ValidateJWT(db *sqlx.DB, r *http.Request, secretKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the token from Authorization header or cookie
			tokenString, err := extractToken(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Parse and validate the JWT
			user, err := parseJWT(db, r, tokenString, secretKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Optionally, you can set the user in the context for use in subsequent handlers
			ctx := r.Context()
			ctx = context.WithValue(ctx, userContextKey, user)
			r = r.WithContext(ctx)

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

func extractToken(r *http.Request) (string, error) {
	// If not in the Authorization header, check the cookie
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", errors.New("missing Authorization cookie")
		}
		return "", fmt.Errorf("error retrieving cookie: %v", err)
	}

	return cookie.Value, nil
}

func parseJWT(db *sqlx.DB, r *http.Request, tokenString, secretKey string) (*models.User, error) {
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Ensure the token is signed with the correct method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	// If token is invalid
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// If token is valid, extract user claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract info from claims
		userID := claims["user_id"].(string)

		// Initialize userStore
		userStore := store.NewUserStore(db)

		//Fetch user from DB
		user, err := userStore.GetByIdFromDB(r.Context(), userID)
		if err != nil {
			return nil, fmt.Errorf("could not fetch user from DB: %v", err)
		}

		// Return the user
		return user, nil
	}

	return nil, errors.New("invalid token claims")
}

// func ValidateJWT(r *http.Request, secretKey string, db *sqlx.DB) (*models.User, error) {
// 	// Extract JWT token from cookie
// 	cookie, err := r.Cookie("Authorization")
// 	if err != nil {
// 		if err == http.ErrNoCookie {
// 			return nil, errors.New("cookie is missing")
// 		}
// 		return nil, fmt.Errorf("error retrieving cookie: %v", err)
// 	}

// 	// Get the token from cookie
// 	tokenString := cookie.Value

// 	// Parse the JWT token
// 	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
// 		// Validate signing method
// 		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
// 		}
// 		return []byte(secretKey), nil
// 	})

// 	// If token is invalid
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid token: %v", err)
// 	}

// 	// If token is valid
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		// Extract info from claims
// 		userID := claims["user_id"].(string)

// 		// Initialize userStore
// 		userStore := store.NewUserStore(db)

// 		//Fetch user from DB
// 		user, err := userStore.GetByIdFromDB(r.Context(), userID)
// 		if err != nil {
// 			return nil, fmt.Errorf("could not fetch user from DB: %v", err)
// 		}

// 		// Return user if found
// 		return user, nil
// 	}

// 	// Return error if invalid token
// 	return nil, errors.New("invalid token")
// }
