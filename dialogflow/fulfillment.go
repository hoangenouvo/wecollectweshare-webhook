package dialogflow

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Platform string

const ActionsOnGoogle Platform = "ACTIONS_ON_GOOGLE"
const Facebook Platform = "FACEBOOK"

// Fulfillment is the response sent back to dialogflow in case of a successful webhook call
type Fulfillment struct {
	FulfillmentText     string             `json:"fulfillmentText,omitempty"`
	FulfillmentMessages Messages           `json:"fulfillmentMessages,omitempty"`
	Source              string             `json:"source,omitempty"`
	Payload             interface{}        `json:"payload,omitempty"`
	OutputContexts      Contexts           `json:"outputContexts,omitempty"`
	FollowupEventInput  FollowupEventInput `json:"followupEventInput,omitempty"`
}

type FacebookPayloadRequest struct {
	Facebook FBRQ `json:"facebook"`
}

type FBRQ struct {
	Text           string   `json:"text"`
	FBQuickReplies QuickRep `json:"quick_replies"`
}
type QuickRep struct {
	ContentType string `json:"content_type"`
}

// FollowupEventInput Optional. Makes the platform immediately invoke another sessions.detectIntent call internally with the specified event as input.
// https://dialogflow.com/docs/reference/api-v2/rest/v2beta1/projects.agent.sessions/detectIntent#EventInput
type FollowupEventInput struct {
	Name         string      `json:"name"`
	LanguageCode string      `json:"languageCode,omitempty"`
	Parameters   interface{} `json:"parameters,omitempty"`
}

// Messages is a simple slice of Message
type Messages []Message

// RichMessage is an interface used in the Message type.
// It is used to send back payloads to dialogflow
type RichMessage interface {
	GetKey() string
}

// Message is a struct holding a platform and a RichMessage.
// Used in the FulfillmentMessages of the response sent back to dialogflow
type Message struct {
	Platform
	RichMessage RichMessage
}

// MarshalJSON implements the Marshaller interface for the JSON type.
// Custom marshalling is necessary since there can only be one rich message
// per Message and the key associated to each type is dynamic
func (m *Message) MarshalJSON() ([]byte, error) {
	var err error
	var b []byte
	buffer := bytes.NewBufferString("{")
	if m.Platform != "" {
		buffer.WriteString(fmt.Sprintf(`"platform": "%s"`, m.Platform))
	}
	if m.Platform != "" && m.RichMessage != nil {
		buffer.WriteString(", ")
	}
	if m.RichMessage != nil {
		if b, err = json.Marshal(m.RichMessage); err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf(`"%s": %s`, m.RichMessage.GetKey(), string(b)))
	}
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// ForGoogle takes a rich message wraps it in a message with the appropriate
// platform set
func ForGoogle(r RichMessage) Message {
	return Message{
		Platform:    ActionsOnGoogle,
		RichMessage: r,
	}
}

func ForFacebook(r RichMessage) Message {
	return Message{
		Platform:    Facebook,
		RichMessage: r,
	}
}

type OriginalDetectIntentRequest struct {
	Source  string      `json:"source,omitempty"`
	Version string      `json:"version,omitempty"`
	Payload PayloadInfo `json:"payload,omitempty"`
}

type FBResponsePayload struct {
	PostBack FBPostBack  `json:"post_back"`
	Sender   interface{} `json:"sender"`
}
type FBPostBack struct {
	Data interface{} `json:"data"`
}

type PayloadInfo struct {
	User              UserInfo    `json:"user,omitempty"`
	Conversation      interface{} `json:"conversation,omitempty"`
	Surface           interface{} `json:"surface,omitempty"`
	Device            DeviceInfo  `json:"device,omitempty"`
	IsInSandbox       interface{} `json:"isInSandbox,omitempty"`
	AvailableSurfaces interface{} `json:"availableSurfaces,omitempty"`
	PostBack          interface{} `json:"postback,omitempty"`
}

type UserInfo struct {
	AccessToken            string   `json:"accessToken,omitempty"`
	Permissions            []string `json:"permissions,omitempty"`
	Locale                 string   `json:"locale,omitempty"`
	LastSeen               string   `json:"lastSeen,omitempty"`
	UserVerificationStatus string   `json:"userVerificationStatus,omitempty"`
}

type DeviceInfo struct {
	LocationInfo LocationInfo `json:"location,omitempty"`
}

type LocationInfo struct {
	Coordinates      Coordinates `json:"coordinates,omitempty"`
	FormattedAddress string      `json:"formattedAddress,omitempty"`
	ZipCode          string      `json:"zipCode,omitempty"`
	City             string      `json:"city,omitempty"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}
