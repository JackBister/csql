package csql

import (
	"fmt"
	"strconv"

	"github.com/araddon/dateparse"
)

func Parse(tokens []Token) [][]Expression {
	fmt.Println(tokens)
	exprs := [][]Expression{}
	currentExprs := []Expression{}
	currentColumn := 0
	var lastExpr Expression = &Nop{}
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		var currentExpr Expression = nil
		if t.Typ == TokenTypeComma {
			currentExprs = append(currentExprs, lastExpr)
			currentColumn++
			lastExpr = &Nop{}
		} else if t.Typ == TokenTypeOperator {
			if len(tokens) < i+2 {
				panic("not enough tokens")
			}
			if t.Str == "=" {
				fmt.Println("op=")
				if lastExpr.Type() == ExpressionNop {
					lastExpr = &ColumnReferenceExpression{
						index: currentColumn,
					}
				}
				currentExpr = &OpEquals{
					lhs: lastExpr,
					rhs: nil,
				}
			} else if t.Str == "$" {
				if len(tokens) < i+2 {
					panic("not enough tokens")
				}
				nextToken := tokens[i+1]
				if nextToken.Typ != TokenTypeString {
					panic("invalid token type after $, must be string")
				}
				index, err := strconv.ParseInt(nextToken.Str, 10, 32)
				if err != nil {
					panic("failed to parse column reference to int")
				}
				currentExpr = &ColumnReferenceExpression{
					index: int(index),
				}
				i++
			} else {
				panic("unknown operator")
			}
		} else if t.Typ == TokenTypeString {
			currentExpr = parseLiteral(t.Str)
		} else if t.Typ == TokenTypeNewLine {
			currentExprs = append(currentExprs, lastExpr)
			exprs = append(exprs, currentExprs)
			currentExprs = []Expression{}
			lastExpr = &Nop{}
			currentColumn = 0
			currentExpr = nil
		}
		if currentExpr != nil {
			if b, ok := lastExpr.(BinaryExpr); ok {
				b.SetRHS(currentExpr)
			} else {
				lastExpr = currentExpr
			}
		}
	}
	currentExprs = append(currentExprs, lastExpr)
	exprs = append(exprs, currentExprs)
	return exprs
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
