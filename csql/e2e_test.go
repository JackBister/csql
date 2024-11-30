package csql_test

import (
	"strings"
	"testing"

	"github.com/jackbister/csql/csql"
)

var testCsv = `1,a,b,c
2,d,e,f
1,a,dafa,ssz
4,a,a,y`

func TestSimpleFilter(t *testing.T) {
	query := "=1"

	tokens := csql.Tokenize(query)
	exprs, err := csql.ParseQuery(tokens)
	if err != nil {
		t.FailNow()
	}
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

func TestSimpleFilterColumnReference(t *testing.T) {
	query := "$1=a"

	tokens := csql.Tokenize(query)
	exprs, err := csql.ParseQuery(tokens)
	if err != nil {
		t.FailNow()
	}
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
	if res[2][2] != "a" {
		t.FailNow()
	}
}

func TestEqualsByReference(t *testing.T) {
	query := "$1=$2"

	tokens := csql.Tokenize(query)
	exprs, err := csql.ParseQuery(tokens)
	if err != nil {
		t.FailNow()
	}
	res, err := csql.Execute(exprs, strings.NewReader(testCsv), csql.NewOptions())
	if err != nil {
		t.FailNow()
	}
	if len(res) != 1 {
		t.FailNow()
	}
	if res[0][3] != "y" {
		t.FailNow()
	}
}

func TestMultipleSteps(t *testing.T) {
	query := ",=a\n,,,=y"

	tokens := csql.Tokenize(query)
	exprs, err := csql.ParseQuery(tokens)
	if err != nil {
		t.FailNow()
	}
	res, err := csql.Execute(exprs, strings.NewReader(testCsv), csql.NewOptions())
	if err != nil {
		t.FailNow()
	}
	if len(res) != 1 {
		t.FailNow()
	}
	if res[0][2] != "a" {
		t.FailNow()
	}
}

func TestProjection(t *testing.T) {
	query := "$0"

	tokens := csql.Tokenize(query)
	exprs, err := csql.ParseQuery(tokens)
	if err != nil {
		t.FailNow()
	}
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
	if res[3][0] != "4" {
		t.FailNow()
	}
}

func TestTrueLiteral(t *testing.T) {
	query := "true"

	tokens := csql.Tokenize(query)
	exprs, err := csql.ParseQuery(tokens)
	if err != nil {
		t.FailNow()
	}
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
	exprs, err := csql.ParseQuery(tokens)
	if err != nil {
		t.FailNow()
	}
	res, err := csql.Execute(exprs, strings.NewReader(testCsv), csql.NewOptions())
	if err != nil {
		t.FailNow()
	}
	if len(res) != 0 {
		t.FailNow()
	}
}

func TestGt(t *testing.T) {
	query := ">1"

	tokens := csql.Tokenize(query)
	exprs, err := csql.ParseQuery(tokens)
	if err != nil {
		t.FailNow()
	}
	res, err := csql.Execute(exprs, strings.NewReader(testCsv), csql.NewOptions())
	if err != nil {
		t.FailNow()
	}
	if len(res) != 2 {
		t.FailNow()
	}
}

func TestLt(t *testing.T) {
	query := "<2"

	tokens := csql.Tokenize(query)
	exprs, err := csql.ParseQuery(tokens)
	if err != nil {
		t.FailNow()
	}
	res, err := csql.Execute(exprs, strings.NewReader(testCsv), csql.NewOptions())
	if err != nil {
		t.FailNow()
	}
	if len(res) != 2 {
		t.FailNow()
	}
}

func TestNeg(t *testing.T) {
	query := "!>1"

	tokens := csql.Tokenize(query)
	exprs, err := csql.ParseQuery(tokens)
	if err != nil {
		t.FailNow()
	}
	res, err := csql.Execute(exprs, strings.NewReader(testCsv), csql.NewOptions())
	if err != nil {
		t.FailNow()
	}
	if len(res) != 2 {
		t.FailNow()
	}
}
