package utils

import (
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
)

type response struct {
	Data []string `json:"data"`
}

func RunServerTesting() {
	e := echo.New()

	basePath := e.Group("/api")
	basePath.GET("/timeout", geTimedOut)
	basePath.GET("/retry", getRetry)
	basePath.GET("/query-params", getWithQueryParams)

	go e.Start(":8080")
}

func geTimedOut(c echo.Context) error {
	time.Sleep(3 * time.Second)
	return c.JSON(http.StatusOK, nil)
}

var try = 0

func getRetry(c echo.Context) error {
	if try == 3 {
		log.Infof("Successful request on the %drd attempt", try)
		return c.JSON(http.StatusOK, nil)
	}

	if try == 0 {
		log.Info("Initial attempt")
	} else {
		log.Infof("Trying to retry %d", try)
	}

	try += 1
	time.Sleep(3 * time.Second)
	return c.JSON(http.StatusOK, nil)
}

func getWithQueryParams(c echo.Context) error {
	return c.JSON(http.StatusOK, response{
		Data: []string{
			c.QueryParam("email"),
			c.QueryParam("username"),
			c.QueryParam("test"),
		},
	})
}
