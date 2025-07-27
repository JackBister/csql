package csql_test

import (
	"fmt"
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

func TestFuncall(t *testing.T) {
	testCsv := `contrary,desktop,1
	continue,printing,2
	established,description,3
	variations,combined,4
	available,repetition,5`
	query := "has(cont)"

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

func TestFuncallColumnRef(t *testing.T) {
	testCsv := `contrary,desktop,1
	continue,printing,2
	established,description,3
	variations,combined,4
	available,repetition,5`
	query := "has($1,des)"

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

func TestAggregationGroupBy(t *testing.T) {
	testCsv := `Peter,Part0,100,50
Peter,Part1,200,60
Peter,Part2,133,220
Peter,Part3,400,10
Peter,Part4,250,30
Peter,Part5,105,40
Charles,Part0,10,50
Charles,Part1,20,60
Charles,Part2,53,220
Charles,Part3,66,10
Charles,Part4,123,30
Charles,Part5,44,40`
	query := "group()"

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
	if res[0][0] != "Peter" {
		t.FailNow()
	}
	if res[1][0] != "Charles" {
		t.FailNow()
	}
}

func TestAggregationSumColumnReference(t *testing.T) {
	testCsv := `Peter,Part0,100,50
Peter,Part1,200,60
Peter,Part2,133,220
Peter,Part3,400,10
Peter,Part4,250,30
Peter,Part5,105,40
Charles,Part0,10,50
Charles,Part1,20,60
Charles,Part2,53,220
Charles,Part3,66,10
Charles,Part4,123,30
Charles,Part5,44,40`
	query := "group(),sum($2)"

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
	if res[0][0] != "Peter" {
		t.FailNow()
	}
	if res[0][1] != "1188" {
		t.FailNow()
	}
	if res[1][0] != "Charles" {
		t.FailNow()
	}
	if res[1][1] != "316" {
		t.FailNow()
	}
}

func TestAggregationSumImplicitColumn(t *testing.T) {
	testCsv := `Peter,Part0,100,50
Peter,Part1,200,60
Peter,Part2,133,220
Peter,Part3,400,10
Peter,Part4,250,30
Peter,Part5,105,40
Charles,Part0,10,50
Charles,Part1,20,60
Charles,Part2,53,220
Charles,Part3,66,10
Charles,Part4,123,30
Charles,Part5,44,40`
	query := "group(),,sum()"

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
	if res[0][0] != "Peter" {
		t.FailNow()
	}
	if res[0][1] != "1188" {
		t.FailNow()
	}
	if res[1][0] != "Charles" {
		t.FailNow()
	}
	if res[1][1] != "316" {
		t.FailNow()
	}
}

func TestAggregationSumGroupByMultiple(t *testing.T) {
	testCsv := `Peter,Part0,100,50
Peter,Part0,200,60
Peter,Part1,133,220
Peter,Part1,400,10
Charles,Part0,10,50
Charles,Part0,20,60
Charles,Part1,53,220
Charles,Part1,66,10`
	query := "group($0,$1),sum($2)"

	tokens := csql.Tokenize(query)
	exprs, err := csql.ParseQuery(tokens)
	if err != nil {
		t.FailNow()
	}
	res, err := csql.Execute(exprs, strings.NewReader(testCsv), csql.NewOptions())
	fmt.Println(res, err)
	if err != nil {
		t.FailNow()
	}
	if len(res) != 4 {
		t.FailNow()
	}
	if res[0][0] != "Peter" {
		t.FailNow()
	}
	if res[0][1] != "Part0" {
		t.FailNow()
	}
	if res[0][2] != "300" {
		t.FailNow()
	}
	if res[1][0] != "Peter" {
		t.FailNow()
	}
	if res[1][1] != "Part1" {
		t.FailNow()
	}
	if res[1][2] != "533" {
		t.FailNow()
	}
	if res[2][0] != "Charles" {
		t.FailNow()
	}
	if res[2][1] != "Part0" {
		t.FailNow()
	}
	if res[2][2] != "30" {
		t.FailNow()
	}
	if res[3][0] != "Charles" {
		t.FailNow()
	}
	if res[3][1] != "Part1" {
		t.FailNow()
	}
	if res[3][2] != "119" {
		t.FailNow()
	}
}

func TestAggregationSumUngrouped(t *testing.T) {
	testCsv := `Peter,Part0,100,50
Peter,Part1,200,60
Peter,Part2,133,220
Peter,Part3,400,10
Peter,Part4,250,30
Peter,Part5,105,40
Charles,Part0,10,50
Charles,Part1,20,60
Charles,Part2,53,220
Charles,Part3,66,10
Charles,Part4,123,30
Charles,Part5,44,40`
	query := "sum($2)"

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
	if res[0][0] != "1504" {
		t.FailNow()
	}
}
