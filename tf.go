package main

import (
    "os"
    "fmt"
    "log"
    "strings"
    "strconv"
    "reflect"
    "io/ioutil"
    flags "github.com/jessevdk/go-flags"
    "gopkg.in/yaml.v2"
    "text/template"
    "bytes"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

var fns = template.FuncMap{
    "last": func(x int, a interface{}) bool {
        return x == reflect.ValueOf(a).Len() - 1
    },
    "join": func(a []interface{}, sep string) string {
        s := make([]string, len(a))
        for i, v := range a {
            switch v.(type) {
                case string:
                    s[i] = v.(string)
                case int, int32, int64:
                    s[i] = strconv.Itoa(v.(int))
            }
        }

        return strings.Join(s, sep)
    },
}

func main() {
    // Initialize log
    l := log.New(os.Stderr, "", 0)

    // Options
    var opts struct {
        Verbose bool `short:"v" long:"verbose" description:"Verbose"`
        Input string `short:"i" long:"input" description:"YAML input"`
        InpFile string `short:"f" long:"input-file" description:"YAML input file"`
        TemplFile string `short:"t" long:"template-file" description:"Template file"`
        OutpFile string `short:"o" long:"output-file" description:"Output file, will use stdout per default"`
    }

    // Parse options
    _, err := flags.Parse(&opts)
    if err != nil {
        os.Exit(1)
    }

    // Get YAML input
    var input []byte
    if opts.Input != "" {
        if opts.InpFile != "" {
            l.Printf("Can't specify both --input (-i) and --input-file (-f)\n")
            os.Exit(1)
        }
        input = []byte(opts.Input)
    } else if opts.InpFile != "" {
        if _, err := os.Stat(opts.InpFile); os.IsNotExist(err) {
            l.Printf("File doesn't exist: %v\n", opts.InpFile)
            os.Exit(1)
        }
        c, err := ioutil.ReadFile(opts.InpFile)
        input = c
        check(err)
    } else {
        l.Printf("You need to specify either --input (-i) or --input-file (-f)\n")
        os.Exit(1)
    }

    // Decode YAML 
    var y map[string]interface{}
    err = yaml.Unmarshal(input, &y)
    check(err)

//    s, err := yaml.Marshal(&y)
//    fmt.Printf("%s\n", string(s))

    // Template file
    if opts.TemplFile != "" {
        if _, err := os.Stat(opts.TemplFile); os.IsNotExist(err) {
            l.Printf("File doesn't exist: %v\n", opts.TemplFile)
            os.Exit(1)
        }

        // Open file
        c, err := ioutil.ReadFile(opts.TemplFile)
        check(err)

        // Parse template
        t := template.Must(template.New("template").Funcs(fns).Parse(string(c)))
//        t, err := template.New("template").Funcs(fns).Parse(string(c))
//        check(err)

        buf := new(bytes.Buffer)
        err = t.Execute(buf, y)
        check(err)

        // Write result
        if opts.OutpFile != "" {
            err := ioutil.WriteFile(opts.OutpFile, buf.Bytes(), 0644)
            check(err)
        } else {
            fmt.Printf("%v\n", buf)
        }
    }
}
