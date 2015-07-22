package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	log "github.com/mickep76/tf/vendor/github.com/Sirupsen/logrus"
	etcd "github.com/mickep76/tf/vendor/github.com/coreos/go-etcd/etcd"
	flags "github.com/mickep76/tf/vendor/github.com/jessevdk/go-flags"
	"github.com/mickep76/tf/vendor/gopkg.in/yaml.v2"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var fns = template.FuncMap{
	"last":       IsLast,
	"islast":     IsLast,
	"join":       Join,
	"split":      Split,
	"repeat":     Repeat,
	"keys":       Keys,
	"type":       Type,
	"ismap":      IsMap,
	"upper":      strings.ToUpper,
	"lower":      strings.ToLower,
	"contains":   strings.Contains,
	"replace":    Replace,
	"trim":       Trim,
	"ltrim":      TrimLeft,
	"rtrim":      TrimRight,
	"default":    Default,
	"center":     Center,
	"random":     Random,
	"capitalize": Capitalize,
	"add":        Add,
	"sub":        Sub,
	"div":        Div,
	"mul":        Mul,
	"lalign":     AlignLeft,
	"ralign":     AlignRight,
	"odd":        Odd,
	"even":       Even,
	"date":       Date,
}

type Input struct {
	Name       *string
	Type       *string
	Path       *string
	EtcdNode   *string
	EtcdPort   *int64
	EtcdDir    *string
	HttpUrl    *string
	HttpHeader *string
	HttpFormat *string
	MysqlUser  *string
	MysqlPass  *string
	MysqlHost  *string
	MysqlPort  *int64
	MysqlDb    *string
	MysqlQry   *string
}

type Merge struct {
	Name   string
	Inputs []interface{}
}

