package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	validator "github.com/asaskevich/govalidator"
	"github.com/go-martini/martini"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

const (
	USER_FIELD_FIRST_NAME  = "firstName"
	USER_FIELD_LAST_NAME   = "lastName"
	USER_FIELD_EMAIL       = "email"
	USER_FIELD_PASSWORD    = "password"
	USER_FIELD_PICTURE_URL = "pictureUrl"
)

func SetupUserRoutes(m *martini.ClassicMartini, db *sql.DB, env *Environment) {
	// Registers a new user
	// Expects a JSON encoded body with the following properties:
	// - firstName (string; no longer than 100 characters)
	// - lastName (string; no longer than 100 characters)
	// - email (string; must be email formatted; no longer than 100 characters)
	// - pictureUrl (string; must be URL formatted; no longer than 500 characters)
	m.Post(API_REGISTER_USER, func(req *http.Request, responder *Responder) {
		// Perform json unmarshalling
		var (
			body       map[string]interface{}
			firstName  string
			lastName   string
			email      string
			password   string
			pictureUrl string
			ok         bool
		)

		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&body); err != nil {
			responder.Error(PUBERR_INVALID_JSON)
			return
		}

		// Basic validation and field extractions
		firstName, ok = String(body[USER_FIELD_FIRST_NAME])
		if !ok {
			responder.Error(NewPublicError(http.StatusBadRequest, ERRCODE_INVALID_FIELD, fmt.Sprintf(ERR_BODY_FIELD_INVALID, USER_FIELD_FIRST_NAME)))
			return
		}
		lastName, ok = String(body[USER_FIELD_LAST_NAME])
		if !ok {
			responder.Error(NewPublicError(http.StatusBadRequest, ERRCODE_INVALID_FIELD, fmt.Sprintf(ERR_BODY_FIELD_INVALID, USER_FIELD_LAST_NAME)))
			return
		}
		email, ok = String(body[USER_FIELD_EMAIL])
		if !ok || !validator.IsEmail(email) {
			responder.Error(NewPublicError(http.StatusBadRequest, ERRCODE_INVALID_FIELD, fmt.Sprintf(ERR_BODY_FIELD_INVALID, USER_FIELD_EMAIL)))
			return
		}
		password, ok = String(body[USER_FIELD_PASSWORD])
		if !ok || !IsValidPassword(password) {
			responder.Error(NewPublicError(http.StatusBadRequest, ERRCODE_INVALID_FIELD, fmt.Sprintf(ERR_BODY_FIELD_INVALID, USER_FIELD_PASSWORD)))
			return
		}
		pictureUrl, ok = String(body[USER_FIELD_PICTURE_URL])
		if !ok {
			responder.Error(NewPublicError(http.StatusBadRequest, ERRCODE_INVALID_FIELD, fmt.Sprintf(ERR_BODY_FIELD_INVALID, USER_FIELD_PICTURE_URL)))
			return
		}

		// Build hashed password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 7)
		if err != nil {
			responder.Error(err)
			return
		}

		// Get the stripe id
		// TODO get the stripe API working
		stripeId := "FAKE_STRIPE_ID"

		// Put the user in the database
		newId, err := CreateNewUser(db, firstName, lastName, email, string(hashedPassword[:]), stripeId, pictureUrl)
		if err != nil {
			responder.Error(err)
			return
		} else {
			newUser, err := GetUser(db, newId)
			if err != nil {
				responder.Error(err)
				return
			} else {
				responder.Json(newUser)
				return
			}
		}
	})
}
