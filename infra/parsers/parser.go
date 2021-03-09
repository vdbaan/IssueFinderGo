package parsers

import (
	"encoding/xml"
	"io"
	"strings"

	"issuefinder/domain"
)

type Parser interface {
	GetName() string
	CanHandle(doc string) bool
	Parse() []domain.Finding
}

var (
	parsers = []Parser{NewNessusParser()}
)

func GetParser(doc string) Parser {
	for _, parser := range parsers {
		if parser.CanHandle(doc) {
			return parser
		}
	}
	return nil
}

func getNameRootToken(xmldoc string) string {
	d := xml.NewDecoder(strings.NewReader(xmldoc))
	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			// EOF means we're done.
			break
		}
		switch ty := tok.(type) {
		case xml.StartElement:
			return ty.Name.Local
		}
	}

	return ""
}
