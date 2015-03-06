package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	// General constants
	SESSION_LENGTH = (time.Hour * 24 * 7) // How long sessions last (right now its one week)

	// JWT claim keys for Session
	JWT_CLAIM_USER_ID         = "userId"      // JWT claim key for id of currently authed user
	JWT_CLAIM_USER_FIRST_NAME = "firstName"   // JWT claim key for first name of currently authed user
	JWT_CLAIM_USER_LAST_NAME  = "lastName"    // JWT claim key for last name of currently authed user
	JWT_CLAIM_USER_EMAIL      = "email"       // JWT claim key for email of currently authed user
	JWT_CLAIM_USER_PIC_URL    = "pictureUrl"  // JWT claim key for picture url of currently authed user
	JWT_CLAIM_EXPIRATION      = "expiration"  // JWT claim key for expiration date in secs since epoch
	JWT_CLAIM_TIME_CREATED    = "timeCreated" // JWT claim key for time created in secs since epoch

	// Session-related errors
	ERR_JWT_INVALID_CLAIMS  = "Could not parse JWT token claims" // Error occurs when there was a JWT parsing error
	ERR_JWT_SESSION_EXPIRED = "Session has expired"              // Error occurs when the session has expired

	// Whitelisted API endpoints
	API_ENDPOINT_REGISTER     = "/api/register"
	API_ENDPOINT_AUTHENTICATE = "/api/authenticate"
)

/***************************** TYPE DECLARATIONS ******************************/

// Session is a service that is injected into a martini helper with information
// about the currently authenticated user and the session thereof.
type Session struct {
	// JWT token with all the data inside
	token *jwt.Token
	// Currently authed user's id
	userId uint64
	// Currently authed user's first name
	firstName string
	// Currently authed user's last name
	lastName string
	// Currently authed user's email
	email string
	// Currently authed user's picture url
	pictureUrl string
	// Time when session expires
	expiration int64
	// Time when the session was created
	timeCreated int64
}

// Returns the duration of the current session in seconds
func (s *Session) Duration() int64 {
	return time.Now().Unix() - s.timeCreated
}

/****************************** HELPER FUNCTIONS ******************************/

func lookupJWTKey(token *jwt.Token) (interface{}, error) {
	// Get secret from env
	secret := os.Getenv(ENV_VAR_JWT_SECRET)
	if secret == "" {
		// The secret is bogus
		log.Println("Could not read the JWT secret environment variable")
		return nil, errors.New("JWT secret environment variable was missing")
	} else {
		return secret, nil
	}
}

/***************************** INTERNAL FUNCTIONS *****************************/

// Marshals a session into a JWT token
func MarshalSession(sesh Session, token *jwt.Token) {
	token.Claims[JWT_CLAIM_USER_ID] = strconv.FormatUint(sesh.userId, 64)
	token.Claims[JWT_CLAIM_USER_FIRST_NAME] = sesh.firstName
	token.Claims[JWT_CLAIM_USER_LAST_NAME] = sesh.lastName
	token.Claims[JWT_CLAIM_USER_EMAIL] = sesh.email
	token.Claims[JWT_CLAIM_USER_PIC_URL] = sesh.pictureUrl
	token.Claims[JWT_CLAIM_EXPIRATION] = strconv.FormatInt(sesh.expiration, 64)
	token.Claims[JWT_CLAIM_TIME_CREATED] = strconv.FormatInt(sesh.timeCreated, 64)
}

// Parses a session struct out of a token
func UnmarshalSession(token *jwt.Token) (*Session, error) {
	expirationStr, ok := token.Claims[JWT_CLAIM_EXPIRATION].(string)
	if !ok {
		return nil, errors.New(ERR_JWT_INVALID_CLAIMS)
	}
	// Convert expiration str immediately so we can find out if the session expired
	expiration, err := strconv.ParseInt(expirationStr, 10, 64)
	if err != nil {
		return nil, errors.New(ERR_JWT_INVALID_CLAIMS)
	} else if time.Now().Unix() > expiration {
		return nil, errors.New(ERR_JWT_SESSION_EXPIRED)
	}

	// Next do all the rest of the fields
	userIdStr, ok := token.Claims[JWT_CLAIM_USER_ID].(string)
	if !ok {
		return nil, errors.New(ERR_JWT_INVALID_CLAIMS)
	}
	timeCreatedStr, ok := token.Claims[JWT_CLAIM_TIME_CREATED].(string)
	if !ok {
		return nil, errors.New(ERR_JWT_INVALID_CLAIMS)
	}
	firstName, ok := token.Claims[JWT_CLAIM_USER_FIRST_NAME].(string)
	if !ok {
		return nil, errors.New(ERR_JWT_INVALID_CLAIMS)
	}
	lastName, ok := token.Claims[JWT_CLAIM_USER_LAST_NAME].(string)
	if !ok {
		return nil, errors.New(ERR_JWT_INVALID_CLAIMS)
	}
	email, ok := token.Claims[JWT_CLAIM_USER_EMAIL].(string)
	if !ok {
		return nil, errors.New(ERR_JWT_INVALID_CLAIMS)
	}
	pictureUrl, ok := token.Claims[JWT_CLAIM_USER_PIC_URL].(string)
	if !ok {
		return nil, errors.New(ERR_JWT_INVALID_CLAIMS)
	}

	// Perform checks for both the numbers
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return nil, errors.New(ERR_JWT_INVALID_CLAIMS)
	}
	timeCreated, err := strconv.ParseInt(timeCreatedStr, 10, 64)
	if err != nil {
		return nil, errors.New(ERR_JWT_INVALID_CLAIMS)
	}

	// Everything is hunky dory
	return &Session{
		token:       token,
		userId:      userId,
		firstName:   firstName,
		lastName:    lastName,
		email:       email,
		pictureUrl:  pictureUrl,
		expiration:  expiration,
		timeCreated: timeCreated,
	}, nil
}

/****************************** PUBLIC FUNCTIONS ******************************/

// Martini middleware that provides the session to martini handlers
func Sessionize(res http.ResponseWriter, req *http.Request, c martini.Context) {
	// First check if the current path is whitelisted
	if (req.URL.Path != API_ENDPOINT_REGISTER) && (req.URL.Path != API_ENDPOINT_AUTHENTICATE) {
		token, err := jwt.ParseFromRequest(req, lookupJWTKey)
		// Check out whether we're fine
		if err != nil || !token.Valid {
			http.Error(res, "Authorization token was invalid", http.StatusUnauthorized)
		} else {
			// Embed the token data in the context
			sesh, err := UnmarshalSession(token)
			if err != nil {
				// Could not marshal the session, send back the 401
				http.Error(res, err.Error(), http.StatusUnauthorized)
			} else {
				// Bind the session to the martini context
				c.Map(sesh)
				// Move on
				c.Next()
			}
		}
	} else {
		// This path is whitelisted
		c.Next()
	}
}

// Creates a new token from session data
func NewSessionToken(
	userId uint64,
	firstName string,
	lastName string,
	email string,
) (string, error) {
	// Create the session
	sesh := Session{
		userId:      userId,
		firstName:   firstName,
		lastName:    lastName,
		email:       email,
		expiration:  time.Now().Add(SESSION_LENGTH).Unix(),
		timeCreated: time.Now().Unix(),
	}
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Marshall the session into the token
	MarshalSession(sesh, token)
	// Get the secret key
	key, err := lookupJWTKey(nil)
	if err != nil {
		return "", err
	}
	// Stringify the token
	return token.SignedString(key)
}
