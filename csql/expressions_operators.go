package csql

import "fmt"

type OpEquals struct {
	lhs Expression
	rhs Expression
}

func (o *OpEquals) Execute(i int, record []Value) (*OperationResult, error) {
	lhs, err := o.lhs.Execute(i, record)
	if err != nil {
		return nil, err
	}
	if lhs == nil {
		// TODO: Should nil == nil?
		return &OperationResult{
			value: &Value{
				typ:   ValueTypeBool,
				value: false,
			},
		}, nil
	}
	rhs, err := o.rhs.Execute(i, record)
	if err != nil {
		return nil, err
	}
	if rhs == nil {
		// TODO: Should nil == nil?
		return &OperationResult{
			value: &Value{
				typ:   ValueTypeBool,
				value: false,
			},
		}, nil
	}
	lhsV := lhs.value
	rhsV := rhs.value

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

func (o *OpEquals) SetLHS(e Expression) {
	o.lhs = e
}

func (o *OpEquals) SetRHS(e Expression) {
	o.rhs = e
}

func (o *OpEquals) Type() ExpressionType {
	return ExpressionEquals
}

func (o *OpEquals) String() string {
	return fmt.Sprintf("(OpEquals: LHS=%v, RHS=%v)", o.lhs, o.rhs)
}
