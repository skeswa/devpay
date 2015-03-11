package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
	"strings"
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
)

/***************************** TYPE DECLARATIONS ******************************/

// Session is a service that is injected into a martini helper with information
// about the currently authenticated user and the session thereof.
type Session struct {
	// JWT token with all the data inside
	token *jwt.Token
	// Currently authed user's id
	UserId int64 `json:"userId"`
	// Currently authed user's first name
	FirstName string `json:"firstName"`
	// Currently authed user's last name
	LastName string `json:"lastName"`
	// Currently authed user's email
	Email string `json:"email"`
	// Currently authed user's picture url
	PictureUrl string `json:"pictureUrl"`
	// Time when session expires
	Expiration int64 `json:"expiration"`
	// Time when the session was created
	TimeCreated int64 `json:"-"`
}

// Returns the duration of the current session in seconds
func (s *Session) Duration() int64 {
	return time.Now().Unix() - s.TimeCreated
}

/***************************** INTERNAL FUNCTIONS *****************************/

// Marshals a session into a JWT token
func MarshalSession(sesh Session, token *jwt.Token) {
	token.Claims[JWT_CLAIM_USER_ID] = strconv.FormatInt(sesh.UserId, 10)
	token.Claims[JWT_CLAIM_USER_FIRST_NAME] = sesh.FirstName
	token.Claims[JWT_CLAIM_USER_LAST_NAME] = sesh.LastName
	token.Claims[JWT_CLAIM_USER_EMAIL] = sesh.Email
	token.Claims[JWT_CLAIM_USER_PIC_URL] = sesh.PictureUrl
	token.Claims[JWT_CLAIM_EXPIRATION] = strconv.FormatInt(sesh.Expiration, 10)
	token.Claims[JWT_CLAIM_TIME_CREATED] = strconv.FormatInt(sesh.TimeCreated, 10)
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
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
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
		UserId:      userId,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		PictureUrl:  pictureUrl,
		Expiration:  expiration,
		TimeCreated: timeCreated,
	}, nil
}

/****************************** PUBLIC FUNCTIONS ******************************/

// Creates a new token from session data
func NewSessionToken(
	env *Environment,
	userId int64,
	firstName string,
	lastName string,
	email string,
	pictureUrl string,
) (string, error) {
	// Create the session
	sesh := Session{
		UserId:      userId,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		PictureUrl:  pictureUrl,
		Expiration:  time.Now().Add(SESSION_LENGTH).Unix(),
		TimeCreated: time.Now().Unix(),
	}
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Marshall the session into the token
	MarshalSession(sesh, token)
	// Stringify the token
	return token.SignedString([]byte(env.jwtSecret))
}

// Martini middleware that provides the session to martini handlers
func Sessionize(res http.ResponseWriter, req *http.Request, env *Environment, c martini.Context) {
	// First check if the current path is not blacklisted
	if strings.Index(req.URL.Path, API_PREFIX) == 0 &&
		(req.URL.Path != API_AUTHENTICATE) &&
		(req.URL.Path != API_REGISTER_USER) {
		// Get the JWT token
		token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
			return []byte(env.jwtSecret), nil
		})
		// Check out whether the token is good
		if err != nil || !token.Valid {
			res.Header().Set(ContentType, ContentJSON)
			res.WriteHeader(http.StatusUnauthorized)
			res.Write(PUBERR_INVALID_AUTH_TOKEN.Json)
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
