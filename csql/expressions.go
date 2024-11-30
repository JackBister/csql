package csql

import (
	"fmt"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
)

type ExpressionType int

const (
	ExpressionNop = iota
	ExpressionEquals
	ExpressionLiteral
	ExpressionColumnReference
)

type ValueType int

const (
	ValueTypeUnknown = iota
	ValueTypeString
	ValueTypeBool
	ValueTypeInt
	ValueTypeDouble
	ValueTypeDate
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
	}
	panic("unhandled value type " + fmt.Sprint(v.typ))
}

type OperationResult struct {
	value *Value
}

type Expression interface {
	Execute(i int, record []Value) (*OperationResult, error)
	Type() ExpressionType
}

type BinaryExpr interface {
	SetLHS(e Expression)
	SetRHS(e Expression)
}
