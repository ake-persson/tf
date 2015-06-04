# tf - Template File

Template files in Bash using YAML input and GO HTML templating.

# Usage

```bash
Usage:
  tf [OPTIONS]

Application Options:
  -v, --verbose        Verbose
  -i, --input=         YAML input
  -f, --input-file=    YAML input file (default.yaml)
  -t, --template-file= Template file
  -o, --output-file=   Output file, will use stdout per default
  -d, --template-dir=  Template files with ext. ".tf" in directory

Help Options:
  -h, --help           Show this help message
```

# Examples

```bash
./tf -f example.yaml -t example.conf.tf -o example.conf
./tf -i '{region: amer, country: us}' -t example.conf.tf
```

# Build

```bash
go get github.com/mickep76/tf
go install github.com/mickep76/tf
$GOPATH/bin/tf
```
