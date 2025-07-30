package csql

import (
	"fmt"
	"strings"
)

type FunctionType int

const (
	FunctionTypeDefault FunctionType = iota
)

type Function struct {
	argumentTypes []ValueType
	fn            func(args ExpressionList, i int, record []Value) (*Value, error)
}

var funcMap = map[string]Function{
	"has": {
		argumentTypes: []ValueType{ValueTypeUnknown, ValueTypeString},
		fn: func(args ExpressionList, i int, record []Value) (*Value, error) {
			col, err := args.exprs[0].Execute(i, record)
			if err != nil {
				return nil, err
			}
			arg0, err := args.exprs[1].Execute(i, record)
			if err != nil {
				return nil, err
			}
			colStr, err := col.value.Convert(ValueTypeString)
			if err != nil {
				return nil, err
			}
			arg0Str, err := arg0.value.Convert(ValueTypeString)
			if err != nil {
				return nil, err
			}
			res := strings.Contains(colStr.String(), arg0Str.String())
			return &Value{
				typ:   ValueTypeBool,
				value: res,
			}, nil
		},
	},
}

type AggregationFunction struct {
	fn        func(a, b Value) (*Value, error)
	valueType ValueType
}

var aggregationFuncMap = map[string]AggregationFunction{
	"sum": {
		valueType: ValueTypeDouble,
		fn: func(a, b Value) (*Value, error) {
			if a.typ != ValueTypeDouble || b.typ != ValueTypeDouble {
				return nil, fmt.Errorf("parameters must be doubles")
			}
			res := a.value.(float64) + b.value.(float64)
			return &Value{
				typ:   ValueTypeDouble,
				value: res,
			}, nil
		},
	},
}

type ExpressionList struct {
	exprs []Expression
}

func (el *ExpressionList) Execute(i int, record []Value) (*OperationResult, error) {
	return nil, fmt.Errorf("expression lists cannot be executed")
}

func (el *ExpressionList) FillNils(e Expression) {
	for i, expr := range el.exprs {
		if expr.Type() == ExpressionNop {
			el.exprs[i] = e
		} else {
			expr.FillNils(e)
		}
	}
}

func (el *ExpressionList) String() string {
	return fmt.Sprintf("(ExprList: {%v})", el.exprs)
}

func (el *ExpressionList) Type() ExpressionType {
	return ExpressionExprList
}

type Funcall struct {
	funcName  string
	arguments ExpressionList
}

func (f *Funcall) Execute(i int, record []Value) (*OperationResult, error) {
	fn, ok := funcMap[f.funcName]
	if !ok {
		return nil, fmt.Errorf("function '%v' not found", fn)
	}
	res, err := fn.fn(f.arguments, i, record)
	if err != nil {
		return nil, err
	}
	return &OperationResult{
		value: res,
	}, nil
}

func (f *Funcall) FillNils(e Expression) {
	if len(f.arguments.exprs) > 0 {
		if f.arguments.exprs[0].Type() != ExpressionColumnReference {
			f.arguments.exprs = append([]Expression{e}, f.arguments.exprs...)
		}
		f.arguments.FillNils(e)
	} else {
		f.arguments = ExpressionList{
			exprs: []Expression{e},
		}
	}
}

func (f *Funcall) String() string {
	return fmt.Sprintf("(Funcall: Name=%v Args={%v})", f.funcName, f.arguments)
}

func (f *Funcall) Type() ExpressionType {
	return ExpressionFuncall
}

type GroupingExpr struct {
	arguments ExpressionList
}

func (f *GroupingExpr) Execute(i int, record []Value) (*OperationResult, error) {
	ret := []Value{}

	for _, a := range f.arguments.exprs {
		res, err := a.Execute(i, record)
		if err != nil {
			return nil, err
		}
		if res != nil {
			ret = append(ret, *res.value)
		}
	}

	return &OperationResult{
		value: &Value{
			typ:   ValueTypeList,
			value: ret,
		},
	}, nil
}

func (f *GroupingExpr) FillNils(e Expression) {
	if len(f.arguments.exprs) > 0 {
		f.arguments.FillNils(e)
	} else {
		f.arguments = ExpressionList{
			exprs: []Expression{e},
		}
	}
}

func (f *GroupingExpr) String() string {
	return fmt.Sprintf("(Grouping: Args={%v})", f.arguments)
}

func (f *GroupingExpr) Type() ExpressionType {
	return ExpressionGrouping
}

type AggregatingExpr struct {
	aggregationName string
	argument        Expression
}

func (f *AggregatingExpr) Execute(i int, record []Value) (*OperationResult, error) {
	return f.argument.Execute(i, record)
}

func (f *AggregatingExpr) FillNils(e Expression) {
	if f.argument.Type() == ExpressionNop {
		f.argument = e
	} else if f.argument != nil {
		f.argument.FillNils(e)
	} else {
		f.argument = e
	}
}

func (f *AggregatingExpr) String() string {
	return fmt.Sprintf("(Aggregating: Name=%v Arg={%v})", f.aggregationName, f.argument)
}

func (f *AggregatingExpr) Type() ExpressionType {
	return ExpressionAggregating
}

const MaxOrderingExprCount = 10

type OrderDirection int

const (
	OrderDirectionAsc OrderDirection = iota
	OrderDirectionDesc
)

type OrderingExpr struct {
	argument  Expression
	direction OrderDirection
}

func (f *OrderingExpr) Execute(i int, record []Value) (*OperationResult, error) {
	res, err := f.argument.Execute(i, record)
	if err != nil {
		return nil, err
	}
	if res == nil || res.value == nil {
		return nil, fmt.Errorf("ordering expression returned nil value")
	}
	return res, nil
}

func (f *OrderingExpr) FillNils(e Expression) {
	if f.argument.Type() == ExpressionNop {
		f.argument = e
	} else if f.argument != nil {
		f.argument.FillNils(e)
	} else {
		f.argument = e
	}
}

func (f *OrderingExpr) String() string {
	return fmt.Sprintf("(Ordering: Arg={%v} Direction=%v)", f.argument, f.direction)
}

func (f *OrderingExpr) Type() ExpressionType {
	return ExpressionOrdering
}

type LimitExpr struct {
	limit int64
}

func (f *LimitExpr) Execute(i int, record []Value) (*OperationResult, error) {
	return nil, fmt.Errorf("limit expressions cannot be executed")
}

func (f *LimitExpr) FillNils(e Expression) {
	// Limit expressions do not have arguments to fill
}

func (f *LimitExpr) String() string {
	return fmt.Sprintf("(Limit: %d)", f.limit)
}

func (f *LimitExpr) Type() ExpressionType {
	return ExpressionLimit
}
