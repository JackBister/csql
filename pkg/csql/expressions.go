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

//go:generate stringer -type=ExpressionType
//go:generate stringer -type=ValueType
package csql

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

type ExpressionType int

const (
	ExpressionNop ExpressionType = iota
	ExpressionOperator
	ExpressionLiteral
	ExpressionColumnReference
	ExpressionExprList
	ExpressionFuncall
	ExpressionGrouping
	ExpressionAggregating
	ExpressionOrdering
	ExpressionLimit
)

type ValueType int

const (
	ValueTypeUnknown ValueType = iota
	ValueTypeString
	ValueTypeBool
	ValueTypeInt
	ValueTypeDouble
	ValueTypeDate
	ValueTypeList
)

type Value struct {
	typ   ValueType
	value any
}

func (v *Value) Convert(targetType ValueType) (*Value, error) {
	if targetType == v.typ {
		return v, nil
	}
	if targetType == ValueTypeString {
		return &Value{
			typ:   ValueTypeString,
			value: v.String(),
		}, nil
	} else if targetType == ValueTypeBool {
		if v.typ == ValueTypeString {
			str := v.value.(string)
			if str == "true" {
				return &Value{
					typ:   ValueTypeBool,
					value: true,
				}, nil
			} else if str == "false" {
				return &Value{
					typ:   ValueTypeBool,
					value: false,
				}, nil
			}
		} else if v.typ == ValueTypeInt {
			i := v.value.(int64)
			return &Value{
				typ:   ValueTypeBool,
				value: i != 0,
			}, nil
		}
	} else if targetType == ValueTypeInt {
		if v.typ == ValueTypeString {
			i, err := strconv.ParseInt(v.value.(string), 10, 64)
			if err == nil {
				return &Value{
					typ:   ValueTypeInt,
					value: i,
				}, nil
			}
		} else if v.typ == ValueTypeBool {
			i := 0
			if v.value.(bool) {
				i = 1
			}
			return &Value{
				typ:   ValueTypeInt,
				value: i,
			}, nil
		} else if v.typ == ValueTypeDouble {
			i := int64(v.value.(float64))
			return &Value{
				typ:   ValueTypeInt,
				value: i,
			}, nil
		}
	} else if targetType == ValueTypeDouble {
		if v.typ == ValueTypeString {
			d, err := strconv.ParseFloat(v.value.(string), 64)
			if err == nil {
				return &Value{
					typ:   ValueTypeDouble,
					value: d,
				}, nil
			}
		} else if v.typ == ValueTypeInt {
			d := float64(v.value.(int64))
			return &Value{
				typ:   ValueTypeDouble,
				value: d,
			}, nil
		}
	} else if targetType == ValueTypeDate {
		if v.typ == ValueTypeString {
			t, err := dateparse.ParseAny(v.value.(string))
			if err == nil {
				return &Value{
					typ:   ValueTypeDate,
					value: t,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("unhandled conversion from type: %v to targetType: %v", v.typ, targetType)
}

func (v *Value) String() string {
	if v.typ == ValueTypeBool {
		if v.value.(bool) {
			return "true"
		}
		return "false"
	} else if v.typ == ValueTypeString {
		return v.value.(string)
	} else if v.typ == ValueTypeInt {
		return strconv.FormatInt(v.value.(int64), 10)
	} else if v.typ == ValueTypeDouble {
		return strconv.FormatFloat(v.value.(float64), 'f', -1, 64)
	} else if v.typ == ValueTypeDate {
		return v.value.(time.Time).String()
	} else if v.typ == ValueTypeList {
		res := strings.Builder{}
		res.WriteRune('[')
		asList := v.value.([]Value)
		for i, v := range asList {
			res.WriteString(v.String())
			if i != len(asList)-1 {
				res.WriteString(", ")
			}
		}
		res.WriteRune(']')
		return res.String()
	}
	panic("unhandled value type " + fmt.Sprint(v.typ))
}

type OperationResult struct {
	value *Value
}

type Expression interface {
	Execute(i int, record []Value) (*OperationResult, error)
	FillNils(e Expression)
	Type() ExpressionType
}

type BinaryExpr interface {
	GetLHS() Expression
	GetRHS() Expression
	SetLHS(e Expression)
	SetRHS(e Expression)
}
