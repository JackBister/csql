// CSQL - A command-line tool for CSV querying
// Copyright (C) 2025  Jack Bister
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

//go:generate stringer -type=TokenType
package csql

import "strings"

type TokenType int

const (
	TokenTypeComma TokenType = iota
	TokenTypeString
	TokenTypeOperator
	TokenTypeNewLine
	TokenTypeLParen
	TokenTypeRParen
)

type Token struct {
	Typ TokenType
	Str string
}

var operators = "$!=><+-*/"

func Tokenize(query string) []Token {
	res := []Token{}

	str := strings.Builder{}
	for _, c := range query {
		if c == ',' {
			if str.Len() > 0 {
				res = append(res, Token{
					Typ: TokenTypeString,
					Str: str.String(),
				})
				str.Reset()
			}
			res = append(res, Token{Typ: TokenTypeComma, Str: ""})
		} else if strings.ContainsRune(operators, c) {
			if str.Len() > 0 {
				res = append(res, Token{
					Typ: TokenTypeString,
					Str: str.String(),
				})
				str.Reset()
			}
			res = append(res, Token{Typ: TokenTypeOperator, Str: string(c)})
		} else if c == '\n' {
			if str.Len() > 0 {
				res = append(res, Token{
					Typ: TokenTypeString,
					Str: str.String(),
				})
			}
			str.Reset()
			res = append(res, Token{Typ: TokenTypeNewLine, Str: ""})
		} else if c == '(' {
			if str.Len() > 0 {
				res = append(res, Token{
					Typ: TokenTypeString,
					Str: str.String(),
				})
				str.Reset()
			}
			res = append(res, Token{Typ: TokenTypeLParen})
		} else if c == ')' {
			if str.Len() > 0 {
				res = append(res, Token{
					Typ: TokenTypeString,
					Str: str.String(),
				})
				str.Reset()
			}
			res = append(res, Token{Typ: TokenTypeRParen})
		} else {
			str.WriteRune(c)
		}
	}
	if str.Len() > 0 {
		res = append(res, Token{
			Typ: TokenTypeString,
			Str: str.String(),
		})
	}
	return res
}
