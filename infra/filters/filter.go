package filters

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"

	"issuefinder/infra/filters/parser"
)

func NewPredicateParser(text string) FindingPredicate {
	// Setup the input
	is := antlr.NewInputStream(text)

	// Create the Lexer
	lexer := parser.NewPredicateLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	// Create the Parser
	parsr := parser.NewPredicateParser(stream)
	parsr.BuildParseTrees = true
	parsr.RemoveErrorListeners()

	// Create the visitor and walk the tree
	visitor := NewPredicateVisitorImpl()
	tree := parsr.Expr()

	return visitor.Visit(tree).(FindingPredicate)
}
