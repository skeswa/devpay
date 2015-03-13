package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
	"time"
)

const (
	CAMPAIGN_FIELD_ID                    = "id"
	CAMPAIGN_FIELD_TITLE                 = "title"
	CAMPAIGN_FIELD_DESCRIPTION           = "description"
	CAMPAIGN_FIELD_COVER_PICTURE_URL     = "coverPictureUrl"
	CAMPAIGN_FIELD_THUMBNAIL_PICTURE_URL = "thumbnailPictureUrl"
	CAMPAIGN_FIELD_AMOUNT                = "amount"
	CAMPAIGN_FIELD_DEADLINE              = "deadline"
)

func SetupCampaignRoutes(m *martini.ClassicMartini, db *sql.DB, env *Environment) {
	// Registers a new campaign
	// Expects a JSON encoded body with the following properties:
	// - firstName (string; no longer than 100 characters)
	// - lastName (string; no longer than 100 characters)
	// - email (string; must be email formatted; no longer than 100 characters)
	// - pictureUrl (string; must be URL formatted; no longer than 500 characters)
	m.Post(API_CREATE_CAMPAIGN, func(req *http.Request, session *Session, responder *Responder) {
		// Perform json unmarshalling
		var (
			body                map[string]interface{}
			title               string
			description         string
			coverPictureUrl     string
			thumbnailPictureUrl string
			deadlineStr         string
			deadlineLong        int64
			deadline            time.Time
			ok                  bool
			err                 error
		)

		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&body); err != nil {
			responder.Error(PUBERR_INVALID_JSON)
			return
		}

		// Basic validation and field extractions
		title, ok = String(body[CAMPAIGN_FIELD_TITLE])
		if !ok {
			responder.Error(NewPublicError(http.StatusBadRequest, ERRCODE_INVALID_FIELD, fmt.Sprintf(ERR_BODY_FIELD_INVALID, CAMPAIGN_FIELD_TITLE)))
			return
		}
		description, ok = String(body[CAMPAIGN_FIELD_DESCRIPTION])
		if !ok {
			responder.Error(NewPublicError(http.StatusBadRequest, ERRCODE_INVALID_FIELD, fmt.Sprintf(ERR_BODY_FIELD_INVALID, CAMPAIGN_FIELD_DESCRIPTION)))
			return
		}
		coverPictureUrl, ok = String(body[CAMPAIGN_FIELD_COVER_PICTURE_URL])
		if !ok {
			responder.Error(NewPublicError(http.StatusBadRequest, ERRCODE_INVALID_FIELD, fmt.Sprintf(ERR_BODY_FIELD_INVALID, CAMPAIGN_FIELD_COVER_PICTURE_URL)))
			return
		}
		thumbnailPictureUrl, ok = String(body[CAMPAIGN_FIELD_THUMBNAIL_PICTURE_URL])
		if !ok {
			responder.Error(NewPublicError(http.StatusBadRequest, ERRCODE_INVALID_FIELD, fmt.Sprintf(ERR_BODY_FIELD_INVALID, CAMPAIGN_FIELD_THUMBNAIL_PICTURE_URL)))
			return
		}
		deadlineStr, ok = String(body[CAMPAIGN_FIELD_DEADLINE])
		if !ok {
			responder.Error(NewPublicError(http.StatusBadRequest, ERRCODE_INVALID_FIELD, fmt.Sprintf(ERR_BODY_FIELD_INVALID, CAMPAIGN_FIELD_DEADLINE)))
			return
		}
		deadlineLong, err = strconv.ParseInt(deadlineStr, 10, 64)
		if err != nil {
			responder.Error(NewPublicError(http.StatusBadRequest, ERRCODE_INVALID_FIELD, fmt.Sprintf(ERR_BODY_FIELD_INVALID, CAMPAIGN_FIELD_DEADLINE)))
			return
		}
		deadline = time.Unix(deadlineLong, 0)

		// Put the campaign in the database
		newId, err := CreateNewCampaign(db, title, description, coverPictureUrl, thumbnailPictureUrl, 0, deadline, session.UserId)
		if err != nil {
			responder.Error(err)
			return
		} else {
			// Return the new campaign
			newCampaign, err := GetCampaign(db, newId)
			if err != nil {
				responder.Error(err)
				return
			} else {
				responder.Json(newCampaign)
				return
			}
		}
	})

	// Gets a info about a specific campaign
	m.Get(API_GET_CAMPAIGN, func(params martini.Params, responder *Responder) {
		id, err := strconv.ParseInt(params[CAMPAIGN_FIELD_ID], 10, 64)
		if err != nil {
			responder.Error(NewPublicError(http.StatusBadRequest, ERRCODE_INVALID_PARAM, fmt.Sprintf(ERR_URL_PARAM_INVALID, CAMPAIGN_FIELD_ID)))
		} else {
			campaign, err := GetFullCampaign(db, id)
			if err != nil {
				responder.Error(err)
			} else {
				responder.Json(campaign)
			}
		}
	})
}
