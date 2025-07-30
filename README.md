# Comma Separated Query Language (CSQL)

CSQL is a language for querying CSV files, with a syntax similar to CSV itself. The result is a very succinct language for quick ad hoc querying of CSV files.

For example, given a CSV containing a list of trades with a format like this:

```
Ticker,Quantity,Price
AAPL,100,100
BRK B,100,500
...
```

You could find the largest trades in AAPL with a price greater than 100 with a CSQL query like this:

```
=AAPL,,>100
order($1*$2,desc)
```

Or you could summarize the traded values by ticker with a query like this:

```
group(),sum($1*$2)
order(,desc)
```

While CSQL started out as a joke, it is pretty powerful. It is useful for queries that are in the middle ground where something is too complex to do using shell tools like grep, but simple enough that loading the CSV into SQLite or another database feels like overkill.


- [Comma Separated Query Language (CSQL)](#comma-separated-query-language-csql)
- [Usage](#usage)
  - [Command Line Flags](#command-line-flags)
    - [`-ops`](#-ops)
    - [`-sep=<STR>`](#-sepstr)
    - [`-skip=<N>`](#-skipn)
    - [`-types`](#-types)
- [Language](#language)
  - [Operations](#operations)
    - [Filtering operations](#filtering-operations)
    - [Projecting operations](#projecting-operations)
    - [Grouping operations](#grouping-operations)
    - [Aggregating operations](#aggregating-operations)
    - [Ordering operations](#ordering-operations)
    - [Limiting operations](#limiting-operations)
  - [Supported operations](#supported-operations)
  - [Operands](#operands)
    - [Literals](#literals)
      - [Datetime literals](#datetime-literals)
    - [Column references](#column-references)
      - [Implicit column references](#implicit-column-references)
- [Examples](#examples)
  - [Find all rows where the first column is equal to "ABC"](#find-all-rows-where-the-first-column-is-equal-to-abc)
  - [Find all rows where the first column contains the string "ABC"](#find-all-rows-where-the-first-column-contains-the-string-abc)
  - [Summarize a column across all rows](#summarize-a-column-across-all-rows)
  - [Summarize a column grouped by another column](#summarize-a-column-grouped-by-another-column)
  - [Find the top 10 rows with the highest value in the first column](#find-the-top-10-rows-with-the-highest-value-in-the-first-column)
  - [Find all rows where the first column is equal to "ABC", and include only the second column in the result](#find-all-rows-where-the-first-column-is-equal-to-abc-and-include-only-the-second-column-in-the-result)
  - [Count the number of rows grouped by the first column](#count-the-number-of-rows-grouped-by-the-first-column)
  - [Find all rows where the first column is NOT equal to "ABC"](#find-all-rows-where-the-first-column-is-not-equal-to-abc)


# Usage

```
csql [-ops] [-sep=<STR>] [-skip=<N>] [-types] <query>
```

The input CSV file to be queried must be provided on stdin. To query a CSV file, you can use the `<` operator in your shell, like so:

```
csql '=' < myfile.csv
```

## Command Line Flags

The following flags are available:

### `-ops`

Prints the parsed operations before executing. Used for debugging.

### `-sep=<STR>`

Sets the column separator to `STR`. Defaults to `,`

### `-skip=<N>`

Skips the first `N` lines in the input. This can be used to skip any header rows in the input. CSQL does not automatically detect column headers, so if your input has them, you must use `-skip=1`.

### `-types`

Prints the types of the columns in the result. Used for debugging.

# Language

A CSQL query consists of multiple steps separated by new lines.

 Each step consists of a list of operations.

When the query executes, it will apply the first step to each line in the input CSV and create an output result set.

The next step will then execute on each line in the output result set of the previous step. This process repeats until there are no steps left.

## Operations

Operations can be divided into the following types:

### Filtering operations
Any operation which returns a boolean will be used as a filtering operation. If the returned boolean value is false, the current line will be excluded from the result set.

```sh
echo '1
2' | csql '=2'
2
```

### Projecting operations
Any operation which returns a non-boolean value will be included into the result set.

```sh
echo '1
2' | csql '1+2,3+4'
3,7
3,7
```

### Grouping operations
`group()` can be used to group rows in the result set:

```sh
echo 'A,1
A,2
B,1
B,2' | csql 'group()'
A
B
```

`group()` can reference multiple columns:

```sh
echo 'A,A,1
A,A,2
A,B,1
B,A,1' | csql 'group($0,$1)'
A,A
A,B
B,A
```

### Aggregating operations
Aggregating operations operate across multiple rows in the input and output aggregated values in the result set.

```sh
echo '1
2
3
4' | csql 'sum()'
10
```

Aggregating operations can be combined with grouping operations:

```sh
echo 'A,1
A,2
B,3
B,4' | csql 'group(),sum()'
A,3
B,7
```

### Ordering operations

`order(<column>,<asc|desc>)` can be used to sort the result set:

```sh
echo '1
2
3
4' | csql 'order(,desc)'
4
3
2
1
```

### Limiting operations

`limit(<n>)` can be used to limit the number of rows in the result set:

```sh
echo '1
2
3
4' | csql 'limit(2)'
1
2
```

## Supported operations

The following operations are currently supported:

* `+`
* `-`
* `*`
* `/`
* `=`
* `!`
* `>`
* `<`
* `has(<haystack>,<needle>)`
* `group()`
* `sum()`
* `order(<x>,<asc|desc>)`
* `limit(<n>)`

## Operands

The operators need something to operate on. Generally speaking, there are two types of operands in CSQL:

### Literals

CSQL supports literals of the following types:

* Booleans (true/false)
* Integers
* Floats
* Datetimes
* Unquoted strings

Quoted strings are NOT supported.

CSQL tries to parse literals in this order:

* Booleans
* Integers
* Floats
* Datetimes
* Unquoted strings

This means that in case there is any ambiguity, the first matching type will be used.

#### Datetime literals
Parsing of datetimes is done using [github.com/araddon/dateparse](https://github.com/araddon/dateparse). This library is very flexible in what formats it accepts and should handle most reasonable strings which look like datetimes. See the linked github page for the full list of supported formats.

### Column references

Column references reference the value in a column on the current row being operated on. `$0` references the first column, `$1` the second, etc.

#### Implicit column references
If a query does not contain a literal or column reference in a spot where one is expected, CSQL will implicitly fill that spot with a reference to the column with the same index as the current operation.

Some examples of this:

| Query           | Equivalent to       |
| --------------- | ------------------- |
| `=`             | `$0=$0`             |
| `+,-`           | `$0+$0,$1-$1`       |
| `group(),sum()` | `group($0),sum($1)` |
| `+$1`           | `$0+$1`             |

# Examples

## Find all rows where the first column is equal to "ABC"

`=ABC`

## Find all rows where the first column contains the string "ABC"

`has($0,ABC)`

## Summarize a column across all rows

`sum($1)`

## Summarize a column grouped by another column

`group($0),sum($1)`

## Find the top 10 rows with the highest value in the first column

```
order($0,desc)
limit(10)
```

## Find all rows where the first column is equal to "ABC", and include only the second column in the result

```
=ABC
$1
```

## Count the number of rows grouped by the first column

```
group(),sum(1)
```

## Find all rows where the first column is NOT equal to "ABC"

```
!=ABC
```
