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
echo 'PATH: {{.Env.PATH}}:{{.Path}}' | ./tf -i '{Path: /usr/local/bin}'
echo '{{range $i, $e := .Etcd}}{{$i}}{{printf "\n"}}{{end}}'|./tf --etcd-node etcd1 --etcd-port 5001 --etcd-key /host
echo '{{ $path := split .Env.PATH ":" }}{{ join $path "|" }}' | ./tf
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
