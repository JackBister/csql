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

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"github.com/jackbister/csql/pkg/csql"
)

var versionString string // This must be set using -ldflags "-X main.versionString=<version>" when building for --version to work

var printOps = flag.Bool("ops", false, "Print operations")
var printTypes = flag.Bool("types", false, "")
var printVersion = flag.Bool("version", false, "Print version and exit")
var separator = flag.String("sep", ",", "")
var skip = flag.Int("skip", 0, "")

func main() {
	flag.Parse()

	args := flag.Args()

	if *printVersion {
		if versionString == "" {
			versionString = "unknown"
		}
		fmt.Println(versionString)
		return
	}

	if len(args) < 1 {
		panic("No query provided")
	}

	options := csql.NewOptions()
	options.PrintOps = *printOps
	options.PrintTypes = *printTypes
	options.Separator = *separator
	options.Skip = *skip

	query := args[0]

	tokens := csql.Tokenize(query)
	operations, err := csql.ParseQuery(tokens)
	if err != nil {
		panic(err)
	}

	rs, err := csql.Execute(operations, os.Stdin, options)
	if err != nil {
		panic(err)
	}
	csvWriter := csv.NewWriter(os.Stdout)
	csvWriter.WriteAll(rs)
}
