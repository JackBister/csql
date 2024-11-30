package csql

import (
	"fmt"
	"time"
)

type OpEquals struct {
	lhs Expression
	rhs Expression
}

func evaluateOperands(l, r Expression, i int, record []Value) (*OperationResult, *Value, *Value, error) {
	lhs, err := l.Execute(i, record)
	if err != nil {
		return nil, nil, nil, err
	}
	if lhs == nil {
		// TODO: Should nil == nil?
		return &OperationResult{
			value: &Value{
				typ:   ValueTypeBool,
				value: false,
			},
		}, nil, nil, nil
	}
	rhs, err := r.Execute(i, record)
	if err != nil {
		return nil, nil, nil, err
	}
	if rhs == nil {
		// TODO: Should nil == nil?
		return &OperationResult{
			value: &Value{
				typ:   ValueTypeBool,
				value: false,
			},
		}, nil, nil, nil
	}
	lhsV := lhs.value
	rhsV := rhs.value
	return nil, lhsV, rhsV, nil
}

func (o *OpEquals) Execute(i int, record []Value) (*OperationResult, error) {
	res, lhsV, rhsV, err := evaluateOperands(o.lhs, o.rhs, i, record)
	if err != nil {
		return nil, err
	}
	if res != nil {
		return res, nil
	}
	if lhsV.typ != rhsV.typ {
		rhsV, err = rhsV.Convert(lhsV.typ)
		if err != nil {
			return &OperationResult{
				value: &Value{
					typ:   ValueTypeBool,
					value: false,
				},
			}, nil
		}
	}
	return &OperationResult{
		value: &Value{
			typ:   ValueTypeBool,
			value: lhsV.value == rhsV.value,
		},
	}, nil
}

func (o *OpEquals) FillNils(e Expression) {
	if o.lhs != nil {
		o.lhs.FillNils(e)
	} else {
		o.lhs = e
	}
	if o.rhs != nil {
		o.rhs.FillNils(e)
	} else {
		o.rhs = e
	}
}

func (o *OpEquals) GetLHS() Expression {
	return o.lhs
}

func (o *OpEquals) GetRHS() Expression {
	return o.lhs
}

func (o *OpEquals) SetLHS(e Expression) {
	o.lhs = e
}

func (o *OpEquals) SetRHS(e Expression) {
	o.rhs = e
}

func (o *OpEquals) Type() ExpressionType {
	return ExpressionOperator
}

func (o *OpEquals) String() string {
	return fmt.Sprintf("(OpEquals: LHS=%v, RHS=%v)", o.lhs, o.rhs)
}

type OpLt struct {
	lhs Expression
	rhs Expression
}

func (o *OpLt) Execute(i int, record []Value) (*OperationResult, error) {
	res, lhsV, rhsV, err := evaluateOperands(o.lhs, o.rhs, i, record)
	if err != nil {
		return nil, err
	}
	if res != nil {
		return res, nil
	}

	if lhsV.typ != rhsV.typ {
		rhsV, err = rhsV.Convert(lhsV.typ)
		if err != nil {
			return &OperationResult{
				value: &Value{
					typ:   ValueTypeBool,
					value: false,
				},
			}, nil
		}
	}
	if lhsV.typ == ValueTypeInt {
		return &OperationResult{
			value: &Value{
				typ:   ValueTypeBool,
				value: lhsV.value.(int64) < rhsV.value.(int64),
			},
		}, nil
	} else if lhsV.typ == ValueTypeDouble {
		return &OperationResult{
			value: &Value{
				typ:   ValueTypeBool,
				value: lhsV.value.(float64) < rhsV.value.(float64),
			},
		}, nil
	} else if lhsV.typ == ValueTypeDate {
		return &OperationResult{
			value: &Value{
				typ:   ValueTypeBool,
				value: lhsV.value.(time.Time).Before(rhsV.value.(time.Time)),
			},
		}, nil
	}
	return nil, fmt.Errorf("operator < is not valid for type %v", lhsV.typ)
}

