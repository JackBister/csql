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
