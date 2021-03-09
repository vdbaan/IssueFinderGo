package filters_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"issuefinder/domain"
	"issuefinder/infra/filters"
)

func TestNewPredicateParser(t *testing.T) {
	f := domain.Finding{
		Ip: "127.0.0.1",
	}
	fp := filters.NewPredicateParser("IP == 127.0.0.1")
	assert.Equal(t, fp.Left, filters.IP)
	assert.Equal(t, fp.Operation, filters.EQ)
	assert.Equal(t, fp.Right, "127.0.0.1")
	assert.True(t, fp.Test(f))
}

func TestRiskPredicate(t *testing.T) {
	f := domain.Finding{
		Ip:       "127.0.0.1",
		Severity: domain.Critical,
		BaseCvss: 3.4,
	}
	fp := filters.NewPredicateParser("risk >= LOW")
	assert.Equal(t, fp.Left, filters.RISK)
	assert.True(t, fp.Test(f))

	fp = filters.NewPredicateParser("CVSS >= 2.0")
	assert.Equal(t, fp.Left, filters.CVSS)
	assert.True(t, fp.Test(f))
}
