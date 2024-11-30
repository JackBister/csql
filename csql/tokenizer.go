package csql

import "strings"

type TokenType int

const (
	TokenTypeComma TokenType = iota
	TokenTypeString
	TokenTypeOperator
	TokenTypeNewLine
)

type Token struct {
	Typ TokenType
	Str string
}

var operators = "$!=><"

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
			res = append(res, Token{
				Typ: TokenTypeString,
				Str: str.String(),
			})
			str.Reset()
			res = append(res, Token{Typ: TokenTypeNewLine, Str: ""})
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
