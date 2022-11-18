package handler

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/antgubarev/pingbot/internal/model"
	"github.com/antgubarev/pingbot/internal/storage"
	"github.com/sirupsen/logrus"
)

const SomethingWentWrong string = "Something went wrong"

type StatusHandler struct {
	logger          *logrus.Logger
	responseStorage storage.ResponseStorage
}

func NewStatusHandler(responseStorage storage.ResponseStorage, logger *logrus.Logger) *StatusHandler {
	return &StatusHandler{responseStorage: responseStorage, logger: logger}
}

type TemplateData struct {
	TargetName string
	Items      []model.ResponseTime
	Period     string
}

func (s *StatusHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	targetName := request.URL.Query().Get("target")
	if len(targetName) == 0 {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	periodParam := request.URL.Query().Get("period")
	if len(periodParam) == 0 {
		periodParam = "hour"
	}
	if periodParam != "hour" &&
		periodParam != "day" &&
		periodParam != "week" {
		writer.Write([]byte("Invalid period"))
		return
	}

	var err error
	var data TemplateData
	data.TargetName = targetName

	from, err := getPeriod(periodParam)
	if err != nil {
		_, _ = writer.Write([]byte(SomethingWentWrong))
		s.logger.Errorf("wrong period: %s: %v", from, err)
	}

	data.Items, err = s.responseStorage.GetResponseTimesByPeriod(
		context.Background(),
		from,
		time.Now(),
		targetName)
	if err != nil {
		_, _ = writer.Write([]byte(SomethingWentWrong))
		s.logger.Errorf("parse template: %s", err.Error())
		return
	}

	t, err := template.ParseFiles("template/status.gohtml")
	if err != nil {
		_, _ = writer.Write([]byte(SomethingWentWrong))
		s.logger.Errorf("parse template: %s", err.Error())
		return
	}

	err = t.Execute(writer, data)
	if err != nil {
		writer.Write([]byte(SomethingWentWrong))
		s.logger.Errorf("render template: %s", err.Error())
		return
	}
}

func getPeriod(period string) (from time.Time, err error) {
	switch period {
	case "hour":
		return time.Now().Add(-time.Hour), nil
	case "day":
		return time.Now().Add(-time.Hour * 24), nil
	case "week":
		return time.Now().Add(-time.Hour * 24 * 7), nil
	}
	return time.Time{}, errors.New("invalid period type")
}
