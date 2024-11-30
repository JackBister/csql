package csql

import (
	"encoding/csv"
	"fmt"
	"io"
	"unicode/utf8"
)

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
		for _, record := range resultSet {
			excluded := false
			projection := []Value{}
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
