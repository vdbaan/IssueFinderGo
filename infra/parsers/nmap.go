package parsers

import (
	"encoding/xml"
	"fmt"
	"strings"

	"issuefinder/domain"
	"issuefinder/infra/parsers/nmap"
)

type NMapParser struct {
	Name string
	Data nmap.Nmaprun
}

func NewNMapParser() Parser {
	p := new(NMapParser)
	p.Name = "NMap"
	p.Data = nmap.Nmaprun{}
	return p
}

func (p *NMapParser) GetName() string {
	return p.Name
}

func (p *NMapParser) CanHandle(doc string) bool {
	if getNameRootToken(doc) != "nmaprun" {
		return false
	}
	if err := xml.Unmarshal([]byte(doc), &p.Data); err != nil {
		return false
	}
	return true
}

func (p *NMapParser) Parse() []domain.Finding {
	result := make([]domain.Finding, 0)

	protocol := p.Data.Scaninfo.Protocol
	for _, host := range p.Data.Host {
		hostIp := ""
		for _, a := range host.Address {
			if a.Addrtype == "ipv4" {
				hostIp = a.Addr
			}
		}
		hostName := ""
		for _, hn := range host.Hostnames.Hostname {
			if hn.Type == "user" || hn.Type == "PTR" {
				hostName = hn.Name
			} else {
				hostName = hostIp
			}
		}
		result = append(result, p.scanInfo(hostIp, hostName, protocol))
		result = append(result, p.runStats(hostIp, hostName, protocol))
		result = append(result, p.summary(hostIp, hostName, protocol))

		if host.Ports.Extraports.State == "closed" {
			result = append(result, domain.Finding{
				Scanner:    p.Name,
				Ip:         hostIp,
				Hostname:   hostName,
				Port:       "0",
				PortStatus: "closed",
				Protocol:   protocol,
				Service:    "none",
				Plugin:     "NMap closed ports",
				Summary:    fmt.Sprintf("Amount of closed ports: %s", host.Ports.Extraports.Count),
				Severity:   domain.Medium,
			})
		}

		for _, port := range host.Ports.Port {
			prod := port.Service.Product
			service := port.Service.Name
			if prod != "" {
				service = fmt.Sprintf("%s (%s)", service, prod)
			}

			risk := domain.Info
			if port.State.State == "closed" {
				risk = domain.Low
			}

			summary := strings.Builder{}
			for _, script := range port.Script {
				summary.WriteString(fmt.Sprintf("script (%s):\n%s\n%s\n\n", script.ID, script.Output, script.Text))
			}

			result = append(result, domain.Finding{
				Scanner:    p.Name,
				Ip:         hostIp,
				Hostname:   hostName,
				Port:       port.Portid,
				PortStatus: port.State.State,
				Protocol:   port.Protocol,
				Service:    service,
				Plugin:     fmt.Sprintf("NMap port (%s) information", port.Portid),
				Summary:    summary.String(),
				Severity:   risk,
			})
		}
	}

	return result
}

func (p *NMapParser) scanInfo(ip string, name string, protocol string) domain.Finding {
	summary := strings.Builder{}
	summary.WriteString(fmt.Sprintf("Protocol       : %s\n", protocol))
	summary.WriteString(fmt.Sprintf("Number of ports: %s\n", p.Data.Scaninfo.Numservices))
	summary.WriteString(fmt.Sprintf("Ports scanned  : %s\n", p.Data.Scaninfo.Services))
	summary.WriteString(fmt.Sprintf("Scan type      : %s\n", p.Data.Scaninfo.Type))
	summary.WriteString(fmt.Sprintf("Nmap command   : %s\n", p.Data.Args))
	return domain.Finding{
		Scanner:    p.Name,
		Ip:         ip,
		Hostname:   name,
		Port:       "0",
		PortStatus: "NA",
		Protocol:   protocol,
		Service:    "none",
		Plugin:     "NMap scan info",
		Summary:    summary.String(),
		Severity:   domain.Info,
	}
}

func (p *NMapParser) runStats(ip string, name string, protocol string) domain.Finding {
	summary := strings.Builder{}
	summary.WriteString("Number of hosts\n")
	summary.WriteString(fmt.Sprintf("Scanned: %s\n", p.Data.Runstats.Hosts.Total))
	summary.WriteString(fmt.Sprintf("Up     : %s\n", p.Data.Runstats.Hosts.Up))
	summary.WriteString(fmt.Sprintf("Down   : %s\n", p.Data.Runstats.Hosts.Down))
	return domain.Finding{
		Scanner:    p.Name,
		Ip:         ip,
		Hostname:   name,
		Port:       "0",
		PortStatus: "NA",
		Protocol:   protocol,
		Service:    "none",
		Plugin:     "NMap scan info",
		Summary:    summary.String(),
		Severity:   domain.Info,
	}
}

func (p *NMapParser) summary(ip string, name string, protocol string) domain.Finding {
	summary := strings.Builder{}
	if p.Data.Runstats.Finished.Summary != "" {
		summary.WriteString(p.Data.Runstats.Finished.Summary)
	} else {
		summary.WriteString("Scan Execution Stats\n")
		summary.WriteString(fmt.Sprintf("Completed: %s\n", p.Data.Runstats.Finished.Timestr))
		summary.WriteString(fmt.Sprintf("Duration : %s\n", p.Data.Runstats.Finished.Elapsed))
	}
	return domain.Finding{
		Scanner:    p.Name,
		Ip:         ip,
		Hostname:   name,
		Port:       "0",
		PortStatus: "NA",
		Protocol:   protocol,
		Service:    "none",
		Plugin:     "NMap scan info",
		Summary:    summary.String(),
		Severity:   domain.Info,
	}
}
