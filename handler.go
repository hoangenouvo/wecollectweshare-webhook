package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"wcws/dialogflow"
)

func welcomeHandler(e echo.Context, source string) error {
	rs := dialogflow.Fulfillment{}
	answer := "Great! Welcome to We Collect We Share application! Do you have something unused?"
	switch source {
	case "facebook":
		rs = dialogflow.Fulfillment{
			FulfillmentMessages: []dialogflow.Message{
				dialogflow.ForFacebook(dialogflow.TextWrapper{Text: []string{answer}}),
				dialogflow.ForFacebook(dialogflow.Card{
					Title:    "We Collect We Share",
					Subtitle: "Here we collect things from those who want to share to give those in need",
					Image: dialogflow.Image{
						ImageURI:          "https://image.freepik.com/free-vector/volunteers-with-charity-icons-illustration_53876-43180.jpg?fbclid=IwAR2bbsMINLoup2HAG8heP1Kq8KF9oimCDQvcrOXqb14d1VlP8UHFDkEMyNA",
						AccessibilityText: "We Collect We Share",
					},
				}),
			},
		}
	default:
		rs = dialogflow.Fulfillment{
			FulfillmentMessages: []dialogflow.Message{
				dialogflow.ForGoogle(dialogflow.SingleSimpleResponse(answer, answer)),
				dialogflow.ForGoogle(dialogflow.BasicCard{
					Title:         "We Collect We Share",
					FormattedText: "We Collect We Share",
					Image: &dialogflow.Image{
						ImageURI:          "https://image.freepik.com/free-vector/volunteers-with-charity-icons-illustration_53876-43180.jpg?fbclid=IwAR2bbsMINLoup2HAG8heP1Kq8KF9oimCDQvcrOXqb14d1VlP8UHFDkEMyNA",
						AccessibilityText: "We Collect We Share",
					},
				}),
			},
		}
	}
	return e.JSON(http.StatusOK, &rs)
}

func addLocationPermissionRequest(e echo.Context, dr dialogflow.Request) error {
	address, err := DoesExistParams(dr, "address")
	if err != nil {
		return ErrResponse(e)
	}
	any, err := DoesExistParams(dr, "any")
	if err != nil {
		return ErrResponse(e)
	}
	if address == any {
		if dr.OriginalDetectIntentRequest.Source == "facebook" {
			rs := dialogflow.Fulfillment{
				FulfillmentText: "PLACEHOLDER_FOR_PERMISSION",
				Payload: dialogflow.FacebookPayloadRequest{
					Facebook: dialogflow.FBRQ{
						Text: "give me your location please",
						FBQuickReplies: dialogflow.QuickRep{
							ContentType: "location",
						},
					},
				},
			}
			return e.JSON(http.StatusOK, &rs)
		}
		rs := dialogflow.Fulfillment{
			FulfillmentText: "PLACEHOLDER_FOR_PERMISSION",
			Payload: dialogflow.DialogFlowResponseData{
				Google: dialogflow.DialogFlowResponseGoogle{
					ExpectUserResponse: true,
					IsSsml:             false,
					SystemIntent: dialogflow.DialogFlowResponseSystemIntent{
						Intent: "actions.intent.PERMISSION",
						Data: dialogflow.DialogFlowResponseSystemIntentData{
							Type:        "type.googleapis.com/google.actions.v2.PermissionValueSpec",
							OptContext:  "Before I do this",
							Permissions: []string{"DEVICE_PRECISE_LOCATION"},
						},
					},
				},
			},
		}
		return e.JSON(http.StatusOK, &rs)
	}
	return ErrResponse(e)
}

func permissionHander(e echo.Context, dr dialogflow.Request) error {

	address, err := DoesExistParams(dr, "address")
	if err != nil {
		return ErrResponse(e)
	}
	if address {
		trans := Transactions{}
		var dfContext map[string]interface{}
		err := dr.GetContext("information", &dfContext)
		if err != nil {
			return ErrResponse(e)
		}
		if dr.OriginalDetectIntentRequest.Source == "facebook" {
			userLocation := dr.OriginalDetectIntentRequest.Payload.PostBack
			lat, err := strconv.ParseFloat(userLocation.(map[string]interface{})["data"].(map[string]interface{})["lat"].(string), 64)
			if err != nil {
				return ErrResponse(e)
			}
			long, err := strconv.ParseFloat(userLocation.(map[string]interface{})["data"].(map[string]interface{})["long"].(string), 64)
			if err != nil {
				return ErrResponse(e)
			}
			coordinates := dialogflow.Coordinates{
				Latitude:  lat,
				Longitude: long,
			}
			address, err := ExtractAddressFromCoordinator(coordinates)
			trans = Transactions{
				Description: dfContext["description"].(string),
				GiverName:   dfContext["person"].(map[string]interface{})["name"].(string),
				PhoneNumber: dfContext["phone-number"].(string),
				Address:     address,
				Long:        coordinates.Longitude,
				Lat:         coordinates.Latitude,
				CreatedDate: time.Now().Unix(),
				Status:      "pending",
				TransactionTime: func() int64 {
					origin := dfContext["transaction-time.original"].(string)
					t, _ := time.Parse(time.RFC3339, origin)
					return t.Unix()
				}(),
			}
		} else {
			userLocation := dr.OriginalDetectIntentRequest.Payload.Device.LocationInfo
			address, err := ExtractAddressFromCoordinator(userLocation.Coordinates)
			if err != nil {
				return ErrResponse(e)
			}
			trans = Transactions{
				Description:     dfContext["any"].(string),
				GiverName:       dfContext["person"].(map[string]interface{})["name"].(string),
				PhoneNumber:     dfContext["phone-number"].(string),
				Address:         address,
				Long:            userLocation.Coordinates.Longitude,
				Lat:             userLocation.Coordinates.Latitude,
				CreatedDate:     time.Now().Unix(),
				Status:          "pending",
				TransactionTime: time.Now().Unix(),
			}
		}
		if err := InsertDataToFirebase(e, trans); err != nil {
			return ErrResponse(e)
		}
		thanksAnswer := GetThanksAnswer(trans.GiverName)
		rs := dialogflow.Fulfillment{
			FulfillmentMessages: func() []dialogflow.Message {
				if dr.OriginalDetectIntentRequest.Source == "facebook" {
					return []dialogflow.Message{
						dialogflow.ForFacebook(dialogflow.TextWrapper{Text: []string{thanksAnswer}}),
					}
				}
				return []dialogflow.Message{
					dialogflow.ForGoogle(dialogflow.SingleSimpleResponse(thanksAnswer, thanksAnswer)),
				}
			}(),
			OutputContexts: func() dialogflow.Contexts {
				for k, v := range dr.QueryResult.OutputContexts {
					dr.QueryResult.OutputContexts[k].Name = v.Name
					dr.QueryResult.OutputContexts[k].LifespanCount = 0
					dr.QueryResult.OutputContexts[k].Parameters = nil
				}
				return dr.QueryResult.OutputContexts
			}(),
		}
		return e.JSON(http.StatusOK, rs)
	}
	return ErrResponse(e)
}
