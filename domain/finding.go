package domain

type Severity int

const (
	Unknown Severity = iota
	Info
	Low
	Medium
	High
	Critical
)

type Finding struct {
	Scanner      string
	Ip           string
	Hostname     string
	Port         string
	PortStatus   string
	Protocol     string
	Location     string
	Service      string
	Plugin       string
	Summary      string
	Description  string
	Reference    string
	PluginOutput string
	Solution     string
	CvssVector   string
	Severity     Severity `json:"severity"` // need to be lowercase in order for Vue datatable to work with colour
	Exploitable  bool
	BaseCvss     float32
}