func (o *OpLt) FillNils(e Expression) {
	if o.lhs != nil {
		o.lhs.FillNils(e)
	} else {
		o.lhs = e
	}
	if o.rhs != nil {
		o.rhs.FillNils(e)
	} else {
		o.rhs = e
	}
}

func (o *OpLt) GetLHS() Expression {
	return o.lhs
}

func (o *OpLt) GetRHS() Expression {
	return o.lhs
}

func (o *OpLt) SetLHS(e Expression) {
	o.lhs = e
}

func (o *OpLt) SetRHS(e Expression) {
	o.rhs = e
}

func (o *OpLt) Type() ExpressionType {
	return ExpressionOperator
}

func (o *OpLt) String() string {
	return fmt.Sprintf("(OpLt: LHS=%v, RHS=%v)", o.lhs, o.rhs)
}

type OpGt struct {
	lhs Expression
	rhs Expression
}

func (o *OpGt) Execute(i int, record []Value) (*OperationResult, error) {
	res, lhsV, rhsV, err := evaluateOperands(o.lhs, o.rhs, i, record)
	if err != nil {
		return nil, err
	}
	if res != nil {
		return res, nil
	}

	if lhsV.typ != rhsV.typ {
		rhsV, err = rhsV.Convert(lhsV.typ)
		if err != nil {
			return &OperationResult{
				value: &Value{
					typ:   ValueTypeBool,
					value: false,
				},
			}, nil
		}
	}
	if lhsV.typ == ValueTypeInt {
		return &OperationResult{
			value: &Value{
				typ:   ValueTypeBool,
				value: lhsV.value.(int64) > rhsV.value.(int64),
			},
		}, nil
	} else if lhsV.typ == ValueTypeDouble {
		return &OperationResult{
			value: &Value{
				typ:   ValueTypeBool,
				value: lhsV.value.(float64) > rhsV.value.(float64),
			},
		}, nil
	} else if lhsV.typ == ValueTypeDate {
		return &OperationResult{
			value: &Value{
				typ:   ValueTypeBool,
				value: lhsV.value.(time.Time).After(rhsV.value.(time.Time)),
			},
		}, nil
	}
	return nil, fmt.Errorf("operator > is not valid for type %v", lhsV.typ)
}

func (o *OpGt) FillNils(e Expression) {
	if o.lhs != nil {
		o.lhs.FillNils(e)
	} else {
		o.lhs = e
	}
	if o.rhs != nil {
		o.rhs.FillNils(e)
	} else {
		o.rhs = e
	}
}

func (o *OpGt) GetLHS() Expression {
	return o.lhs
}

func (o *OpGt) GetRHS() Expression {
	return o.lhs
}

func (o *OpGt) SetLHS(e Expression) {
	o.lhs = e
}

func (o *OpGt) SetRHS(e Expression) {
	o.rhs = e
}

func (o *OpGt) Type() ExpressionType {
	return ExpressionOperator
}

func (o *OpGt) String() string {
	return fmt.Sprintf("(OpGt: LHS=%v, RHS=%v)", o.lhs, o.rhs)
}

type OpNeg struct {
	inner Expression
}

func (o *OpNeg) Execute(i int, record []Value) (*OperationResult, error) {
	res, err := o.inner.Execute(i, record)
	if err != nil {
		return nil, err
	}
	if res.value.typ != ValueTypeBool {
		return nil, fmt.Errorf("cannot negate non-boolean value: %v", res.value)
	}
	return &OperationResult{
		value: &Value{
			typ:   ValueTypeBool,
			value: !res.value.value.(bool),
		},
	}, nil
}

func (o *OpNeg) FillNils(e Expression) {
	if o.inner != nil {
		o.inner.FillNils(e)
	} else {
		o.inner = e
	}
}

func (o *OpNeg) Type() ExpressionType {
	return ExpressionOperator
}

func (o *OpNeg) String() string {
	return fmt.Sprintf("(OpNeg: Inner=%v)", o.inner)
}
