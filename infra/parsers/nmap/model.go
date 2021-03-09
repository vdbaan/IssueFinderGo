package nmap

import "encoding/xml"

type Nmaprun struct {
	XMLName  xml.Name `xml:"nmaprun"`
	Args     string   `xml:"args,attr"`
	Scaninfo Scaninfo `xml:"scaninfo"`
	Host     []Host   `xml:"host"`
	Runstats Runstats `xml:"runstats"`
}

type Scaninfo struct {
	Text        string `xml:",chardata"`
	Type        string `xml:"type,attr"`
	Protocol    string `xml:"protocol,attr"`
	Numservices string `xml:"numservices,attr"`
	Services    string `xml:"services,attr"`
}

type Host struct {
	Address   []Address `xml:"address"`
	Hostnames Hostnames `xml:"hostnames"`
	Ports     Ports     `xml:"ports"`
}

type Address struct {
	Addr     string `xml:"addr,attr"`
	Addrtype string `xml:"addrtype,attr"`
}

type Hostnames struct {
	Hostname []Hostname `xml:"hostname"`
}

type Hostname struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type Ports struct {
	Text       string     `xml:",chardata"`
	Extraports Extraports `xml:"extraports"`
	Port       []Port     `xml:"port"`
}

type Port struct {
	Protocol string   `xml:"protocol,attr"`
	Portid   string   `xml:"portid,attr"`
	State    State    `xml:"state"`
	Service  Service  `xml:"service"`
	Script   []Script `xml:"script"`
}

type State struct {
	State string `xml:"state,attr"`
}

type Service struct {
	Name    string `xml:"name,attr"`
	Product string `xml:"product,attr"`
}

type Script struct {
	Text   string `xml:",chardata"`
	ID     string `xml:"id,attr"`
	Output string `xml:"output,attr"`
}

type Extraports struct {
	State string `xml:"state,attr"`
	Count string `xml:"count,attr"`
}

type Runstats struct {
	Text     string   `xml:",chardata"`
	Finished Finished `xml:"finished"`
	Hosts    Hosts    `xml:"hosts"`
}

type Finished struct {
	Text    string `xml:",chardata"`
	Time    string `xml:"time,attr"`
	Timestr string `xml:"timestr,attr"`
	Elapsed string `xml:"elapsed,attr"`
	Summary string `xml:"summary,attr"`
	Exit    string `xml:"exit,attr"`
}

type Hosts struct {
	Text  string `xml:",chardata"`
	Up    string `xml:"up,attr"`
	Down  string `xml:"down,attr"`
	Total string `xml:"total,attr"`
}
