package main

import (
	"net/http"
	"strings"
)

const (
	ERRCODE_EMAIL_TAKEN         = "EMAIL_TAKEN"
	ERRCODE_INTERNAL_ERROR      = "INTERNAL_ERROR"
	ERRCODE_INVALID_JSON        = "INVALID_JSON"
	ERRCODE_INVALID_FIELD       = "INVALID_FIELD"
	ERRCODE_INVALID_AUTH_TOKEN  = "INVALID_AUTH_TOKEN"
	ERRCODE_INVALID_CREDENTIALS = "INVALID_CREDENTIALS"
	ERRCODE_ENTITY_NOT_FOUND    = "ENTITY_NOT_FOUND"

	ERR_INTERNAL_SERVER_ERROR = "There was an internal issue"
	ERR_INVALID_AUTH_TOKEN    = "Authorization token is invalid"
	ERR_USER_CREATION_FAILED  = "Could not create new user: "
	ERR_INVALID_CREDENTIALS   = "Given credentials were invalid"
	ERR_TABLE_CREATION_FAILED = "Failed to create database table \"%s\": %s"
	ERR_ENV_VAR_MISSING       = "The environment variable \"%s\" was either missing or invalid"
	ERR_COULDNT_START         = "Couldn't start the the server: "
	ERR_JWT_INVALID_CLAIMS    = "Could not parse JWT token claims" // Error occurs when there was a JWT parsing error
	ERR_JWT_SESSION_EXPIRED   = "Session has expired"              // Error occurs when the session has expired
	ERR_BODY_INVALID_JSON     = "Body was invalid JSON"
	ERR_BODY_FIELD_INVALID    = "The \"%s\" field is invalid or ill-formatted"
	ERR_COULD_NOT_HASH_PASS   = "Failed to hash the password field"
	ERR_COULD_CREATE_USER     = "Failed to create a new user"
	ERR_ENTITY_NOT_FOUND      = "Could not find entity matching provided information"
)

var (
	PUBERR_INTERNAL_SERVER_ERROR            = NewPublicError(http.StatusInternalServerError, ERRCODE_INTERNAL_ERROR, ERR_INTERNAL_SERVER_ERROR)
	PUBERR_INVALID_AUTH_TOKEN               = NewPublicError(http.StatusUnauthorized, ERRCODE_INVALID_AUTH_TOKEN, ERR_INVALID_AUTH_TOKEN)
	PUBERR_USER_CREATION_FAILED_EMAIL_TAKEN = NewPublicError(http.StatusInternalServerError, ERRCODE_EMAIL_TAKEN, ERR_USER_CREATION_FAILED+" email address is already in user")
	PUBERR_INVALID_JSON                     = NewPublicError(http.StatusBadRequest, ERRCODE_INVALID_JSON, ERR_BODY_INVALID_JSON)
	PUBERR_INVALID_CREDENTIALS              = NewPublicError(http.StatusUnauthorized, ERRCODE_INVALID_CREDENTIALS, ERR_INVALID_CREDENTIALS)
	PUBERR_ENTITY_NOT_FOUND                 = NewPublicError(http.StatusNotFound, ERRCODE_ENTITY_NOT_FOUND, ERR_ENTITY_NOT_FOUND)
)

type PublicError struct {
	Code    string
	Message string
	Json    []byte
	Status  int
}

func (err *PublicError) Error() string {
	return err.Message
}

func escapeStringForJson(str string) string {
	return strings.Replace(str, "\"", "\\\"", -1)
}

func NewPublicError(status int, code string, message string) *PublicError {
	pubErr := PublicError{
		Code:    code,
		Message: message,
		Status:  status,
		Json:    []byte("{\"code\":\"" + code + "\",\"message\":\"" + escapeStringForJson(message) + "\"}"),
	}
	return &pubErr
}

func IsPublicError(err error) bool {
	if err == nil {
		return false
	} else {
		_, ok := err.(*PublicError)
		return ok
	}
}
