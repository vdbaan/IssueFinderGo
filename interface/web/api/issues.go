package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"issuefinder/domain"
	filter "issuefinder/infra/filters"
)

func (h *APIHandler) GetIssues(context echo.Context) error {
	dataFilter := make([]string, 0)
	predfilters := make([]filter.FindingPredicate, 0)
	if req, err := h.getJSONmap(context); err == nil {
		filters := req["filters"].([]interface{})
		for _, f := range filters {
			fm := f.(map[string]interface{})
			val := fm["text"].(string)
			dataFilter = append(dataFilter, val)
			predfilters = append(predfilters, filter.NewPredicateParser(val))
		}
		return context.JSON(http.StatusOK, h.buildIssueResult(predfilters))
	}

	return context.JSON(http.StatusOK, h.buildIssueResult(nil))
}

func (h *APIHandler) buildIssueResult(predfilters []filter.FindingPredicate) FindingResult {
	result := FindingResult{}
	if predfilters == nil {
		result.Findings = h.Findings
	} else {
		result.Findings = make([]domain.Finding, 0)
		for _, finding := range h.Findings {
			if h.testFilters(finding, predfilters) {
				result.Findings = append(result.Findings, finding)
			}
		}
	}
	result.UniqueIps = h.calcUniqueIps(result)
	result.UniqueIpPorts = h.calcUniqueIpPorts(result)
	return result
}
