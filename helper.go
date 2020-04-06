package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/labstack/echo/v4"

	"wcws/dialogflow"
)

func DoesExistParams(df dialogflow.Request, param string) (bool, error) {
	if len(df.QueryResult.Parameters) > 0 {
		for k, v := range df.QueryResult.Parameters {
			if k == param && v != "" {
				return true, nil
			}
		}
	}
	var contextParams interface{}
	err := df.GetContext("information", &contextParams)
	if err != nil {
		return false, err
	}
	for k, v := range contextParams.(map[string]interface{}) {
		switch v.(type) {
		case string:
			if k == param && v != "" {
				return true, nil
			}
		}
	}
	return false, nil
}

func GetThanksAnswer(name string) string {
	thanksArr := []string{
		fmt.Sprintf("Great! thank you %s", name),
		fmt.Sprintf("Thank you so much %s, have a good day!", name),
		fmt.Sprintf("Thank you very much, %s!", name),
		fmt.Sprintf("Glad you contacted us, %s, thank you!", name),
		fmt.Sprintf("Glad to have you, %s, thanks for doing great things", name),
	}
	tt := len(thanksArr)
	index := rand.Intn(tt)
	return thanksArr[index]
}

func InsertDataToFirebase(e echo.Context, trans Transactions) error {
	ctx := e.Request().Context()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = client.Collection("transactions").Add(ctx, trans)
	if err != nil {
		return err
	}
	return nil
}

func ErrResponse(e echo.Context) error {
	return e.JSON(http.StatusOK, nil)
}
func ExtractAddressFromCoordinator(coordinator dialogflow.Coordinates) (string, error) {
	var payload map[string]interface{}
	ApiKey := os.Getenv("OPENCAGE_API_KEY")
	ApiUrl := fmt.Sprintf("https://api.opencagedata.com/geocode/v1/json?q=%f+%f&key=%s", coordinator.Latitude, coordinator.Longitude, ApiKey)
	res, err := http.Get(ApiUrl)
	if err != nil {
		return "", err
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(resBody, &payload)
	if err != nil {
		return "", err
	}
	address := payload["results"].([]interface{})[0].(map[string]interface{})["formatted"].(string)
	return address, nil
}
