# tf - Template File

Template Dockerfile or any file in Bash using YAML input and GO text template (http://golang.org/pkg/text/template/).

# Usage

```
Usage:
  tf [OPTIONS]

Application Options:
  -v, --verbose         Verbose
      --version         Version
  -c, --config=         YAML, TOML or JSON config file
  -i, --input=          Input, defaults to using YAML
  -F, --input-format=   Data serialization format YAML, TOML or JSON (YAML)
  -f, --input-file=     Input file, data serialization format used is based on the file extension
  -t, --template=       Template file
  -o, --output-file=    Output file (STDOUT)
  -p, --permission=     File permissions in octal (644)
  -O, --owner=          File Owner
      --etcd-host=      Etcd Host
      --etcd-port=      Etcd Port (2379)
      --etcd-dir=       Etcd Dir (/)
      --http-url=       HTTP Url
      --http-header=    HTTP Header (Accept: application/json)
      --http-format=    HTTP Format (JSON)
      --mysql-user=     MySQL user
      --mysql-password= MySQL password
      --mysql-host=     MySQL host
      --mysql-port=     MySQL port (3306)
      --mysql-database= MySQL database
      --mysql-query=    MySQL query

Help Options:
  -h, --help            Show this help message
```

Input will have it's own namespace such as Arg, File, Env, Etcd. you can also get this by:

```bash
echo '{{keys .}} | tf
```

Argument input will also be in the root scope for convenience.

# Configuration file

Configuration file is also a template i.e. you can use .Env and .Arg for customizing inputs.

## Defaults

Key       | Description | Default
--------- | ----------- | -------
etcd_host | Default Etcd node. |
etcd_port | Default Etcd port. | 2379
http_header | HTTP accept header. | application/json
http_format | Format used by the http response JSON, YAML or TOML. | JSON
mysql_user | Default MySQL user. |
mysql_password | Default MySQL password. |
mysql_host | Default MySQL host. |
mysql_port | Default MySql port | 3306
mysql_database | MySQL database. |

**Example:**

```
[defaults]
mysql_user = "test"
mysql_pass = "test"
mysql_host = "mysql.example.com"
mysql_port = 3306
mysql_database = "test"
```

## Inputs

### Generic

Key | Description | Default
----| ----------- | -------
name | Name of input in data namespace. | Name given in [inputs.<name>].
type | Type of input file, etcd, http, mysql. |

### Specific

Type | Key | Description | Default
---- | --- | ----------- | -------
file | path | Path to input file, format will be determined by file extension .yaml, .json or .toml.
etcd | etcd_host | Etcd node to connect to.
etcd | etcd_port | Etcd port to connect to. | 2379
etcd | etcd_dir | Etcd directory to query, this will be queried recursively.
http | http_url | HTTP url to request.
http | http_header | HTTP accept headers to use for request. Optional will default to JSON.
http | http_format | Format used by the http response JSON, YAML or TOML.
mysql | mysql_user | MySQL user for connection.
mysql | mysql_password | MySQL password for connection.
mysql | mysql_host | MySQL host to connect to.
mysql | mysql_port | MySQL post to connect to.
mysql | mysql_database | MySQL database to connect to.
mysql | mysql_query | MySQL SQL query.

## Example with Etcd

```bash
IP=localhost
```

If you're using Boot2Docker or docker-machine:

```bash
IP=$(echo ${DOCKER_HOST##tcp://} | awk -F: '{ print $1}')
```

```bash
docker run -d --name etcd -p 4001:4001 coreos/etcd:v0.4.6
curl -L http://${IP}:4001/v2/keys/hosts/host1.example.com/serialno -XPUT -d value="abc123"
curl -L http://${IP}:4001/v2/keys/hosts/host2.example.com/serialno -XPUT -d value="def456"
curl -L http://${IP}:4001/v2/keys/hosts/host3.example.com/serialno -XPUT -d value="fgh789"
printf '{{ range $k, $e := .Etcd }}| {{ $k | lalign 20 }} | {{ $e.serialno | lalign 10 }} |\n{{end}}' | \
tf --etcd-host ${IP} --etcd-port 4001 --etcd-dir /hosts
tf --config examples/etcd/tf.toml --template examples/etcd/hosts.tf --input "{ EtcdHost: ${IP} }"
```

You can find a more complete example in "examples/docker" for templating Dockerfile and configuration files, this was
primary use-case for this project. However it's pretty generic and could be used for any templating in Bash.

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
echo '{{range $i, $e := .Apples}}Apple: {{$e}}{{if last $i $.Apples | not}}{{printf ",\n"}}{{end}}{{end}}' | tf -i '{ Apples: [ 1, 2, 3] }'
echo '{{range $k, $e := .Oranges}}{{if map $e | not }}{{printf "%s: %v\n" $k $e}}{{end}}{{end}}' | tf -i '{ Oranges: { a: 1, b: 2, c: { a: 1, b: 2 } } }'
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
echo '{{ 2 | add 2 | sub 2 | mul 5 | div 5}}' | tf 
```

# Build

```bash
git clone https://github.com/mickep76/tf.git
cd tf
./build
bin/tf --version
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

- Add test's
- LDAP support
- Add sort array asc. and desc. templ. func.
- Examples with Etcd data
- Examples with MySQL data
- Validation of data input using schema
- Config that is evaluated in sorted file order for consequtive queries using prev. values
  "/etc/tf.d/01-http", "/etc/tf.d/02-http" etc.
