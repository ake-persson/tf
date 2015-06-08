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
    etcd "github.com/coreos/go-etcd/etcd"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func map_print(y map[string]interface{}, dir string, pad string) {
    for key, val := range y {
        if reflect.ValueOf(val).Kind() == reflect.Map {
            fmt.Printf("%v[%v]\n", pad, key)
            map_print(val.(map[string]interface{}), dir + "/" + key, pad + "    ")
        } else {
            fmt.Printf("%v%v: %v (%v)\n", pad, key, val, dir)
        }
    }
}

func node_tree(node *etcd.Node, rootKey string, vars map[string]interface{}) error {
    for _, node := range node.Nodes {
        keys := strings.Split(node.Key, "/")
        key := keys[len(keys) - 1]
        if node.Dir {
            vars[key] = make(map[string]interface{})
            node_tree(node, fmt.Sprintf("%s/%s", rootKey, key), vars[key].(map[string]interface{}))
        } else {
            vars[key] = node.Value
        }
    }
    return nil
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
    // get the FileInfo struct describing the standard input.
    fi, _ := os.Stdin.Stat()

    // Initialize log
    log := log.New(os.Stderr, "", 0)

    // Options
    var opts struct {
        Verbose bool `short:"v" long:"verbose" description:"Verbose"`
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
    if err != nil {
        os.Exit(1)
    }

    // Get YAML input
    var input []byte
    if opts.Input != "" {
        if opts.InpFile != "" {
            log.Printf("Can't specify both --input (-i) and --input-file (-f)\n")
            os.Exit(1)
        }
        input = []byte(opts.Input)
    } else if opts.InpFile != "" {
        if _, err := os.Stat(opts.InpFile); os.IsNotExist(err) {
            log.Printf("File doesn't exist: %v\n", opts.InpFile)
            os.Exit(1)
        }
        c, err := ioutil.ReadFile(opts.InpFile)
        input = c
        check(err)
    } else {
        input = []byte("{}")
    }

    // Decode YAML 
    var y map[string]interface{}
    err = yaml.Unmarshal(input, &y)
    check(err)

    env := make(map[string]string)
    for _, e := range os.Environ() {
        v := strings.Split(e, "=")
        env[v[0]] = v[1]
    }

    y["Env"] = env

//    s, err := yaml.Marshal(&y)
//    fmt.Printf("%s\n", string(s))

    vars := make(map[string]interface{})
    if opts.EtcdNode != "" {
        node := []string{fmt.Sprintf("http://%v:%v", opts.EtcdNode, opts.EtcdPort)}
        client := etcd.NewClient(node)
        res, _ := client.Get(opts.EtcdKey, true, true)
        err = node_tree(res.Node, opts.EtcdKey, vars)
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
