package controllers

import (
	"github.com/Oxyrus/shopping/internal/models"
	"github.com/Oxyrus/shopping/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// UserController provides handlers to deal with user registration and login.
type UserController struct {
	DB *sqlx.DB
}

// Login authenticates a user against the database and returns
// a signed JWT in case the login was successful.
func (c *UserController) Login(ctx *gin.Context) {
	credentials := &models.LoginCredentials{}
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{}
	err := c.DB.Get(&user, "SELECT * FROM users WHERE email=$1", credentials.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not retrieve the user from the database"})
		return
	}

	if user.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Wrong credentials"})
		return
	}

	// ToDo compare passwords
	if passwordsAreEqual := arePasswordsEqual([]byte(user.PasswordHash), []byte(credentials.Password)); !passwordsAreEqual {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Wrong credentials"})
		return
	}

	token, err := utils.GenerateToken(user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating JWT"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// Register attempts to create a new user in the database. Internally it validates against
// all the fields that are required, matches the passwords and stores a hashed version in the
// database.
func (c *UserController) Register(ctx *gin.Context) {
	newUser := &models.RegisterPayload{}
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newUser.Password != newUser.PasswordConfirmation {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Passwords must match"})
		return
	}

	// Validate if a user with the same email address already exists
	existingUser := models.User{}
	err := c.DB.Get(&existingUser, "SELECT * FROM users WHERE email=$1", newUser.Email)
	if err == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An user with that email address already exists"})
		return
	}

	password, err := generatePasswordHash(newUser.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user in the database"})
		return
	}

	ctx.Status(http.StatusCreated)
}

// generatePasswordHash takes a password in its string representation
// and returns a hashed version.
func generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", nil
	}

	return string(bytes), nil
}

// arePasswordsEqual validates a password against its hashed counterpart
// and returns whether the passwords are equivalent.
func arePasswordsEqual(hashedPassword, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		return false
	}
	return true
}
