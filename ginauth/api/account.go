package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Account struct {
	ID            int     `json:"id"`
	UserID        int     `json:"user_id"`
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
}

// var (
// 	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
// )

// func generateRandString(size int) string {

// 	b := make([]rune, size)
// 	for i := range b {
// 		b[i] = letterRunes[rand.Intn(len(letterRunes))]
// 	}
// 	return string(b)
// }

// func extractRole(c *gin.Context) (string, error) {
// 	claims, ok := c.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
// 	if !ok {
// 		return "", errors.New("JWT Information Missing")
// 	}

// 	// claims.RegisteredClaims.Audience
// 	//Here is how we pull the custom claims from the JWT token
// 	cciClaims, ok := claims.CustomClaims.(*CCIClaim)
// 	if !ok {
// 		return "", errors.New("CCI Custom Claims Missing")
// 	}
// 	return cciClaims.CCI.Account.Role, nil
// }

func getAccount(c *gin.Context) {
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

	account, err := getAccountFromDB(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func payBill(c *gin.Context) {
	id := c.Param("id")

	err := markBillPaid(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark bill as paid"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bill payment successful"})
}

func getDueDate(c *gin.Context) {
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

	dueDate, err := getPaymentDueDate(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"due_date": dueDate})
}

func getCreditScore(c *gin.Context) {
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

	creditScore, err := getCreditScoreFromDB(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Credit score not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"credit_score": creditScore})
}

// Testing
func getPublicCreditScore(c *gin.Context) {
	id := c.Param("id")

	creditScore, err := getCreditScoreFromDB(id)
	if err != nil {
		log.Println("Get Credit Score :", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Credit score not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"credit_score": creditScore})
}

// Database operations

func getAccountFromDB(accountID string) (*Account, error) {
	var account Account
	err := db.QueryRow("SELECT id, user_id, account_number, balance FROM accounts WHERE id = $1", accountID).Scan(
		&account.ID, &account.UserID, &account.AccountNumber, &account.Balance,
	)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func markBillPaid(accountID string) error {
	_, err := db.Exec("UPDATE accounts SET paid = true WHERE id = $1", accountID)
	if err != nil {
		return err
	}

	return nil
}

func getPaymentDueDate(accountID string) (string, error) {
	var dueDate string
	err := db.QueryRow("SELECT due_date FROM payments WHERE account_id = $1", accountID).Scan(&dueDate)
	if err != nil {
		return "", err
	}

	return dueDate, nil
}

func getCreditScoreFromDB(userID string) (float64, error) {
	var creditScore float64
	// result := db.QueryRow("SELECT COUNT(*) FROM credit_scores")
	// var count int
	// err1 := result.Scan(&count)
	err := db.QueryRow("SELECT score FROM credit_scores WHERE id = $1", userID).Scan(&creditScore)
	if err != nil {
		return 0, err
	}

	// if err1 != nil {
	// 	return 0, err
	// }

	return creditScore, nil
}
