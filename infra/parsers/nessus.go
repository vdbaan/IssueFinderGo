package parsers

import (
	"encoding/xml"
	"fmt"
	"strings"

	"issuefinder/domain"
	"issuefinder/infra/parsers/nessus"
)

type nessusParser struct {
	Name string
	Data nessus.NessusClientData_v2
}

func NewNessusParser() Parser {
	p := new(nessusParser)
	p.Name = "Nessus"
	p.Data = nessus.NessusClientData_v2{}
	return p
}

func (p *nessusParser) GetName() string {
	return p.Name
}

func (p *nessusParser) CanHandle(doc string) bool {
	if getNameRootToken(doc) != "NessusClientData_v2" {
		return false
	}
	if err := xml.Unmarshal([]byte(doc), &p.Data); err != nil {
		return false
	}
	return true
}

func (p *nessusParser) Parse() []domain.Finding {
	result := make([]domain.Finding, 0)

	for _, host := range p.Data.Report.ReportHost {
		tags := p.mapTags(host)
		hostName := tags["host-ip"]
		hostIp := hostName
		if fqdn, ok := tags["host-fqdn"]; ok {
			hostName = fqdn
		}
		for _, issue := range host.ReportItem {
			risk := domain.Unknown
			switch issue.Severity {
			case "0":
				risk = domain.Info
			case "1":
				risk = domain.Low
			case "2":
				risk = domain.Medium
			case "3":
				risk = domain.High
			case "4":
				risk = domain.Critical
			}

			references := strings.Builder{}
			if issue.Cve != nil {
				references.WriteString("CVE references  : ")
				references.WriteString(strings.Join(issue.Cve, ", "))
				references.WriteByte('\n')
			}
			if issue.Bid != nil {
				references.WriteString("BID references  : ")
				references.WriteString(strings.Join(issue.Bid, ", "))
				references.WriteByte('\n')
			}
			if issue.Xref != nil {
				references.WriteString("Other references: ")
				references.WriteString(strings.Join(issue.Xref, ", "))
				references.WriteByte('\n')
			}

			summary := strings.Builder{}
			summary.WriteString(issue.Synopsis)
			summary.WriteString("\n")
			summary.WriteString(fmt.Sprintf("RiskFactor           : %s\n", issue.RiskFactor))
			summary.WriteString(fmt.Sprintf("Exploit available    : %v\n", issue.ExploitAvailable))
			summary.WriteString(fmt.Sprintf("Ease of exploit      : %s\n", issue.EaseOfExploit))
			summary.WriteString(fmt.Sprintf("Patch availabls since: %s\n", issue.PatchDate))
			summary.WriteString(fmt.Sprintf("CVSS temporal vector : %s\n", issue.CvssTemporalVector))
			summary.WriteString(fmt.Sprintf("CVSS temporal score  : %2.1f\n", issue.CvssTemporalScore))

			result = append(result, domain.Finding{
				Scanner:      p.Name,
				Ip:           hostIp,
				Hostname:     hostName,
				Port:         issue.Port,
				PortStatus:   "open",
				Protocol:     issue.Protocol,
				Service:      issue.Service,
				Plugin:       fmt.Sprintf("%s:%s", issue.PluginID, issue.PluginName),
				Summary:      summary.String(),
				Description:  issue.Description,
				Reference:    references.String(),
				PluginOutput: issue.PluginOutput,
				Solution:     issue.Solution,
				CvssVector:   issue.CvssVector,
				Severity:     risk,
				Exploitable:  issue.ExploitAvailable,
				BaseCvss:     issue.CvssBaseScore,
			})
		}
	}

	return result
}

func (p *nessusParser) mapTags(host nessus.ReportHost) map[string]string {
	result := make(map[string]string)
	for _, tag := range host.HostProperties.Tag {
		result[tag.Name] = tag.Value
	}
	return result
}
