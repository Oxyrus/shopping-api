package models

import (
	"database/sql"
	"time"
)

// User represents a user of the application and is the template
// for the objects stored in the database.
type User struct {
	ID           int          `json:"id" db:"id"`
	FirstName    string       `json:"firstName" db:"first_name"`
	LastName     string       `json:"lastName" db:"last_name"`
	PasswordHash string       `json:"password" db:"password_hash"`
	Email        string       `json:"email" db:"email"`
	CreatedOn    time.Time    `json:"createdOn" db:"created_on"`
	LastLogin    sql.NullTime `json:"LastLogin" db:"last_login"`
}

// LoginCredentials represents the payload used during a login attempt.
type LoginCredentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterPayload represents the payload used during user registration.
type RegisterPayload struct {
	FirstName            string `json:"firstName" binding:"required"`
	LastName             string `json:"lastName" binding:"required"`
	Password             string `json:"password" binding:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required"`
	Email                string `json:"email" binding:"required"`
}
