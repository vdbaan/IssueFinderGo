package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"

	"issuefinder/domain"
	"issuefinder/infra/config"
	"issuefinder/infra/parsers"
)

type API interface {
	UploadFile(context echo.Context) error
	Reset(context echo.Context) error
	findings
	filters
}
type APIHandler struct {
	Config   config.Handler
	Filters  []string
	Findings []domain.Finding
}

func NewAPIHandler(cfg config.Handler) API {
	h := new(APIHandler)
	h.Config = cfg
	h.Findings = make([]domain.Finding, 0)
	return h
}

func (h *APIHandler) UploadFile(context echo.Context) error {
	if file, err := context.FormFile("file1"); err == nil {
		src, err := file.Open()
		if err != nil {
			return context.JSON(http.StatusInternalServerError, nil)
		}
		defer src.Close()

		var contents []byte
		contents, err = ioutil.ReadAll(src)
		if err != nil {
			return context.JSON(http.StatusInternalServerError, nil)
		}

		parser := parsers.GetParser(string(contents))
		h.Findings = append(h.Findings, parser.Parse()...)

		return context.JSON(http.StatusOK, parser.GetName())
	}
	return context.JSON(http.StatusBadRequest, nil)
}

func (h *APIHandler) getJSONmap(context echo.Context) (map[string]interface{}, error) {
	json_map := make(map[string]interface{})
	err := json.NewDecoder(context.Request().Body).Decode(&json_map)
	if err != nil {
		return nil, err
	}
	return json_map, nil
}

func (h *APIHandler) Reset(context echo.Context) error {
	h.Findings = nil
	return context.JSON(http.StatusOK, nil)
}
