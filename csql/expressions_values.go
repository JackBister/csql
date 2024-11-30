package csql

import "fmt"

type Nop struct{}

func (o *Nop) Execute(i int, record []Value) (*OperationResult, error) {
	return &OperationResult{}, nil
}

func (o *Nop) FillNils(e Expression) {
}

func (o *Nop) Type() ExpressionType {
	return ExpressionNop
}

func (o *Nop) String() string {
	return "(Nop)"
}

type LiteralExpression struct {
	value Value
}

func (l *LiteralExpression) Execute(i int, record []Value) (*OperationResult, error) {
	return &OperationResult{
		value: &l.value,
	}, nil
}

func (o *LiteralExpression) FillNils(e Expression) {
}

func (l *LiteralExpression) Type() ExpressionType {
	return ExpressionLiteral
}

func (l *LiteralExpression) String() string {
	return fmt.Sprintf("(Literal: Value=%v)", l.value)
}

type ColumnReferenceExpression struct {
	index int
}

func (c *ColumnReferenceExpression) Execute(i int, record []Value) (*OperationResult, error) {
	if c.index >= len(record) {
		return nil, fmt.Errorf("index out of range, index: %v, record length: %v, record: %v", c.index, len(record), record)
	}
	return &OperationResult{
		value: &record[c.index],
	}, nil
}

func (o *ColumnReferenceExpression) FillNils(e Expression) {
}

func (c *ColumnReferenceExpression) Type() ExpressionType {
	return ExpressionColumnReference
}

func (c *ColumnReferenceExpression) String() string {
	return fmt.Sprintf("(ColumnRef: Index=%v)", c.index)
}
