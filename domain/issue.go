package domain

type Issue struct {
	Title           string
	Severity        Severity
	Description     string
	Recommendations string
	Location        string
}
