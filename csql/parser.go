package csql

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/araddon/dateparse"
)

var aggregationNames = []string{
	"sum",
}

func Parse(tokens []Token) (Expression, int, error) {
	if len(tokens) == 0 {
		return nil, 0, nil
	}

	tok := tokens[0]
	consumed := 0

	var head Expression = &Nop{}
	if tok.Typ == TokenTypeString {
		literal := parseLiteral(tok.Str)
		if literal == nil {
			return nil, consumed, fmt.Errorf("failed to parse literal. string was: %v", tok.Str)
		}
		head = literal
		consumed++
		tokens = tokens[consumed:]
		if len(tokens) > 0 && tokens[0].Typ == TokenTypeLParen {
			argList, consumed2, err := Parse(tokens)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to parse argument list for function call: %w", err)
			}
			if argList.Type() != ExpressionExprList {
				return nil, 0, fmt.Errorf("failed to parse function call: expected expression list but got: %v", argList)
			}
			consumed += consumed2
			if tok.Str == "group" {
				head = &GroupingExpr{
					arguments: *argList.(*ExpressionList),
				}
			} else if slices.Contains(aggregationNames, tok.Str) {
				head = &AggregatingExpr{
					aggregationName: tok.Str,
					argument:        argList.(*ExpressionList).exprs[0],
				}
			} else {
				head = &Funcall{
					funcName:  tok.Str,
					arguments: *argList.(*ExpressionList),
				}
			}
		}
	} else if tok.Typ == TokenTypeOperator {
		if tok.Str == "$" {
			if len(tokens) < 1 {
				return nil, consumed, fmt.Errorf("not enough tokens, expected at least 1 after $ operator")
			}
			nextTok := tokens[1]
			if nextTok.Typ != TokenTypeString {
				return nil, consumed, fmt.Errorf("expected string token after $ operator")
			}
			index, err := strconv.ParseInt(nextTok.Str, 10, 32)
			if err != nil {
				return nil, consumed, fmt.Errorf("failed to parse column reference to int. string was: %v", nextTok.Str)
			}
			head = &ColumnReferenceExpression{
				index: int(index),
			}
			consumed += 2
		} else if tok.Str == "=" {
			consumed += 1
			tokens = tokens[1:]
			rhs, consumed2, err := Parse(tokens)
			if err != nil {
				return nil, 0, err
			}
			head = &OpEquals{
				rhs: rhs,
			}
			consumed += consumed2
		} else if tok.Str == "<" {
			consumed += 1
			tokens = tokens[1:]
			rhs, consumed2, err := Parse(tokens)
			if err != nil {
				return nil, 0, err
			}
			head = &OpLt{
				rhs: rhs,
			}
			consumed += consumed2
		} else if tok.Str == ">" {
			consumed += 1
			tokens = tokens[1:]
			rhs, consumed2, err := Parse(tokens)
			if err != nil {
				return nil, 0, err
			}
			head = &OpGt{
				rhs: rhs,
			}
			consumed += consumed2
		} else if tok.Str == "!" {
			consumed += 1
			tokens = tokens[1:]
			inner, consumed2, err := Parse(tokens)
			if err != nil {
				return nil, 0, err
			}
			head = &OpNeg{
				inner: inner,
			}
			consumed += consumed2
		}
	} else if tok.Typ == TokenTypeLParen {
		exprs := []Expression{}
		consumed += 1
		tokens = tokens[1:]
		for {
			arg, consumed2, err := Parse(tokens)
			if err != nil {
				return nil, 0, err
			}
			exprs = append(exprs, arg)
			consumed += consumed2
			tokens = tokens[consumed2:]
			if len(tokens) == 0 {
				return nil, 0, fmt.Errorf("failed to parse expression list: expected right paren but ran out of tokens")
			}
			if tokens[0].Typ == TokenTypeRParen {
				consumed += 1
				break
			}
			if tokens[0].Typ == TokenTypeComma {
				consumed += 1
				tokens = tokens[1:]
			}
		}
		head = &ExpressionList{
			exprs: exprs,
		}
	}

	return head, consumed, nil
}

func ParseLine(tokens []Token) ([]Expression, int, error) {
	res := []Expression{}
	columnIdx := 0
	consumed := 0
	for len(tokens) > 0 {
		expr, consumed2, err := Parse(tokens)
		if err != nil {
			return nil, 0, err
		}
		consumed += consumed2
		tokens = tokens[consumed2:]

		if len(tokens) > 0 && tokens[0].Typ == TokenTypeOperator {
			expr2, consumed2, err := Parse(tokens)
			if err != nil {
				return nil, 0, err
			}

			if b, ok := expr2.(BinaryExpr); ok {
				b.SetLHS(expr)
				expr = expr2
			} else {
				return nil, 0, fmt.Errorf("expected binary expression but got: %v", expr2)
			}

			consumed += consumed2
			tokens = tokens[consumed2:]
		}
		expr.FillNils(&ColumnReferenceExpression{
			index: columnIdx,
		})
		res = append(res, expr)

		if len(tokens) > 0 {
			if tokens[0].Typ == TokenTypeNewLine {
				return res, consumed, nil
			}
			if tokens[0].Typ != TokenTypeComma {
				return nil, 0, fmt.Errorf("unexpected token type, expected comma but got: %v", tokens[0])
			}
			columnIdx++
			consumed += 1
			tokens = tokens[1:]
		}
	}
	return res, consumed, nil
}

func ParseQuery(tokens []Token) ([][]Expression, error) {
	res := [][]Expression{}
	for len(tokens) > 0 {
		exprs, consumed, err := ParseLine(tokens)
		if err != nil {
			return nil, err
		}
		res = append(res, exprs)
		tokens = tokens[consumed:]
		if len(tokens) > 0 {
			if tokens[0].Typ != TokenTypeNewLine {
				return nil, fmt.Errorf("unexpected token type, expected new line but got: %v", tokens[0])
			}
			tokens = tokens[1:]
		}
	}
	return res, nil
}

func parseLiteral(str string) *LiteralExpression {
	if str == "true" {
		return &LiteralExpression{
			value: Value{
				typ:   ValueTypeBool,
				value: true,
			},
		}
	} else if str == "false" {
		return &LiteralExpression{
			value: Value{
				typ:   ValueTypeBool,
				value: false,
			},
		}
	} else if i, err := strconv.ParseInt(str, 10, 64); err == nil {
		return &LiteralExpression{
			value: Value{
				typ:   ValueTypeInt,
				value: i,
			},
		}
	} else if d, err := strconv.ParseFloat(str, 64); err == nil {
		return &LiteralExpression{
			value: Value{
				typ:   ValueTypeDouble,
				value: d,
			},
		}
	} else if t, err := dateparse.ParseAny(str); err == nil {
		return &LiteralExpression{
			value: Value{
				typ:   ValueTypeDate,
				value: t,
			},
		}

	} else {
		return &LiteralExpression{
			value: Value{
				typ:   ValueTypeString,
				value: str,
			},
		}
	}
}
