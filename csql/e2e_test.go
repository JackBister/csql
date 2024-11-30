package csql_test

import (
	"strings"
	"testing"

	"github.com/jackbister/csql/csql"
)

var testCsv = `1,a,b,c
2,d,e,f
1,a,dafa,ssz
a,a,x,y`

func TestSimpleFilter(t *testing.T) {
	query := "=1"

	tokens := csql.Tokenize(query)
	exprs := csql.Parse(tokens)
	res, err := csql.Execute(exprs, strings.NewReader(testCsv), csql.NewOptions())
	if err != nil {
		t.FailNow()
	}
	if len(res) != 2 {
		t.FailNow()
	}
	if res[0][2] != "b" {
		t.FailNow()
	}
	if res[1][2] != "dafa" {
		t.FailNow()
	}
}

func TestSimlpeFilterColumnReference(t *testing.T) {
	query := "$1=a"

	tokens := csql.Tokenize(query)
	exprs := csql.Parse(tokens)
	res, err := csql.Execute(exprs, strings.NewReader(testCsv), csql.NewOptions())
	if err != nil {
		t.FailNow()
	}
	if len(res) != 3 {
		t.FailNow()
	}
	if res[0][2] != "b" {
		t.FailNow()
	}
	if res[1][2] != "dafa" {
		t.FailNow()
	}
	if res[2][2] != "x" {
		t.FailNow()
	}
}

func TestEqualsByReference(t *testing.T) {
	query := "$0=$1"

	tokens := csql.Tokenize(query)
	exprs := csql.Parse(tokens)
	res, err := csql.Execute(exprs, strings.NewReader(testCsv), csql.NewOptions())
	if err != nil {
		t.FailNow()
	}
	if len(res) != 1 {
		t.FailNow()
	}
	if res[0][2] != "x" {
		t.FailNow()
	}
}

func TestMultipleSteps(t *testing.T) {
	query := ",=a\n,,=x"

	tokens := csql.Tokenize(query)
	exprs := csql.Parse(tokens)
	res, err := csql.Execute(exprs, strings.NewReader(testCsv), csql.NewOptions())
	if err != nil {
		t.FailNow()
	}
	if len(res) != 1 {
		t.FailNow()
	}
	if res[0][2] != "x" {
		t.FailNow()
	}
}

func TestProjection(t *testing.T) {
	query := "$0"

	tokens := csql.Tokenize(query)
	exprs := csql.Parse(tokens)
	res, err := csql.Execute(exprs, strings.NewReader(testCsv), csql.NewOptions())
	if err != nil {
		t.FailNow()
	}
	if len(res) != 4 {
		t.FailNow()
	}
	if res[0][0] != "1" {
		t.FailNow()
	}
	if res[1][0] != "2" {
		t.FailNow()
	}
	if res[2][0] != "1" {
		t.FailNow()
	}
	if res[3][0] != "a" {
		t.FailNow()
	}
}

func TestTrueLiteral(t *testing.T) {
	query := "true"

	tokens := csql.Tokenize(query)
	exprs := csql.Parse(tokens)
	res, err := csql.Execute(exprs, strings.NewReader(testCsv), csql.NewOptions())
	if err != nil {
		t.FailNow()
	}
	if len(res) != 4 {
		t.FailNow()
	}
}

func TestFalseLiteral(t *testing.T) {
	query := "false"

	tokens := csql.Tokenize(query)
	exprs := csql.Parse(tokens)
	res, err := csql.Execute(exprs, strings.NewReader(testCsv), csql.NewOptions())
	if err != nil {
		t.FailNow()
	}
	if len(res) != 0 {
		t.FailNow()
	}
}
