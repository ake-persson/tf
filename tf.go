package main

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/go-etcd/etcd"
	flags "github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var fns = template.FuncMap{
	"last":       arrLast,
	"join":       arrJoin,
	"split":      strSplit,
	"repeat":     strRepeat,
	"keys":       intfKeys,
	"type":       intfType,
	"ismap":      intfIsMap,
	"upper":      strings.ToUpper,
	"lower":      strings.ToLower,
	"contains":   strings.Contains,
	"replace":    strReplace,
	"trim":       strTrim,
	"ltrim":      strTrimLeft,
	"rtrim":      strTrimRight,
	"default":    intfDefault,
	"center":     strCenter,
	"random":     intRandom,
	"capitalize": strCapitalize,
	"add":        intAdd,
	"sub":        intSub,
	"div":        intDiv,
	"mul":        intMul,
	"lalign":     strAlignLeft,
	"ralign":     strAlignRight,
	"odd":        odd,
	"even":       even,
	"date":       date,
}

type Input struct {
	Name string
	Type string
	Path string
}

type Merge struct {
    Name string
    Inputs []interface{}
}

func main() {
	// get the FileInfo struct describing the standard input.
	fi, _ := os.Stdin.Stat()

	// Set log options
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)

	// Options
	var opts struct {
		Verbose    bool   `short:"v" long:"verbose" description:"Verbose"`
		Version    bool   `long:"version" description:"Version"`
		Config     string `short:"c" long:"config" description:"YAML, TOML or JSON config file"`
		Input      string `short:"i" long:"input" description:"Input, defaults to using YAML"`
		InpFormat  string `short:"F" long:"input-format" description:"Data serialization format YAML, TOML or JSON" default:"YAML"`
		InpFile    string `short:"f" long:"input-file" description:"Input file, data serialization format used is based on the file extension"`
		TemplFile  string `short:"t" long:"template-file" description:"Template file"`
		OutpFile   string `short:"o" long:"output-file" description:"Output file (STDOUT)"`
		Permission string `short:"p" long:"permission" description:"File permissions in octal" default:"644"`
		Owner      string `short:"O" long:"owner" description:"File Owner"`
		EtcdNode   string `short:"n" long:"etcd-node" description:"Etcd Node"`
		EtcdPort   int    `short:"P" long:"etcd-port" description:"Etcd Port" default:"2379"`
		EtcdKey    string `short:"k" long:"etcd-key" description:"Etcd Key" default:"/"`
	}

	// Parse options
	_, err := flags.Parse(&opts)
	if err != nil {
		ferr := err.(*flags.Error)
		if ferr.Type == flags.ErrHelp {
			os.Exit(1)
		}
	}
	check(err)

	// Print version
	if opts.Version == true {
		fmt.Println(version)
		os.Exit(0)
	}

	// Get argument input
	y := make(map[string]interface{})
	if opts.Input != "" {
		var fmt DataFmt
		switch opts.InpFormat {
		case "YAML":
			fmt = YAML
		case "TOML":
			fmt = TOML
		case "JSON":
			fmt = JSON
		default:
			log.Fatal("Unsupported data format, needs to be YAML, JSON or TOML")
		}
		v, err := UnmarshalData([]byte(opts.Input), fmt)
		check(err)
		y = v
		y["Arg"] = v
	}

	// Get file input
	if opts.InpFile != "" {
		v, err := LoadFile(opts.InpFile)
		check(err)
		y["File"] = v
	}

	// Get environment
	env := GetOSEnv()
	y["Env"] = env

	// Load config file
	if opts.Config != "" {
		c, err := LoadFile(opts.Config)
		check(err)
		y["Cfg"] = c

		if reflect.ValueOf(c["inputs"]).Kind() == reflect.Map {
			for k1, v1 := range c["inputs"].(map[string]interface{}) {
				var inp Input
				inp.Name = k1
				for k2, v2 := range v1.(map[string]interface{}) {
					switch k2 {
					case "name":
						inp.Name = v2.(string)
					case "type":
						inp.Type = v2.(string)
					case "path":
						inp.Path = v2.(string)
                    default:
                        log.Fatalf("Invalid key in configuration file inputs.%v.%v", k1, k2)
					}
				}
				switch inp.Type {
				case "file":
					c, err := LoadFile(inp.Path)
					check(err)
					if y[inp.Name] == nil {
						y[inp.Name] = c
					} else {
						log.Fatalf("Namespace already exist's: %v", inp.Name)
					}
				default:
					log.Fatalf("Unknown type in config inputs.%v.Type: %v", inp.Name, inp.Type)
				}
			}
		}

        if reflect.ValueOf(c["merge"]).Kind() == reflect.Map {
            for k1, v1 := range c["merge"].(map[string]interface{}) {
				var m Merge
				m.Name = k1
                for k2, v2 := range v1.(map[string]interface{}) {
                    switch k2 {
					case "name":
						m.Name = v2.(string)
					case "inputs":
						m.Inputs = v2.([]interface{})
					default:
						log.Fatalf("Invalid key in configuration file merge.%v.%v", k1, k2)
					}
				}

				for i := range m.Inputs {
	                if y[m.Name] == nil {
						y2 := make(map[string]interface{})
                        for k, v := range y[m.Inputs[i].(string)].(map[string]interface{}) {
			                y2[k] = v
						}
						y[m.Name] = y2
					} else {
						y2 := y[m.Name].(map[string]interface{})
		                for k, v := range y[m.Inputs[i].(string)].(map[string]interface{}) {
			                y2[k] = v
		                }
					}
				}
			}
		}
	}

	if opts.EtcdNode != "" {
		node := []string{fmt.Sprintf("http://%v:%v", opts.EtcdNode, opts.EtcdPort)}
		client := etcd.NewClient(node)
		res, _ := client.Get(opts.EtcdKey, true, true)
		e := EtcdMap(res.Node)
		y["Etcd"] = e
	}

	if opts.Verbose {
		s, _ := yaml.Marshal(&y)
		fmt.Printf("%s\n", string(s))
	}

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
		p, err := strconv.ParseUint(opts.Permission, 8, 32)
		check(err)

		w, err := os.Create(opts.OutpFile)
		check(err)

		w.Chmod(os.FileMode(p))

		if opts.Owner != "" {
			u, err := user.Lookup(opts.Owner)
			check(err)

			uid, err := strconv.Atoi(u.Uid)
			check(err)

			gid, err := strconv.Atoi(u.Gid)
			check(err)

			err = w.Chown(uid, gid)
			check(err)
		}

		_, err = w.Write(buf.Bytes())
		check(err)

		w.Close()
	} else {
		fmt.Printf("%v\n", buf)
	}
}
