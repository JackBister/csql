package csql

import (
	"encoding/csv"
	"fmt"
	"io"
	"slices"
	"unicode/utf8"
)

type GroupOperations struct {
	groupExpr       *GroupingExpr
	projectionExprs []*AggregatingExpr
}

func Execute(operations [][]Expression, reader io.Reader, options Options) ([][]string, error) {
	fmt.Println(operations)
	csvReader := csv.NewReader(reader)
	sep, _ := utf8.DecodeRuneInString(options.Separator)
	if sep == utf8.RuneError {
		panic("invalid separator")
	}
	csvReader.Comma = sep
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	if options.Skip != 0 {
		records = records[options.Skip:]
	}

	resultSet := make([][]Value, len(records))
	for i, r := range records {
		resultSet[i] = make([]Value, len(r))
		for j, v := range r {
			literalExpr := parseLiteral(v)
			resultSet[i][j] = literalExpr.value
		}
	}

	for _, ops := range operations {
		nextResultSet := [][]Value{}
		groupOperations := GroupOperations{
			projectionExprs: make([]*AggregatingExpr, 0),
		}
		groupingValues := map[string][]Value{}
		groupedResults := map[string][]Value{}
		orderOperations := make([]*OrderingExpr, 0)
		limitOperations := make([]*LimitExpr, 0)

		for _, op := range ops {
			if op.Type() == ExpressionGrouping {
				if groupOperations.groupExpr != nil {
					panic("cannot have more than one grouping expr in a line")
				}
				groupOperations.groupExpr = op.(*GroupingExpr)
			} else if op.Type() == ExpressionAggregating {
				fnc := op.(*AggregatingExpr)
				groupOperations.projectionExprs = append(groupOperations.projectionExprs, fnc)
			} else if op.Type() == ExpressionOrdering {
				orderOperations = append(orderOperations, op.(*OrderingExpr))
				if len(orderOperations) > MaxOrderingExprCount {
					return nil, fmt.Errorf("too many ordering expressions, maximum is %d", MaxOrderingExprCount)
				}
			} else if op.Type() == ExpressionLimit {
				limitOperations = append(limitOperations, op.(*LimitExpr))
				if len(limitOperations) > 1 {
					return nil, fmt.Errorf("cannot have more than one limit expression in a line")
				}
			}
		}

		for _, record := range resultSet {
			projection := []Value{}
			if groupOperations.groupExpr != nil || len(groupOperations.projectionExprs) > 0 {
				groupString := ""
				groupResults := []Value{}

				if groupOperations.groupExpr != nil {
					res, err := groupOperations.groupExpr.Execute(0, record)
					if err != nil {
						panic(err)
					}
					if res.value != nil {
						groupString = res.value.String()
						asList := res.value.value.([]Value)
						groupingValues[groupString] = asList
						groupResults = append(groupResults, asList...)
					}
				}

				for i, op := range groupOperations.projectionExprs {
					res, err := op.argument.Execute(i, record)
					if err != nil {
						panic(err)
					}
					if res.value != nil {
						groupResults = append(groupResults, *res.value)
					}
				}

				if existing, ok := groupedResults[groupString]; ok {
					startIndex := 0
					if l, ok := groupingValues[groupString]; ok {
						startIndex = len(l)
					}
					for i, v := range existing[startIndex:] {
						next := groupResults[i+startIndex]
						aggr := groupOperations.projectionExprs[i]
						aggrFn, ok := aggregationFuncMap[aggr.aggregationName]
						if !ok {
							panic("aggregation function '" + aggr.aggregationName + "' not found")
						}
						va, err := v.Convert(aggrFn.valueType)
						if err != nil {
							return nil, err
						}
						vb, err := next.Convert(aggrFn.valueType)
						if err != nil {
							return nil, err
						}
						vr, err := aggrFn.fn(*va, *vb)
						if err != nil {
							return nil, err
						}
						existing[i+startIndex] = *vr
					}
				} else {
					groupedResults[groupString] = groupResults
				}
			} else if len(orderOperations) > 0 || len(limitOperations) > 0 {
				nextResultSet = append(nextResultSet, record)
			} else {
				excluded := false
				for i, op := range ops {
					res, err := op.Execute(i, record)
					if err != nil {
						panic(err)
					}
					if res.value != nil {
						if res.value.typ == ValueTypeBool {
							if !res.value.value.(bool) {
								excluded = true
								break
							}
						} else {
							projection = append(projection, *res.value)
						}
					}
				}
				if !excluded {
					if len(projection) > 0 {
						nextResultSet = append(nextResultSet, projection)
					} else {
						nextResultSet = append(nextResultSet, record)
					}
				}
			}
		}
		if len(groupedResults) > 0 {
			for _, gr := range groupedResults {
				nextResultSet = append(nextResultSet, gr)
			}
		}
		if len(orderOperations) > 0 {
			slices.SortFunc(nextResultSet, func(iv, jv []Value) int {
				var sortValueI int64 = 0
				var sortValueJ int64 = 0
				for i, op := range orderOperations {
					resI, err := op.argument.Execute(i, iv)
					if err != nil {
						panic(err)
					}
					resJ, err := op.argument.Execute(i, jv)
					if err != nil {
						panic(err)
					}
					if resI.value == nil || resJ.value == nil {
						continue
					}

					eq := OpEquals{
						lhs: &LiteralExpression{
							value: *resI.value,
						},
						rhs: &LiteralExpression{
							value: *resJ.value,
						},
					}
					eqRes, err := eq.Execute(0, nil)
					if err != nil {
						panic(err)
					}
					if eqRes.value.value.(bool) {
						continue
					}

					lt := OpLt{
						lhs: &LiteralExpression{
							value: *resI.value,
						},
						rhs: &LiteralExpression{
							value: *resJ.value,
						},
					}

					ltRes, err := lt.Execute(0, nil)
					if err != nil {
						panic(err)
					}

					if op.direction == OrderDirectionAsc {
						if !ltRes.value.value.(bool) {
							sortValueI += powInt64(10, MaxOrderingExprCount-int64(i)+1)
						} else {
							sortValueJ += powInt64(10, MaxOrderingExprCount-int64(i)+1)
						}
					} else {
						if ltRes.value.value.(bool) {
							sortValueI += powInt64(10, MaxOrderingExprCount-int64(i)+1)
						} else {
							sortValueJ += powInt64(10, MaxOrderingExprCount-int64(i)+1)
						}
					}
				}
				return int(sortValueI - sortValueJ)
			})
		}
		if len(limitOperations) > 0 {
			limit := limitOperations[0].limit
			if limit < 0 {
				return nil, fmt.Errorf("limit cannot be negative")
			}
			nextResultSet = nextResultSet[:limit]
		}
		resultSet = nextResultSet
	}

	valueTypes := []ValueType{}
	resultSetStrings := make([][]string, len(resultSet))
	for i, r := range resultSet {
		resultSetStrings[i] = make([]string, len(r))
		for j, v := range r {
			if j >= len(valueTypes) {
				valueTypes = append(valueTypes, v.typ)
			} else {
				if valueTypes[j] != v.typ {
					valueTypes[j] = ValueTypeUnknown
				}
			}
			resultSetStrings[i][j] = v.String()
		}
	}

	if options.PrintTypes {
		fmt.Println(valueTypes)
	}

	return resultSetStrings, nil
}

func powInt64(base, exp int64) int64 {
	if exp < 0 {
		return 0
	}
	result := int64(1)
	for i := int64(0); i < exp; i++ {
		result *= base
	}
	return result
}
