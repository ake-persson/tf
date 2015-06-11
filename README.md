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
echo '{{ split ":" .Env.PATH | join ",\n" }}' | ./tf
./tf -f examples/example.yaml -t examples/example.conf.tf -o example.conf
./tf -i '{region: amer, country: us}' -t examples/example.conf.tf
./tf -i '{Apples: [1,2,3]}' -t examples/apples.tf
echo '{{range $i, $e := .Etcd}}{{$i}}{{printf "\n"}}{{end}}'|./tf --etcd-node etcd1 --etcd-port 5001 --etcd-key /host
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

# Extended functions amd tests

## Tests

| Test | Arguments | Description |
| - | -  | - |
| last | index (int), array ([]interface{}) | Determine if index is the last element in the array |
| ismap | variable (interface{}) | Test if type is a map (nested data structure) i.e. not printable |

### Examples

## Functions

| Function | Arguments | Description |
| - | - | - |
| join | separator (string), array ([]interface{}) | Join elements in an array to a string |
| split | separator (string), string | Split string into an array |
| repeat | count (int) string | Repeat string x number of times |
| keys | variable (interface{}) | Get keys from interface{} |
| type | variable (interface{}) | Get data type (usefull for debugging templates) |
| nl | | Return new-line |

### Examples
