package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

// HERE ARE SOME CUSTOM CLAIMS WE CAN USE - THEY Look Like This
/*
"cci": {
    "client_id": "c1",
    "account": {
      "id": 1234,
      "group": "*",
      "role": "customer-premium"
    }
  }

  --- OR

  "cci": {
    "subject": "s1",
    "account": {
      "id": 5678,
      "group": "*",
      "role": "user-basic"
    }
*/
type CCIClaim struct {
	CCI CCICustomClaims `json:"cci"`
}

// Note the CCI claims expects either a client_id or a subject
// we have a validator below to handle this
type CCICustomClaims struct {
	ClientID string       `json:"client_id,omitempty"`
	Subject  string       `json:"subject,omitempty"`
	Account  AccountClaim `json:"account"`
}

type AccountClaim struct {
	ID    int    `json:"id"`
	Group string `json:"group"`
	Role  string `json:"role"`
}

// Validate errors out if `ShouldReject` is true.
func (c *CCIClaim) Validate(ctx context.Context) error {
	if c.CCI.Subject == "" && c.CCI.ClientID == "" {
		return errors.New("cci token should include either a subject or a client_id")
	}
	return nil
}

// This package implements JWT middleware for the gin framework.  It is
// somewhat basic but is useful for playing with JWTs and understanding
// how they work.  It is not suitable for production use.
type JwtMiddleware struct {
	mw     *jwtmiddleware.JWTMiddleware
	mwFunc gin.HandlerFunc
}

func NewJwtMiddleware(issuer string, audience []string) *JwtMiddleware {
	mw := setupHandlerMiddleware(issuer, audience)
	mwFunc := checkJWT(mw)
	return &JwtMiddleware{
		mw:     mw,
		mwFunc: mwFunc,
	}
}

func (j *JwtMiddleware) CheckJWT() gin.HandlerFunc {
	return j.mwFunc
}

func setupHandlerMiddleware(issuer string, audience []string) *jwtmiddleware.JWTMiddleware {
	issuerURL, err := url.Parse(issuer)
	if err != nil {
		log.Fatalf("failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	customClaims := func() validator.CustomClaims {
		return &CCIClaim{}
	}

	// Set up the validator with support for our custom claims
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		audience,
		validator.WithCustomClaims(customClaims),
		validator.WithAllowedClockSkew(30*time.Second),
	)
	if err != nil {
		log.Fatalf("failed to set up the validator: %v", err)
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Encountered error while validating JWT: %v", err)
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return middleware
}

func checkJWT(m *jwtmiddleware.JWTMiddleware) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		encounteredError := true
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			ctx.Request = r
			ctx.Next()
		}

		m.CheckJWT(handler).ServeHTTP(ctx.Writer, ctx.Request)

		if encounteredError {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				map[string]string{"message": "JWT is invalid."},
			)
		}
	}
}
