package filters

import (
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"issuefinder/infra/filters/parser"
)

type PredicateVisitorImpl struct {
	parser.BasePredicateVisitor
}

func NewPredicateVisitorImpl() parser.PredicateVisitor {
	return PredicateVisitorImpl{}
}

func (f PredicateVisitorImpl) Visit(tree antlr.ParseTree) interface{} {
	switch val := tree.(type) {
	case *parser.AndExprContext:
		return f.VisitAndExpr(val)
	case *parser.AssignContext:
		return f.VisitAssign(val)
	case *parser.EnclosedExprContext:
		return f.VisitEnclosedExpr(val)
	case *parser.ColumnContext:
		return f.VisitColumn(val)
	case *parser.ExploitableExprContext:
		return f.VisitExploitableExpr(val)
	case *parser.GroupContext:
		return f.VisitGroup(val)
	case *parser.GroupOperatorContext:
		return f.VisitGroupOperator(val)
	case *parser.OrExprContext:
		return f.VisitOrExpr(val)
	case *parser.OperatorContext:
		return f.VisitOperator(val)
	}
	return nil
}

func (f PredicateVisitorImpl) VisitEnclosedExpr(ctx *parser.EnclosedExprContext) interface{} { // LPAREN expr RPAREN
	return f.Visit(ctx.Expr())
}

func (f PredicateVisitorImpl) VisitRange(ctx *parser.RangeContext) interface{} { // column rangeOperator RANGE
	return FindingPredicate{
		Left:      f.Visit(ctx.Column()),
		Operation: f.Visit(ctx.RangeOperator()).(LogicalOperation),
		Right:     f.buildList(ctx.GROUP().GetText()),
	}
}

func (f PredicateVisitorImpl) VisitOrExpr(ctx *parser.OrExprContext) interface{} { // expr OR expr
	return FindingPredicate{
		Left:      f.Visit(ctx.Expr(0)),
		Operation: OR,
		Right:     f.Visit(ctx.Expr(1)),
	}
}

func (f PredicateVisitorImpl) VisitExploitableExpr(ctx *parser.ExploitableExprContext) interface{} { //EXPLOITABLE
	result := FindingPredicate{Left: EXPLOITABLE}
	if ctx.GetChildCount() == 2 {
		result.Operation = NOT
	}
	return result
}

func (f PredicateVisitorImpl) VisitAssign(ctx *parser.AssignContext) interface{} { // column operator STRING
	operation := f.Visit(ctx.Operator())
	right := f.stripQuotes(ctx.STRING().GetText())
	return FindingPredicate{
		Left:      ColumnName(f.Visit(ctx.Column()).(string)),
		Operation: LogicalOperation(operation.(string)),
		Right:     right}
}

func (f PredicateVisitorImpl) VisitGroup(ctx *parser.GroupContext) interface{} { // column groupOperator GROUP
	return FindingPredicate{
		Left:      f.Visit(ctx.Column()),
		Operation: f.Visit(ctx.GroupOperator()).(LogicalOperation),
		Right:     f.buildList(ctx.GROUP().GetText()),
	}
}

func (f PredicateVisitorImpl) VisitAndExpr(ctx *parser.AndExprContext) interface{} { // expr AND expr
	return FindingPredicate{
		Left:      f.Visit(ctx.Expr(0)),
		Operation: AND,
		Right:     f.Visit(ctx.Expr(1)),
	}
}

func (f PredicateVisitorImpl) VisitColumn(ctx *parser.ColumnContext) interface{} {
	return strings.ToUpper(ctx.GetText())
}

func (f PredicateVisitorImpl) VisitGroupOperator(ctx *parser.GroupOperatorContext) interface{} {
	return strings.ToUpper(ctx.GetText())
}

func (f PredicateVisitorImpl) VisitRangeOperator(ctx *parser.RangeOperatorContext) interface{} {
	return strings.ToUpper(ctx.GetText())
}

func (f PredicateVisitorImpl) VisitOperator(ctx *parser.OperatorContext) interface{} {
	if ctx.GetChildCount() == 2 {
		return NLIKE
	}
	return strings.ToUpper(ctx.GetText())
}

func (f PredicateVisitorImpl) stripQuotes(input string) string {
	if strings.HasPrefix(input, "'") || strings.HasPrefix(input, "\"") {
		sb := strings.Builder{}
		l := len(input)
		sb.WriteString(input[1 : l-1])
		return sb.String()
	}
	return input
}

func (f PredicateVisitorImpl) buildList(text string) []string {
	result := make([]string, 0)
	l := len(text)
	workable := text[1 : l-1]
	for _, s := range strings.Split(workable, ",") {
		result = append(result, f.stripQuotes(s))
	}
	return result
}
