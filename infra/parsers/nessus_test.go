package parsers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"issuefinder/infra/parsers"
)

func TestNessusParser_ParseEmpty(t *testing.T) {
	n := parsers.NewNessusParser()
	doc := "<NessusClientData_v2></NessusClientData_v2>"
	assert.True(t, n.CanHandle(doc))
	issues := n.Parse()
	assert.Equal(t, 0, len(issues))
}

func TestNessusParser_Parse(t *testing.T) {
	n := parsers.NewNessusParser()
	doc := getDoc("Nessus.nessus")
	assert.True(t, n.CanHandle(doc))
	issues := n.Parse()
	assert.Equal(t, 241, len(issues))
}

func TestNessusParser_ParseBig(t *testing.T) {
	n := parsers.NewNessusParser()
	doc := getDoc("Nessus - 5000 findings.nessus")
	assert.True(t, n.CanHandle(doc))
	issues := n.Parse()
	assert.Equal(t, 5784, len(issues))
}
