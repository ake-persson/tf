package main

import (
    "os"
    "fmt"
    "log"
    "path/filepath"
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
    var opts struct {
        Debug bool `short:"d" long:"debug" description:"Debug"`
        Input string `short:"y" long:"input" description:"YAML input"`
        InpFile string `short:"i" long:"input-file" description:"YAML input file" default:"default.yaml"`
        TemplFile string `short:"f" long:"template-file" description:"Template file"`
        OutpFile string `short:"o" long:"output-file" description:"Output file"`
        TemplDir string `short:"t" long:"template-dir" description:"Template files with ext. \".tf\" in directory"`
    }

    _, err := flags.Parse(&opts)
    if err != nil {
        os.Exit(1)
    }

    l := log.New(os.Stderr, "", 0)

    var input []byte
    if opts.Input != "" {
        if opts.InpFile != "default.yaml" {
            l.Printf("Can't specify both --input and --input-file\n")
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
//    fmt.Printf("%v\n", input)

    // Decode YAML
    var y map[string]interface{}
    err = yaml.Unmarshal(input, &y)
    check(err)

    os.Exit(0)

    files, _ := filepath.Glob(opts.TemplDir + "/*.yaml")

    gy := make(map[string]interface{})
    for _, file := range files {
        // Read YAML file
        fmt.Printf("### Read File ###\n\n%s\n\n", file)
        c, err := ioutil.ReadFile(file)
        check(err)
        if opts.Debug { fmt.Printf("### File Content ###\n\n%s\n", c) }

        // Parse template
        tmpl, err := template.New("template").Parse(string(c))
        check(err)

        buf := new(bytes.Buffer)
        err = tmpl.Execute(buf, gy)
        check(err)
        if opts.Debug { fmt.Printf("### Compiled Template ###\n\n%s\n", buf) }

        // Decode YAML
        var y map[string]interface{}
            err = yaml.Unmarshal(buf.Bytes(), &y)
        check(err)

        // Merge global map
        for k, v := range y {
            gy[k] = v
        }
    }

    // Result
    s, err := yaml.Marshal(&gy)
    check(err)
    fmt.Printf("### Result ###\n\n%s\n\n", string(s))
}
