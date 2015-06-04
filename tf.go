package main

import (
    "os"
    "fmt"
    "log"
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

func main() {
    // Initialize log
    l := log.New(os.Stderr, "", 0)

    // Options
    var opts struct {
        Debug bool `short:"d" long:"debug" description:"Debug"`
        Input string `short:"y" long:"input" description:"YAML input"`
        InpFile string `short:"i" long:"input-file" description:"YAML input file" default:"default.yaml"`
        TemplFile string `short:"f" long:"template-file" description:"Template file"`
        OutpFile string `short:"o" long:"output-file" description:"Output file"`
        TemplDir string `short:"t" long:"template-dir" description:"Template files with ext. \".tf\" in directory"`
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
            l.Printf("Can't specify both --input (-y) and --input-file (-i)\n")
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

    // Template file
    if opts.TemplFile != "" {
        if opts.TemplDir != "" {
            l.Printf("Can't specify both --template-file (-i) and --template-dir (-t)\n")
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
        err = t.Execute(buf, input) 
        check(err)

        fmt.Printf("%v\n", input)
        fmt.Printf("%v\n", buf)
    }
}
