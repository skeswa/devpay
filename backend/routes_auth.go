package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func SetupAuthRoutes(m *martini.ClassicMartini, db *sql.DB, env *Environment) {
	// Log's a user in; creates a session token
	m.Post(API_AUTHENTICATE, func(req *http.Request, responder *Responder) {
		// Perform json unmarshalling
		var (
			body     map[string]interface{}
			email    string
			password string
			ok       bool
		)

		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&body); err != nil {
			responder.Error(PUBERR_INVALID_JSON)
			return
		}

		// Basic validation and field extractions
		email, ok = String(body[USER_FIELD_EMAIL])
		if !ok {
			responder.Error(NewPublicError(http.StatusBadRequest, ERRCODE_INVALID_FIELD, fmt.Sprintf(ERR_BODY_FIELD_INVALID, USER_FIELD_EMAIL)))
			return
		}
		password, ok = String(body[USER_FIELD_PASSWORD])
		if !ok {
			responder.Error(NewPublicError(http.StatusBadRequest, ERRCODE_INVALID_FIELD, fmt.Sprintf(ERR_BODY_FIELD_INVALID, USER_FIELD_PASSWORD)))
			return
		}

		// Find the user
		user, err := FindUserByEmail(db, email)
		if err != nil {
			responder.Error(PUBERR_INVALID_CREDENTIALS)
		} else {
			err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
			if err != nil {
				responder.Error(PUBERR_INVALID_CREDENTIALS)
			} else {
				token, err := NewSessionToken(
					user.Id,
					user.FirstName,
					user.LastName,
					user.Email,
					user.PictureUrl,
				)
				if err != nil {
					responder.Error(err)
				} else {
					responder.Text(token)
				}
			}
		}
	})
}
