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
