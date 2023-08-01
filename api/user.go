package main

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	DOB       string `json:"dob"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := insertUserIntoDB(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	user.ID = id
	c.JSON(http.StatusCreated, user)
}

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func generateRandString(size int) string {

	b := make([]rune, size)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func extractRole(c *gin.Context) (string, error) {
	claims, ok := c.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	if !ok {
		return "", errors.New("JWT Information Missing")
	}

	// claims.RegisteredClaims.Audience
	//Here is how we pull the custom claims from the JWT token
	cciClaims, ok := claims.CustomClaims.(*CCIClaim)
	if !ok {
		return "", errors.New("CCI Custom Claims Missing")
	}
	return cciClaims.CCI.Account.Role, nil
}

func getUser(c *gin.Context) {

	id := c.Param("id")
	role, err := extractRole(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	if role == "guest" {
		s := c.Query("size")
		sz, _ := strconv.Atoi(s)

		c.JSON(http.StatusOK, gin.H{"message": generateRandString(sz)})
		return
	}
	user, err := getUserFromDB(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func updateUser(c *gin.Context) {
	id := c.Param("id")

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := updateUserInDB(id, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := deleteUserFromDB(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Database operations

func insertUserIntoDB(user *User) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO users (firstname, lastname, email, username, password, dob, phone, address) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		user.FirstName, user.LastName, user.Email, user.Username, user.Password, user.DOB, user.Phone, user.Address).Scan(&id)
	if err != nil {
		log.Println("insertUserIntoDB", err)
		return 0, err
	}

	return id, nil
}

func getUserFromDB(userID string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, firstname, lastname, email, username, dob, phone, address FROM users WHERE id = $1", userID).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.DOB, &user.Phone, &user.Address,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func updateUserInDB(userID string, user *User) error {
	_, err := db.Exec("UPDATE users SET firstname = $1, lastname = $2, email = $3, username = $4, dob = $5, phone = $6, address = $7 WHERE id = $8",
		user.FirstName, user.LastName, user.Email, user.Username, user.DOB, user.Phone, user.Address, userID)
	if err != nil {
		return err
	}

	return nil
}

func deleteUserFromDB(userID string) error {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return err
	}

	return nil
}