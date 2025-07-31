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

type Options struct {
	PrintOps   bool
	PrintTypes bool
	Separator  string
	Skip       int
}

func NewOptions() Options {
	return Options{
		PrintOps:   false,
		PrintTypes: false,
		Separator:  ",",
		Skip:       0,
	}
}
