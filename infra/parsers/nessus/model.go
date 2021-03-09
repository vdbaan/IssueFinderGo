package nessus

// I only add those elements into the struct that I am interested in :D
type NessusClientData_v2 struct {
	Report Report
}

type Report struct {
	Name       string `xml:"name,attr"`
	ReportHost []ReportHost
}

type ReportHost struct {
	Name           string `xml:"name,attr"`
	HostProperties HostProperties
	ReportItem     []ReportItem
}

type HostProperties struct {
	Tag []Tag `xml:"tag"`
}

type Tag struct {
	Value string `xml:",chardata"`
	Name  string `xml:"name,attr"`
}

type ReportItem struct {
	Port               string   `xml:"port,attr"`
	Protocol           string   `xml:"protocol,attr"`
	PluginID           string   `xml:"pluginID,attr"`
	PluginName         string   `xml:"pluginName,attr"`
	Service            string   `xml:"svc_name,attr"`
	Severity           string   `xml:"severity,attr"`
	Description        string   `xml:"description"`
	Solution           string   `xml:"solution"`
	PluginOutput       string   `xml:"plugin_output"`
	CvssVector         string   `xml:"cvss_vector"`
	CvssBaseScore      float32  `xml:"cvss_base_score"`
	Synopsis           string   `xml:"synopsis"`
	Cve                []string `xml:"cve"`
	Bid                []string `xml:"bid"`
	Xref               []string `xml:"xref"`
	RiskFactor         string   `xml:"risk_factor"`
	ExploitAvailable   bool     `xml:"exploit_available"`
	EaseOfExploit      string   `xml:"exploitability_ease"`
	PatchDate          string   `xml:"patch_publication_date"`
	CvssTemporalVector string   `xml:"cvss_temporal_vector"`
	CvssTemporalScore  float32  `xml:"cvss_temporal_score"`
}
