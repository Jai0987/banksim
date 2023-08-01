package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	IdpDefaultLocation = "http://host.docker.internal:1234"
)

var (
	oauthConfig *oauth2.Config
	db          *sql.DB
)

func OAuthTest(c *gin.Context) {
	//Demo showing pulling the claims from the JWT token
	//Note we will not even get here if the JWT token is not valid
	claims, ok := c.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			map[string]string{"message": "Failed to get validated JWT claims."},
		)
		return
	}

	// claims.RegisteredClaims.Audience
	//Here is how we pull the custom claims from the JWT token
	cciClaims, ok := claims.CustomClaims.(*CCIClaim)
	if !ok {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			map[string]string{"message": "CCI Custom claims missing."},
		)
		return
	}
	switch cciClaims.CCI.Account.Role {
	case "customer-premium":
		c.JSON(http.StatusOK, gin.H{"message": "Role = Customer-premium"})
	case "customer-basic":
		c.JSON(http.StatusOK, gin.H{"message": "Role = Customer-basic"})
	case "guest":
		c.JSON(http.StatusOK, gin.H{"message": "Role = Guest"})
	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "This is a page feature"})
	}
	// c.JSON(http.StatusOK, cciClaims)
}

// func OAuthTest(c *gin.Context) {
// 	// Claims should be available in the context after the JwtMiddleware.
// 	claims, exists := c.Get("jwt-claims")
// 	if !exists {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get JWT claims."})
// 		return
// 	}

// 	// You can access the custom claims from the JWT token using type assertion.
// 	customClaims, ok := claims.(*CCIClaim)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get custom JWT claims."})
// 		return
// 	}

// 	// Now you can access any custom fields you added to the JWT token.
// 	// For example: userID := customClaims.UserID

// 	c.JSON(http.StatusOK, gin.H{"message": "OAuth test successful!", "claims": customClaims})
// }

func main() {
	var err error

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Use localhost as the default database URL if DATABASE_URL environment variable is not set
		databaseURL = "postgres://jaikash12:jaikash12@localhost/ginauth?sslmode=disable"
	}

	// Database connection
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	defer db.Close()

	r := gin.Default()

	// User details API endpoints
	r.POST("/users", createUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)

	// Account operations API endpoints

	r.POST("/accounts/:id/pay", payBill)

	IdpUrl := os.Getenv("IDP_URL")
	//This handles the default condition
	if IdpUrl == "" {
		IdpUrl = IdpDefaultLocation
	}
	// OAuth2 endpoints
	//mw := api.SetupHandlerMiddleware("http://localhost:9999", []string{"https://api.cci.drexel.edu"})
	log.Println("IDP_URL ", IdpUrl)
	mw := NewJwtMiddleware(IdpUrl,
		[]string{"https://api.cci.drexel.edu", "cci.drexel.edu/api"})

	r.GET("/oauth", mw.CheckJWT(), OAuthTest)
	r.GET("/users/:id", mw.CheckJWT(), getUser)
	r.GET("/accounts/:id", mw.CheckJWT(), getAccount)
	r.GET("/accounts/:id/due", getDueDate)
	r.GET("/accounts/:id/score", mw.CheckJWT(), getCreditScore)
	r.GET("/uaccounts/:id/score", getPublicCreditScore)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}
	fmt.Println("Server is running on port:", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
