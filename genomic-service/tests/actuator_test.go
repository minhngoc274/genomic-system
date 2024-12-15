package tests

import (
	"net/http"
	"testing"
	golibtest "github.com/golibs-starter/golib-test"
)

func TestActuatorInfo_ShouldSuccess(t *testing.T) {
	golibtest.NewRestAssured(t).
		When().
		Get("/actuator/info").
		Then().
		Status(http.StatusOK).
		Body("meta.code", 200).
		Body("data.service_name", "Application Name")
}

func TestActuatorHealth_ShouldSuccess(t *testing.T) {
	golibtest.NewRestAssured(t).
		When().
		Get("/actuator/health").
		Then().
		Status(http.StatusOK).
		Body("meta.code", 200).
		Body("data.status", "UP")
}
