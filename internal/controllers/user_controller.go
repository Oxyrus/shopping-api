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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	credentials := &LoginCredentials{}
	if err := render.DecodeJSON(r.Body, credentials); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	user := models.User{}
	err := c.DB.Get(&user, "SELECT * FROM users WHERE email=$1", credentials.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Could not retrieve the user from the database")
		return
	}

	if user.ID == 0 {
		utils.WriteError(w, http.StatusInternalServerError, "Wrong credentials")
		return
	}

	if passwordsAreEqual := arePasswordsEqual([]byte(user.PasswordHash), []byte(credentials.Password)); !passwordsAreEqual {
		utils.WriteError(w, http.StatusBadRequest, "Wrong credentials")
		return
	}

	claims := make(map[string]interface{})
	claims["firstName"] = user.FirstName
	claims["email"] = user.Email
	token, err := c.generateToken(claims)
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
		utils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if newUser.Password != newUser.PasswordConfirmation {
		utils.WriteError(w, http.StatusBadRequest, "Passwords must match")
		return
	}

	// Validate if a user with the same email address already exists
	existingUser := models.User{}
	err := c.DB.Get(&existingUser, "SELECT * FROM users WHERE email=$1", newUser.Email)
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

func arePasswordsEqual(hashedPassword, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		return false
	}
	return true
}

func (c *UserController) generateToken(claims map[string]interface{}) (string, error) {
	_, tokenString, err := c.TokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
