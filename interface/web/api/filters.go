package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type filters interface {
	GetFilters(context echo.Context) error
	Addfilter(context echo.Context) error
	ShowSettings(context echo.Context) error
}

func (h *APIHandler) GetFilters(context echo.Context) error {
	return context.JSON(http.StatusOK, h.Filters)
}

func (h *APIHandler) Addfilter(context echo.Context) error {
	if jmap, err := h.getJSONmap(context); err == nil {
		filter := jmap["text"].(string)
		cfg := h.Config.GetConfig()
		cfg.Filters = append(cfg.Filters, filter)
		h.Config.UpdateConfig(cfg)
		return context.JSON(http.StatusCreated, nil)
	}
	return context.JSON(http.StatusBadRequest, nil)
}

func (h *APIHandler) ShowSettings(context echo.Context) error {
	return context.JSON(http.StatusOK, nil)
}