func main() {
	// Get the FileInfo struct describing the standard input.
	fi, _ := os.Stdin.Stat()

	// Set log options
	log.SetOutput(os.Stderr)
	log.SetLevel(log.WarnLevel)

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
		EtcdDir    string `short:"k" long:"etcd-dir" description:"Etcd Dir" default:"/"`
		HttpUrl    string `short:"u" long:"http-url" description:"HTTP Url"`
		HttpHeader string `short:"H" long:"http-header" description:"HTTP Header" default:"Accept: application/json"`
		HttpFormat string `long:"http-format" description:"HTTP Format" default:"JSON"`
		MysqlUser  string `long:"mysql-user" description:"MySql user"`
		MysqlPass  string `long:"mysql-pass" description:"MySQL password"`
		MysqlHost  string `long:"mysql-host" description:"MySQL host"`
		MysqlPort  int64  `long:"mysql-port" description:"MySQL port" default:"3306"`
		MysqlDb    string `long:"mysql-db" description:"MySQL database"`
		MysqlQry   string `long:"mysql-query" description:"MySQL query"`
	}

	// Parse options
	if _, err := flags.Parse(&opts); err != nil {
		ferr := err.(*flags.Error)
		if ferr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			log.Fatal(err.Error())
		}
	}

	// Print version
	if opts.Version {
		fmt.Printf("tf %s\n", version)
		os.Exit(0)
	}

	// Set verbose
	if opts.Verbose {
		log.SetLevel(log.InfoLevel)
	}

	// Get environment variables
	data := make(map[string]interface{})
	data["Env"] = GetOSEnv()

	// Get argument input
	if opts.Input != "" {
		var f DataFmt
		switch opts.InpFormat {
		case "YAML":
			f = YAML
		case "TOML":
			f = TOML
		case "JSON":
			f = JSON
		default:
			log.Fatal("Unsupported data format, needs to be YAML, JSON or TOML")
		}

		var err error
		data["Arg"], err = UnmarshalData([]byte(opts.Input), f)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Copy .Arg namespace to . for convenience
		for k, v := range data["Arg"].(map[string]interface{}) {
			data[k] = v
		}
	}

	// Get file input
	if opts.InpFile != "" {
		var err error
		data["File"], err = LoadFile(opts.InpFile, data)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// Get Etcd input
	if opts.EtcdNode != "" {
		// Add error handling
		node := []string{fmt.Sprintf("http://%v:%v", opts.EtcdNode, opts.EtcdPort)}
		client := etcd.NewClient(node)
		res, _ := client.Get(opts.EtcdDir, true, true)
		data["Etcd"] = EtcdMap(res.Node)
	}

	// Get http input
	if opts.HttpUrl != "" {
		var f DataFmt
		switch opts.HttpFormat {
		case "YAML":
			f = YAML
		case "TOML":
			f = TOML
		case "JSON":
			f = JSON
		default:
			log.Fatal("Unsupported data format, needs to be YAML, JSON or TOML")
		}

		var err error
		data["Http"], err = GetHTTP(opts.HttpUrl, opts.HttpHeader, f)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// Get MySQL input
	if opts.MysqlHost != "" {
		var err error
		data["Mysql"], err = GetMySQL(opts.MysqlUser, opts.MysqlPass, opts.MysqlHost, opts.MysqlPort, opts.MysqlDb, opts.MysqlQry)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// Load config file
	if opts.Config != "" {
		cfg, err := LoadFile(opts.Config, data)
		data["Cfg"] = cfg
		if err != nil {
			log.Fatal(err.Error())
		}

		if cfg["inputs"] == nil {
			log.Fatal("No inputs specified in configuration file")
		}

		if reflect.ValueOf(cfg["inputs"]).Kind() != reflect.Map {
			log.Fatal("Incorrect definition of inputs, it needs to be a map of values")
		}

		var defs CfgDefault
		if cfg["defaults"] != nil {
			defs, err = GetDefaults(cfg["defaults"].(map[string]interface{}))
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		for k, v := range cfg["inputs"].(map[string]interface{}) {
			i, err := GetInput(k, v.(map[string]interface{}), defs)
			if err != nil {
				log.Fatal(err.Error())
			}

			if data[*i.Name] != nil {
				log.Fatalf("Input name already exist's: %s", i.Name)
			}

			switch *i.Type {
			case "file":
				var err error
				data[*i.Name], err = LoadFile(*i.Path, data)
				if err != nil {
					log.Fatal(err.Error())
				}
			case "etcd":
				// Add error handling
				node := []string{fmt.Sprintf("http://%v:%v", i.EtcdNode, i.EtcdPort)}
				client := etcd.NewClient(node)
				res, _ := client.Get(*i.EtcdDir, true, true)
				data[*i.Name] = EtcdMap(res.Node)
			case "http":
				var f DataFmt
				switch *i.HttpFormat {
				case "YAML":
					f = YAML
				case "TOML":
					f = TOML
				case "JSON":
					f = JSON
				default:
					log.Fatal("Unsupported data format, needs to be YAML, JSON or TOML")
				}

				var err error
				data[*i.Name], err = GetHTTP(*i.HttpUrl, *i.HttpHeader, f)
				if err != nil {
					log.Fatal(err.Error())
				}
			case "mysql":
				var err error
				data[*i.Name], err = GetMySQL(*i.MysqlUser, *i.MysqlPass, *i.MysqlHost, *i.MysqlPort, *i.MysqlDb, *i.MysqlQry)
				if err != nil {
					log.Fatal(err.Error())
				}
			default:
				log.Fatalf("Unknown type in configuration file .%v.Type: %v", *i.Name, *i.Type)
			}
		}

		if reflect.ValueOf(cfg["merge"]).Kind() == reflect.Map {
			for k1, v1 := range cfg["merge"].(map[string]interface{}) {
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
					if data[m.Name] == nil {
						data2 := make(map[string]interface{})
						for k, v := range data[m.Inputs[i].(string)].(map[string]interface{}) {
							data2[k] = v
						}
						data[m.Name] = data2
					} else {
						data2 := data[m.Name].(map[string]interface{})
						for k, v := range data[m.Inputs[i].(string)].(map[string]interface{}) {
							data2[k] = v
						}
					}
				}
			}
		}
	}

	// If verbose print data structure as YAML
	if opts.Verbose {
		s, _ := yaml.Marshal(&data)
		log.Printf("Input data\n%s", string(s))
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
	err := t.Execute(buf, data)
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
