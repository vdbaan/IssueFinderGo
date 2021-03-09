package parsers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"issuefinder/infra/parsers"
)

func TestNMapParser_ParseEmpty(t *testing.T) {
	n := parsers.NewNMapParser()
	doc := "<nmaprun></nmaprun>"
	assert.True(t, n.CanHandle(doc))
	issues := n.Parse()
	assert.Equal(t, 0, len(issues))
}

func TestNMapParser_Parse(t *testing.T) {
	n := parsers.NewNMapParser()
	doc := getDoc("Nmap.xml")
	assert.True(t, n.CanHandle(doc))
	issues := n.Parse()
	assert.Equal(t, 54, len(issues))
}
