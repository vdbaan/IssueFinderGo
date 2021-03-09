package filters

import (
	"fmt"
	"strconv"
	"strings"

	"issuefinder/domain"
)

type LogicalOperation string
type ColumnName string

const (
	LT      LogicalOperation = "<"
	LE      LogicalOperation = "<="
	GT      LogicalOperation = ">"
	GE      LogicalOperation = ">="
	EQ      LogicalOperation = "=="
	NE      LogicalOperation = "!="
	LIKE    LogicalOperation = "LIKE"
	APROX   LogicalOperation = "~="
	NOT     LogicalOperation = "!"
	AND     LogicalOperation = "&&"
	OR      LogicalOperation = "||"
	IN      LogicalOperation = "IN"
	BETWEEN LogicalOperation = "BETWEEN"
	NLIKE   LogicalOperation = "NOT LIKE"

	SCANNER     ColumnName = "SCANNER"
	IP          ColumnName = "IP"
	PORT        ColumnName = "PORT"
	SERVICE     ColumnName = "SERVICE"
	RISK        ColumnName = "RISK"
	EXPLOITABLE ColumnName = "EXPLOITABLE"
	DESCRIPTION ColumnName = "DESCRIPTION"
	PLUGIN      ColumnName = "PLUGIN"
	STATUS      ColumnName = "STATUS"
	PROTOCOL    ColumnName = "PROTOCOL"
	HOSTNAME    ColumnName = "HOSTNAME"
	CVSS        ColumnName = "CVSS"
)

func (c ColumnName) get(finding domain.Finding) interface{} {
	switch c {
	case SCANNER:
		return finding.Scanner
	case IP:
		return finding.Ip
	case PORT:
		return finding.Port
	case SERVICE:
		return finding.Service
	case RISK:
		return finding.Severity
	case EXPLOITABLE:
		return finding.Exploitable
	case DESCRIPTION:
		return finding.Description
	case STATUS:
		return finding.PortStatus
	case PROTOCOL:
		return finding.Protocol
	case HOSTNAME:
		return finding.Hostname
	case CVSS:
		return finding.BaseCvss
	default:
		return nil
	}
}

type FindingPredicate struct {
	Left      interface{}
	Operation LogicalOperation
	Right     interface{}
}

func (fp *FindingPredicate) Test(finding domain.Finding) bool {
	lval := fp.Left
	rval := fp.Right

	if _, ok := fp.Left.(ColumnName); ok {
		lval = fp.Left.(ColumnName).get(finding)
	}
	if _, ok := fp.Right.(ColumnName); ok {
		rval = fp.Right.(ColumnName).get(finding)
	}

	// if we compare the RISK column with a string ( e.g. RISK >= LOW
	if _, ok := fp.Left.(ColumnName); ok {
		_, isString := fp.Right.(string)
		if fp.Left.(ColumnName) == RISK && isString {
			switch rval {
			case "LOW":
				rval = domain.Low
			case "MEDIUM":
				rval = domain.Medium
			case "HIGH":
				rval = domain.High
			case "CRITICAL":
				rval = domain.Critical
			default:
				rval = domain.Unknown
			}
		}
	}
	if lval == nil || (rval == nil && fp.Operation != NOT) {
		return false
	}

	switch fp.Operation {
	case EQ:
		return lval == rval
	case NE:
		return lval != rval
	case LIKE:
		return strings.Contains(fmt.Sprintf("%v", lval), fmt.Sprintf("%v", rval))
	case NOT:
		if _, ok := lval.(FindingPredicate); ok {
			l := lval.(FindingPredicate)
			return !l.Test(finding)
		}
		return !lval.(bool)

	case AND:
		l := lval.(FindingPredicate)
		return l.and(rval.(FindingPredicate), finding)
	case OR:
		l := lval.(FindingPredicate)
		return l.or(rval.(FindingPredicate), finding)
	case LT:
		return fp.compare(lval, rval) < 0
	case LE:
		return fp.compare(lval, rval) <= 0
	case GT:
		return fp.compare(lval, rval) > 0
	case GE:
		return fp.compare(lval, rval) >= 0
	}

	return false
}

func (fp *FindingPredicate) and(other FindingPredicate, finding domain.Finding) bool {
	return fp.Test(finding) && other.Test(finding)
}

func (fp *FindingPredicate) or(other FindingPredicate, finding domain.Finding) bool {
	return fp.Test(finding) || other.Test(finding)
}

func (fp *FindingPredicate) compare(lval interface{}, rval interface{}) int {
	if _, ok := lval.(domain.Severity); ok {
		l := lval.(domain.Severity)
		r := rval.(domain.Severity)
		result := l - r
		return int(result)
	}
	if _, ok := lval.(float32); ok {

		if r, err := strconv.ParseFloat(rval.(string), 32); err == nil {
			l := float64(lval.(float32))
			return int(l - r)
		}

	}
	if _, ok := lval.(int); ok {
		l := lval.(int)
		r := rval.(int)
		return l - r
	}

	return strings.Compare(fmt.Sprintf("%v", lval), fmt.Sprintf("%v", rval))
}
