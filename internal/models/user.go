package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int          `json:"id" db:"id"`
	FirstName    string       `json:"firstName" db:"first_name"`
	LastName     string       `json:"lastName" db:"last_name"`
	PasswordHash string       `json:"password" db:"password_hash"`
	Email        string       `json:"email" db:"email"`
	CreatedOn    time.Time    `json:"createdOn" db:"created_on"`
	LastLogin    sql.NullTime `json:"LastLogin" db:"last_login"`
}
