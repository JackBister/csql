package main

import (
	"encoding/csv"
	"flag"
	"os"

	"github.com/jackbister/csql/csql"
)

var printTypes = flag.Bool("types", false, "")
var separator = flag.String("sep", ",", "")
var skip = flag.Int("skip", 0, "")

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		panic("No query provided")
	}

	options := csql.NewOptions()
	options.PrintTypes = *printTypes
	options.Separator = *separator
	options.Skip = *skip

	query := args[0]

	tokens := csql.Tokenize(query)
	operations := csql.Parse(tokens)

	rs, err := csql.Execute(operations, os.Stdin, options)
	if err != nil {
		panic(err)
	}
	csvWriter := csv.NewWriter(os.Stdout)
	csvWriter.WriteAll(rs)
}
