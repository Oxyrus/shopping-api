package controllers

import (
	"github.com/Oxyrus/shopping/internal/models"
	"github.com/Oxyrus/shopping/internal/utils"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserController struct {
	TokenAuth *jwtauth.JWTAuth
	DB        *sqlx.DB
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	type LoginCredentials struct {
		Username string `json:"username"`
	}

	credentials := &LoginCredentials{}
	if err := render.DecodeJSON(r.Body, credentials); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := c.generateToken(credentials.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	type RegisterPayload struct {
		FirstName            string `json:"firstName"`
		LastName             string `json:"lastName"`
		Password             string `json:"password"`
		PasswordConfirmation string `json:"passwordConfirmation"`
		Email                string `json:"email"`
	}

	newUser := &RegisterPayload{}
	if err := render.DecodeJSON(r.Body, newUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if newUser.Password != newUser.PasswordConfirmation {
		utils.WriteError(w, http.StatusBadRequest, "Passwords must match")
		return
	}

	// Validate if a user with the same email address already exists
	existingUser := models.User{}
	err := c.DB.Get(&existingUser, "SELECT * FROM users LIMIT 1")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Error connecting verifying if an user with the same email already exists")
		return
	}

	password, err := generatePasswordHash(newUser.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	user := models.User{
		FirstName:    newUser.FirstName,
		LastName:     newUser.LastName,
		PasswordHash: password,
		Email:        newUser.Email,
	}

	tx := c.DB.MustBegin()
	tx.MustExec("INSERT INTO users (first_name, last_name, password_hash, email) VALUES ($1, $2, $3, $4)", user.FirstName, user.LastName, user.PasswordHash, user.Email)
	err = tx.Commit()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Error saving user in the database")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *UserController) Profile(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Profile"))
}

func generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", nil
	}

	return string(bytes), nil
}

func (c *UserController) generateToken(username string) (string, error) {
	_, tokenString, err := c.TokenAuth.Encode(map[string]interface{}{"username": username})
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
