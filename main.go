package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"wcws/dialogflow"
)

var globalSession map[string]string

func init() {
	globalSession = make(map[string]string)
}

func main() {

	_ = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "cred.json")
	_ = godotenv.Load()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Routes
	e.GET("/", test)
	e.POST("/webhook", webhook)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

}

func webhook(e echo.Context) error {
	dr := dialogflow.Request{}
	err := e.Bind(&dr)
	if err != nil {
		log.Println("got err:", err)
		return err
	}
	action := dr.QueryResult.Action
	switch action {
	case "welcome":
		return welcomeHandler(e, dr.OriginalDetectIntentRequest.Source)
	case "collect":
		return addLocationPermissionRequest(e, dr)
	case "getPermission":
		return permissionHander(e, dr)
	}
	return ErrResponse(e)
}

func test(e echo.Context) error {
	return e.String(http.StatusOK, "It's worked!")
}
