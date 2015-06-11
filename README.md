# tf - Template File

Template file in Bash using YAML input and GO text template (http://golang.org/pkg/text/template/).

# Usage

```bash
age:
  tf [OPTIONS]

Application Options:
  -v, --verbose        Verbose
  -i, --input=         YAML input
  -f, --input-file=    YAML input file
  -t, --template-file= Template file
  -o, --output-file=   Output file (STDOUT)
  -p, --permission=    Permission for output file (644)

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

# Extended functions amd tests

## Tests

Test  | Argument 1             | Argument 2            | Description
----- | ---------------------- | --------------------- | -----------
last  | index (int)            | array ([]interface{}) | Determine if index is the last element in the array
ismap | variable (interface{}) |                       | Test if type is a map (nested data structure) i.e. not printable

### Examples

```bash
echo '{{range $i, $e := .Apples}}Apple: {{$e}}{{if last $i $.Apples | not}}{{printf ",\n"}}{{end}}{{end}}' | tf -i '{ Apples: [ 1, 2, 3] }'
echo '{{range $k, $e := .Oranges}}{{if ismap $e | not }}{{printf "%s: %v\n" $k $e}}{{end}}{{end}}' | tf -i '{ Oranges: { a: 1, b: 2, c: { a: 1, b: 2 } } }'
```

## Functions

Function | Argument 1             | Argument 2            | Description
-------- | ---------------------- | --------------------- | -----------
join     | separator (string)     | array ([]interface{}) | Join elements in an array to a string
split    | separator (string)     | string                | Split string into an array
repeat   | count (int)            | string                | Repeat string x number of times
keys     | variable (interface{}) |                       | Get keys from interface{}
type     | variable (interface{}) |                       | Get data type (usefull for debugging templates)
nl       |                        |                       | Return new-line

### Examples

```bash
echo '{{split ":" .Env.PATH | join ",\n"}}' | tf
echo '{{repeat 20 "-"}} HELLO WORLD! {{"-" | repeat 20}}' | tf
echo '{{keys .Env | join "\n"}}' | tf
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
for file in $(find . -name '*.tf'); do
    tf -i ${input} -t ${file} -o ${file%%.tf}
done
```

## Cleanup

```bash
for file in $(find . -name '*.tf'); do
    rm -f ${file%%.tf}
done
```
