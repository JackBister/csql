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

package csql

import (
	"fmt"
	"strings"
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
	} else if lhsV.typ == ValueTypeString {
		return &OperationResult{
			value: &Value{
				typ:   ValueTypeBool,
				value: strings.Compare(lhsV.value.(string), rhsV.value.(string)) < 0,
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
	} else if lhsV.typ == ValueTypeString {
		return &OperationResult{
			value: &Value{
				typ:   ValueTypeBool,
				value: strings.Compare(lhsV.value.(string), rhsV.value.(string)) > 0,
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

type OpAdd struct {
	lhs Expression
	rhs Expression
}

func (o *OpAdd) Execute(i int, record []Value) (*OperationResult, error) {
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
				typ:   ValueTypeInt,
				value: lhsV.value.(int64) + rhsV.value.(int64),
			},
		}, nil
	} else if lhsV.typ == ValueTypeDouble {
		return &OperationResult{
			value: &Value{
				typ:   ValueTypeDouble,
				value: lhsV.value.(float64) + rhsV.value.(float64),
			},
		}, nil
	}
	return nil, fmt.Errorf("operator + is not valid for type %v", lhsV.typ)
}

func (o *OpAdd) FillNils(e Expression) {
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

func (o *OpAdd) GetLHS() Expression {
	return o.lhs
}

func (o *OpAdd) GetRHS() Expression {
	return o.lhs
}

func (o *OpAdd) SetLHS(e Expression) {
	o.lhs = e
}
func (o *OpAdd) SetRHS(e Expression) {
	o.rhs = e
}

func (o *OpAdd) Type() ExpressionType {
	return ExpressionOperator
}

func (o *OpAdd) String() string {
	return fmt.Sprintf("(OpAdd: LHS=%v, RHS=%v)", o.lhs, o.rhs)
}

type OpSub struct {
	lhs Expression
	rhs Expression
}

func (o *OpSub) Execute(i int, record []Value) (*OperationResult, error) {
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
				typ:   ValueTypeInt,
				value: lhsV.value.(int64) - rhsV.value.(int64),
			},
		}, nil
	} else if lhsV.typ == ValueTypeDouble {
		return &OperationResult{
			value: &Value{
				typ:   ValueTypeDouble,
				value: lhsV.value.(float64) - rhsV.value.(float64),
			},
		}, nil
	}
	return nil, fmt.Errorf("operator - is not valid for type %v", lhsV.typ)
}

func (o *OpSub) FillNils(e Expression) {
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

func (o *OpSub) GetLHS() Expression {
	return o.lhs
}

func (o *OpSub) GetRHS() Expression {
	return o.lhs
}

func (o *OpSub) SetLHS(e Expression) {
	o.lhs = e
}

func (o *OpSub) SetRHS(e Expression) {
	o.rhs = e
}

func (o *OpSub) Type() ExpressionType {
	return ExpressionOperator
}

func (o *OpSub) String() string {
	return fmt.Sprintf("(OpSub: LHS=%v, RHS=%v)", o.lhs, o.rhs)
}

type OpMul struct {
	lhs Expression
	rhs Expression
}

func (o *OpMul) Execute(i int, record []Value) (*OperationResult, error) {
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
				typ:   ValueTypeInt,
				value: lhsV.value.(int64) * rhsV.value.(int64),
			},
		}, nil
	} else if lhsV.typ == ValueTypeDouble {
		return &OperationResult{
			value: &Value{
				typ:   ValueTypeDouble,
				value: lhsV.value.(float64) * rhsV.value.(float64),
			},
		}, nil
	}
	return nil, fmt.Errorf("operator * is not valid for type %v", lhsV.typ)
}

func (o *OpMul) FillNils(e Expression) {
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

func (o *OpMul) GetLHS() Expression {
	return o.lhs
}

func (o *OpMul) GetRHS() Expression {
	return o.lhs
}

func (o *OpMul) SetLHS(e Expression) {
	o.lhs = e
}

func (o *OpMul) SetRHS(e Expression) {
	o.rhs = e
}

func (o *OpMul) Type() ExpressionType {
	return ExpressionOperator
}

func (o *OpMul) String() string {
	return fmt.Sprintf("(OpMul: LHS=%v, RHS=%v)", o.lhs, o.rhs)
}

type OpDiv struct {
	lhs Expression
	rhs Expression
}

func (o *OpDiv) Execute(i int, record []Value) (*OperationResult, error) {
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
				typ:   ValueTypeDouble,
				value: float64(lhsV.value.(int64)) / float64(rhsV.value.(int64)),
			},
		}, nil
	} else if lhsV.typ == ValueTypeDouble {
		return &OperationResult{
			value: &Value{
				typ:   ValueTypeDouble,
				value: lhsV.value.(float64) / rhsV.value.(float64),
			},
		}, nil
	}
	return nil, fmt.Errorf("operator / is not valid for type %v", lhsV.typ)
}

func (o *OpDiv) FillNils(e Expression) {
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

func (o *OpDiv) GetLHS() Expression {
	return o.lhs
}

func (o *OpDiv) GetRHS() Expression {
	return o.lhs
}

func (o *OpDiv) SetLHS(e Expression) {
	o.lhs = e
}

func (o *OpDiv) SetRHS(e Expression) {
	o.rhs = e
}

func (o *OpDiv) Type() ExpressionType {
	return ExpressionOperator
}

func (o *OpDiv) String() string {
	return fmt.Sprintf("(OpDiv: LHS=%v, RHS=%v)", o.lhs, o.rhs)
}
