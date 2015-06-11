package main

import (
    "os"
    "fmt"
    "log"
    "strings"
    "io/ioutil"
    flags "github.com/jessevdk/go-flags"
    "gopkg.in/yaml.v2"
    "text/template"
    "bytes"
    etcd "github.com/coreos/go-etcd/etcd"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

var fns = template.FuncMap{
    "last": arrLast,
    "join": arrJoin,
    "split": strSplit,
    "repeat": strRepeat,
    "keys": intfKeys,
    "type": intfType,
    "ismap": intfIsMap,
    "upper": strings.ToUpper,
    "lower": strings.ToLower,
    "contains": strings.Contains,
    "replace": strReplace,
    "trim": strTrim,
    "ltrim": strTrimLeft,
    "rtrim": strTrimRight,
    "default": intfDefault,
    "center": strCenter,
    "random": intRandom,
    "capitalize": strCapitalize,
    "add": intAdd,
    "sub": intSub,
    "div": intDiv,
    "mul": intMul,
    "lalign": strAlignLeft,
    "ralign": strAlignRight,
}

func main() {
    // get the FileInfo struct describing the standard input.
    fi, _ := os.Stdin.Stat()

    // Initialize log
    log := log.New(os.Stderr, "", 0)

    // Options
    var opts struct {
        Verbose bool `short:"v" long:"verbose" description:"Verbose"`
        Version bool `long:"version" description:"Version"`
        Input string `short:"i" long:"input" description:"YAML input"`
        InpFile string `short:"f" long:"input-file" description:"YAML input file"`
        TemplFile string `short:"t" long:"template-file" description:"Template file"`
        OutpFile string `short:"o" long:"output-file" description:"Output file (STDOUT)"`
        Permission int `short:"p" long:"permission" description:"Permission for output file" default:"644"`
        EtcdNode string `short:"n" long:"etcd-node" description:"Etcd Node"`
        EtcdPort int `short:"P" long:"etcd-port" description:"Etcd Port" default:"2379"`
        EtcdKey string `short:"k" long:"etcd-key" description:"Etcd Key" default:"/"`
    }

    // Parse options
    _, err := flags.Parse(&opts)
    check(err)

    // Print version
    if opts.Version == true {
        fmt.Printf("%s\n", version)
        os.Exit(0)
    }

    // Get YAML input
    var inp []byte
    if opts.Input != "" {
        if opts.InpFile != "" {
            log.Fatal("Can't specify both --input (-i) and --input-file (-f)\n")
            os.Exit(1)
        }
        inp = []byte(opts.Input)
    } else if opts.InpFile != "" {
        if _, err := os.Stat(opts.InpFile); os.IsNotExist(err) {
            log.Printf("File doesn't exist: %v\n", opts.InpFile)
            os.Exit(1)
        }
        content, err := ioutil.ReadFile(opts.InpFile)
        check(err)
        inp = content
    } else {
        inp = []byte("{}")
    }

    // Decode YAML 
    var y map[string]interface{}
    err = yaml.Unmarshal(inp, &y)
    check(err)

    // Get environment
    env := make(map[string]string)
    for _, e := range os.Environ() {
        v := strings.Split(e, "=")
        env[v[0]] = v[1]
    }

    y["Env"] = env

//    s, err := yaml.Marshal(&y)
//    fmt.Printf("%s\n", string(s))

    vars := make(map[interface{}]interface{})
    if opts.EtcdNode != "" {
        node := []string{fmt.Sprintf("http://%v:%v", opts.EtcdNode, opts.EtcdPort)}
        client := etcd.NewClient(node)
        res, _ := client.Get(opts.EtcdKey, true, true)
        etcdNestedMap(res.Node, vars)
        y["Etcd"] = vars
    }

//    s, err := yaml.Marshal(&y)
//    fmt.Printf("%s\n", string(s))

    // Template input
    var templ string
    if (fi.Mode() & os.ModeCharDevice) == 0 {
        bytes, _ := ioutil.ReadAll(os.Stdin)
        templ = string(bytes)
    } else if opts.TemplFile != "" {
        if _, err := os.Stat(opts.TemplFile); os.IsNotExist(err) {
            log.Printf("File doesn't exist: %v\n", opts.TemplFile)
            os.Exit(1)
        }

        // Open file
        c, err := ioutil.ReadFile(opts.TemplFile)
        check(err)
        templ = string(c)
    } else {
        log.Printf("No template specified using --template-file (-t) or piped to STDIN\n")
        os.Exit(1)
    }

    // Parse template
    t := template.Must(template.New("template").Funcs(fns).Parse(templ))

    buf := new(bytes.Buffer)
    err = t.Execute(buf, y)
    check(err)

    // Write result
    if opts.OutpFile != "" {
        err := ioutil.WriteFile(opts.OutpFile, buf.Bytes(), os.FileMode(opts.Permission))
        check(err)
    } else {
        fmt.Printf("%v\n", buf)
    }
}
