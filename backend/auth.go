package main

import (
	"errors"
	"net/http"
	"os"
)

const (
	ENV_VAR_JWT_SECRET = "JWT_SECRET"
)

func lookupJWTKey(token *jwt.Token) (interface{}, error) {
	// Get secret from env
	secret := os.Getenv(ENV_VAR_JWT_SECRET)
	if secret == "" {
		// The secret is bogus
		return nil, errors.New("JWT secret environment variable was missing")
	} else {
		return secret, nil
	}
}

func JWTMiddleware(res http.ResponseWriter, req *http.Request) {
	token, err := jwt.ParseFromRequest(req, lookupJWTKey)
	// Check out whether we're fine
	if err != nil {
		http.Error(res, "Authorization token was invalid", http.StatusUnauthorized)
	} else {
		// Embed the token data in the context
	}
}
