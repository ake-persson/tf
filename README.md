# tf - Template File

Template file in Bash using YAML input and GO text template (http://golang.org/pkg/text/template/).

# Usage

```bash
Usage:
  tf [OPTIONS]

Application Options:
  -v, --verbose        Verbose
      --version        Version
  -i, --input=         YAML input
  -f, --input-file=    YAML input file
  -t, --template-file= Template file
  -o, --output-file=   Output file (STDOUT)
  -p, --permission=    File permissions in octal (644)
  -O, --owner=         File Owner
  -n, --etcd-node=     Etcd Node
  -P, --etcd-port=     Etcd Port (2379)
  -k, --etcd-key=      Etcd Key (/)

Help Options:
  -h, --help           Show this help message
```

# Examples

```bash
./tf -f examples/example.yaml -t examples/example.conf.tf -o example.conf
./tf -i '{region: amer, country: us}' -t examples/example.conf.tf
./tf -i '{Apples: [1,2,3]}' -t examples/apples.tf
echo '{{keys .Etcd | join "\n"}}' | tf --etcd-node etcd1 --etcd-port 5001 --etcd-key /host
```

# Extended functions and tests

## Tests

Test     | Arguments           | Types              | Description
-------- | ------------------- | ------------------ | -----------
last     | $index $array       | int, []interface{} | Determine if $index is the last element in the $array
ismap    | $variable           | $interface{}       | Test if $variable type is a map (nested data structure)
contains | $string $sub-string | string, string     | Test if $string contains $sub-string
even     | $x                  | int                | Test if $x is even
odd      | $x                  | int                | Test if $x is odd

### Examples

```bash
echo '{{range $i, $e := .Inp.Apples}}Apple: {{$e}}{{if last $i $.Inp.Apples | not}}{{printf ",\n"}}{{end}}{{end}}' | tf -i '{ Apples: [ 1, 2, 3] }'
echo '{{range $k, $e := .Inp.Oranges}}{{if ismap $e | not }}{{printf "%s: %v\n" $k $e}}{{end}}{{end}}' | tf -i '{ Oranges: { a: 1, b: 2, c: { a: 1, b: 2 } } }'
echo '{{1 | even }} | tf
```

## Functions

Function   | Arguments          | Types                       | Description
---------- | ------------------ | --------------------------- | -----------
join       | $separator $array  | string, []interface{}       | Join elements in an array to a string
split      | $separator $string | string, string              | Split string into an array
repeat     | $count $string     | int, string                 | Repeat string x number of times
keys       | $variable          | interface{}                 | Get keys from interface{}
type       | $variable          | interface{}                 | Get data type (usefull for debugging templates)
lower      | $string            | string                      | Convert string to lower case
upper      | $string            | string                      | Convert string to upper case
replace    | $old $new $string  | string, string, string      | Replace old with new in string
trim       | $trim $string      | string, string              | Trim preceding and trailing characters
ltrim      | $trim $string      | string, string              | Trim preceding characters
rtrim      | $trim $string      | string, string              | Trim trailing characters
default    | $default $optional | interface{}, ...interface{} | If no value is passed for the second arg. it returns the default
center     | $size $string      | int, string                 | Center text
ralign     | $size $string      | int, string                 | Right align text
lalign     | $size $string      | int, string                 | Left align text
capitalize | $string            | string                      | Capitalize first character in string
add        | $y $x              | int, int                    | Addition, arguments are in reverse order to allow pipeline
sub        | $y $x              | int, int                    | Subtraction, arguments are in reverse order to allow pipeline
div        | $y $x              | int, int                    | Division, arguments are in reverse order to allow pipeline
mul        | $y $x              | int, int                    | Multiplication, arguments are in reverse order to allow pipeline
date       | $fmt               | ...interface{}              | Print date/time, optional argument strftime syntax

### Examples

```bash
echo '{{split ":" .Env.PATH | join ",\n"}}' | tf
echo '{{repeat 20 "-"}} HELLO WORLD! {{"-" | repeat 20}}' | tf
echo '{{keys .Env | join "\n"}}' | tf
echo '{{ "UPPER" | lower}} {{ "lower" | upper }}' | tf
echo '{{ "Doo Doo" | replace "Doo" "Doo is extinct" }}' | tf
echo '{{ "!!! TRIM !!!" | trim "! " }}' | tf
echo '{{ 2 | add 2 | sub 2 | mul 5 | div 5}}' | ./tf 
```

# Build

```bash
go get github.com/mickep76/tf
go install github.com/mickep76/tf
$GOPATH/bin/tf
```

# Install using Homebrew

```bash
brew tap mickep76/funk-gnarge
brew install mickep76/funk-gnarge/tf
```

# Template a directory structure

## Template

```bash
input='input.yaml'
for file in $(find . -type f -name '*.tf'); do
    tf -i ${input} -t ${file} -o ${file%%.tf}
done
```

## Cleanup

```bash
for file in $(find . -type f -name '*.tf'); do
    rm -f ${file%%.tf}
done
```

## Use in Makefile

```
INPUT=input.yaml

all: build

clean:
        for file in $$(find . -type f -name '*.tf'); do \
                rm -f $${file%%.tf} ; \
        done

build: clean
        for file in $$(find . -type f -name '*.tf'); do \
                tf -f ${INPUT} -t $${file} -o $${file%%.tf} ; \
        done
```

# Issues

Currently Go text/template doesn't have a way to [suppress newlines](https://github.com/golang/go/issues/9969).

# Roadmap

- Add multiple input sources
  + JSON/TOML support
  + REST query support
- Add config file
  + Add support for multiple input sources
  + Allow for sequencing multiple input source to have hierarhical substitution
- Add templating inside input files
- Add sort array asc. and desc. templ. func.
