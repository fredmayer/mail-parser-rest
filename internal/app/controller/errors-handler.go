package controller

import (
	"errors"
	"net/http"

	"github.com/fredmayer/mail-parser-rest/pkg/logging"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func notFoundError(path string, err error, params map[string]interface{}) *echo.HTTPError {
	status := http.StatusBadRequest
	fields := logrus.Fields{
		"status": status,
		"msg":    err.Error(),
	}

	for k, val := range params {
		fields[k] = val
	}

	logging.Log().WithFields(fields).Warn(path)
	return echo.NewHTTPError(status, errors.New("not found"))
}

func badRequestError(path string, err error, params map[string]interface{}) *echo.HTTPError {
	status := http.StatusBadRequest
	fields := logrus.Fields{
		"status": status,
		"msg":    err.Error(),
	}
	for k, val := range params {
		fields[k] = val
	}

	logging.Log().WithFields(fields).Warn(path)
	return echo.NewHTTPError(status, err)
}
