package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"issuefinder/domain"
	filter "issuefinder/infra/filters"
)

type findings interface {
	GetFindings(context echo.Context) error
}

func (h *APIHandler) GetFindings(context echo.Context) error {

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
		return context.JSON(http.StatusOK, h.buildFindingResult(predfilters))
	}

	return context.JSON(http.StatusOK, h.buildFindingResult(nil))
}

type FindingResult struct {
	Findings      []domain.Finding
	UniqueIps     []string
	UniqueIpPorts []string
}

func (h *APIHandler) buildFindingResult(predfilters []filter.FindingPredicate) FindingResult {
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

func (h *APIHandler) testFilters(finding domain.Finding, predfilters []filter.FindingPredicate) bool {
	for _, filter := range predfilters {
		if !filter.Test(finding) {
			return false
		}
	}
	return true
}

func (h *APIHandler) calcUniqueIps(result FindingResult) []string {
	ips := make(map[string]bool)
	for _, finding := range result.Findings {
		ips[finding.Ip] = true
	}
	unique := make([]string, 0)
	for ip := range ips {
		if ip != "" {
			unique = append(unique, ip)
		}
	}
	return unique
	//return len(ips)
}

func (h *APIHandler) calcUniqueIpPorts(result FindingResult) []string {
	ips := make(map[string]bool)
	for _, finding := range result.Findings {
		ips[fmt.Sprintf("%s:%s", finding.Ip, finding.Port)] = true
	}
	unique := make([]string, 0)
	for ip := range ips {
		if ip != "" {
			unique = append(unique, ip)
		}
	}
	return unique
}
