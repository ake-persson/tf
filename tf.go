package main

import (
    "os"
    "fmt"
    "log"
//    "bufio"
//    "path/filepath"
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

// Accept template from stdin

func main() {
    // Initialize log
    l := log.New(os.Stderr, "", 0)

    // Options
    var opts struct {
        Debug bool `short:"d" long:"debug" description:"Debug"`
        Input string `short:"i" long:"input" description:"YAML input"`
        InpFile string `short:"I" long:"input-file" description:"YAML input file" default:"default.yaml"`
        TemplFile string `short:"t" long:"template-file" description:"Template file"`
        OutpFile string `short:"o" long:"output-file" description:"Output file, will use stdout per default"`
//        TemplDir string `short:"T" long:"template-dir" description:"Template files with ext. \".tf\" in directory"`
    }

    // Parse options
    _, err := flags.Parse(&opts)
    if err != nil {
        os.Exit(1)
    }

    // Get YAML input
    var input []byte
    if opts.Input != "" {
        if opts.InpFile != "default.yaml" {
            l.Printf("Can't specify both --input (-i) and --input-file (-I)\n")
            os.Exit(1)
        }
        input = []byte(opts.Input)
    } else {
        if _, err := os.Stat(opts.InpFile); os.IsNotExist(err) {
            l.Printf("File doesn't exist: %v\n", opts.InpFile)
            os.Exit(1)
        }
        c, err := ioutil.ReadFile(opts.InpFile)
        input = c
        check(err)
    }

    // Decode YAML 
    var y map[string]interface{}
    err = yaml.Unmarshal(input, &y)
    check(err)

    s, err := yaml.Marshal(&y)
    fmt.Printf("%s\n", string(s))

    // Template file
    if opts.TemplFile != "" {
        if opts.TemplDir != "" {
            l.Printf("Can't specify both --template-file (-t) and --template-dir (-T)\n")
            os.Exit(1)
        }

        if _, err := os.Stat(opts.TemplFile); os.IsNotExist(err) {
            l.Printf("File doesn't exist: %v\n", opts.TemplFile)
            os.Exit(1)
        }

        // Open file
        c, err := ioutil.ReadFile(opts.TemplFile)
        check(err)

        // Parse template
        t, err := template.New("template").Parse(string(c))
        check(err)

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
