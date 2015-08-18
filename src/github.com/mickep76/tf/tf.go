package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"reflect"
	"strconv"

	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/go-etcd/etcd"
	"github.com/flosch/pongo2"
	flags "github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"

	"github.com/mickep76/tf/input"
	"github.com/mickep76/tf/template"
)

// Input options.
type Input struct {
	Name          *string
	Type          *string
	Path          *string
	EtcdHost      *string
	EtcdPort      *int64
	EtcdDir       *string
	HTTPUrl       *string
	HTTPHeader    *string
	HTTPFormat    *string
	MySQLUser     *string
	MySQLPassword *string
	MySQLHost     *string
	MySQLPort     *int64
	MySQLDatabase *string
	MySQLQuery    *string
}

// Merge namespaces.
type Merge struct {
	Name   string
	Inputs []interface{}
}

func main() {
	// Get the FileInfo struct describing the standard input.
	fi, _ := os.Stdin.Stat()

	// Set log options.
	log.SetOutput(os.Stderr)
	log.SetLevel(log.WarnLevel)

	// Options.
	var opts struct {
		Verbose       bool    `short:"v" long:"verbose" description:"Verbose"`
		Version       bool    `long:"version" description:"Version"`
		Config        string  `short:"c" long:"config" description:"YAML, TOML or JSON config file"`
		Input         *string `short:"i" long:"input" description:"Input"`
		InpFormat     string  `short:"F" long:"input-format" description:"Data serialization format YAML, TOML or JSON" default:"YAML"`
		InpFile       *string `short:"f" long:"input-file" description:"Input file, data serialization format used is based on the file extension"`
		TemplFile     *string `short:"t" long:"template" description:"Template file"`
		TemplLang     string  `short:"l" long:"template-lang" description:"Template language text or pongo2" default:"pongo2"`
		OutpFile      *string `short:"o" long:"output" description:"Output file (STDOUT)"`
		Permission    string  `short:"p" long:"permission" description:"File permissions in octal" default:"644"`
		Owner         *string `short:"O" long:"owner" description:"File Owner"`
		EtcdHost      *string `long:"etcd-host" description:"Etcd Host"`
		EtcdPort      int     `long:"etcd-port" description:"Etcd Port" default:"2379"`
		EtcdDir       string  `long:"etcd-dir" description:"Etcd Dir" default:"/"`
		HTTPUrl       *string `long:"http-url" description:"HTTP Url"`
		HTTPHeader    string  `long:"http-header" description:"HTTP Header" default:"Accept: application/json"`
		HTTPFormat    string  `long:"http-format" description:"HTTP Format" default:"JSON"`
		MySQLUser     *string `long:"mysql-user" description:"MySQL user"`
		MySQLPassword *string `long:"mysql-password" description:"MySQL password"`
		MySQLHost     *string `long:"mysql-host" description:"MySQL host"`
		MySQLPort     int64   `long:"mysql-port" description:"MySQL port" default:"3306"`
		MySQLDatabase *string `long:"mysql-database" description:"MySQL database"`
		MySQLQuery    *string `long:"mysql-query" description:"MySQL query"`
	}

	// Parse options.
	if _, err := flags.Parse(&opts); err != nil {
		ferr := err.(*flags.Error)
		if ferr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			log.Fatal(err.Error())
		}
	}

	// Print version.
	if opts.Version {
		fmt.Printf("tf %s\n", Version)
		os.Exit(0)
	}

	// Set verbose.
	if opts.Verbose {
		log.SetLevel(log.InfoLevel)
	}

	// Get environment variables.
	data := make(map[string]interface{})
	data["Env"] = input.GetOSEnv()

	// Get argument input.
	if opts.Input != nil {
		var f input.DataFmt
		switch opts.InpFormat {
		case "YAML":
			f = input.YAML
		case "TOML":
			f = input.TOML
		case "JSON":
			f = input.JSON
		default:
			log.Fatal("Unsupported data format, needs to be YAML, JSON or TOML")
		}

		var err error
		data["Arg"], err = input.UnmarshalData([]byte(*opts.Input), f)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Copy .Arg namespace to . for conveniencea.
		for k, v := range data["Arg"].(map[string]interface{}) {
			data[k] = v
		}
	}

	// Get file input.
	if opts.InpFile != nil {
		var err error
		data["File"], err = input.LoadFile(*opts.InpFile, data)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// Get Etcd input.
	if opts.EtcdHost != nil {
		node := []string{fmt.Sprintf("http://%v:%v", *opts.EtcdHost, opts.EtcdPort)}
		client := etcd.NewClient(node)
		res, err := client.Get(opts.EtcdDir, true, true)
		if err != nil {
			log.Fatal(err.Error())
		}
		data["Etcd"] = input.EtcdMap(res.Node)
	}

	// Get http input.
	if opts.HTTPUrl != nil {
		var f input.DataFmt
		switch opts.HTTPFormat {
		case "YAML":
			f = input.YAML
		case "TOML":
			f = input.TOML
		case "JSON":
			f = input.JSON
		default:
			log.Fatal("Unsupported data format, needs to be YAML, JSON or TOML")
		}

		var err error
		data["HTTP"], err = input.GetHTTP(*opts.HTTPUrl, opts.HTTPHeader, f)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// Get MySQL input.
	if opts.MySQLHost != nil {
		var err error

		if opts.MySQLUser == nil {
			log.Fatal("For input \"--mysql-host\" you need to specify \"--mysql-user\"")
		}
		if opts.MySQLPassword == nil {
			log.Fatal("For input \"--mysql-host\" you need to specify \"--mysql-password\"")
		}
		if opts.MySQLDatabase == nil {
			log.Fatal("For input \"--mysql-host\" you need to specify \"--mysql-database\"")
		}
		if opts.MySQLQuery == nil {
			log.Fatal("For input \"--mysql-host\" you need to specify \"--mysql-query\"")
		}

		data["MySQL"], err = input.GetMySQL(*opts.MySQLUser, *opts.MySQLPassword, *opts.MySQLHost, opts.MySQLPort, *opts.MySQLDatabase, *opts.MySQLQuery)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// Load config file.
	if opts.Config != "" {
		cfg, err := input.LoadFile(opts.Config, data)
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
				log.Fatalf("Input name already exist's: %s", *i.Name)
			}

			switch *i.Type {
			case "file":
				var err error
				data[*i.Name], err = input.LoadFile(*i.Path, data)
				if err != nil {
					log.Fatal(err.Error())
				}
			case "etcd":
				node := []string{fmt.Sprintf("http://%v:%v", *i.EtcdHost, *i.EtcdPort)}
				client := etcd.NewClient(node)
				res, err := client.Get(*i.EtcdDir, true, true)
				if err != nil {
					log.Fatal(err.Error())
				}
				data[*i.Name] = input.EtcdMap(res.Node)
			case "http":
				var f input.DataFmt
				switch *i.HTTPFormat {
				case "YAML":
					f = input.YAML
				case "TOML":
					f = input.TOML
				case "JSON":
					f = input.JSON
				default:
					log.Fatal("Unsupported data format, needs to be YAML, JSON or TOML")
				}

				var err error
				data[*i.Name], err = input.GetHTTP(*i.HTTPUrl, *i.HTTPHeader, f)
				if err != nil {
					log.Fatal(err.Error())
				}
			case "mysql":
				var err error
				data[*i.Name], err = input.GetMySQL(*i.MySQLUser, *i.MySQLPassword, *i.MySQLHost, *i.MySQLPort, *i.MySQLDatabase, *i.MySQLQuery)
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

	// If verbose print data structure as YAML.
	if opts.Verbose {
		s, _ := yaml.Marshal(&data)
		log.Printf("Input data\n%s", string(s))
	}

	// Template input.
	var templ string
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		bytes, _ := ioutil.ReadAll(os.Stdin)
		templ = string(bytes)
	} else if opts.TemplFile != nil {
		if _, err := os.Stat(*opts.TemplFile); os.IsNotExist(err) {
			log.Printf("File doesn't exist: %v\n", opts.TemplFile)
			os.Exit(1)
		}

		// Open file.
		c, err := ioutil.ReadFile(*opts.TemplFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		templ = string(c)
	} else {
		log.Printf("No template specified using --template-file (-t) or piped to STDIN\n")
		os.Exit(1)
	}

	var outp string
	if opts.TemplLang == "text" {
		// Parse and compile template.
		buf, err := template.Compile(templ, data)
		if err != nil {
			log.Fatal(err.Error())
		}
		outp = buf.String()
	} else {
		var err error
		// Compile template
		tpl, err := pongo2.FromString(templ)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Render template.
		outp, err = tpl.Execute(data) //pongo2.Context{"name": "florian"})
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// Write result
	if opts.OutpFile != nil {
		p, err := strconv.ParseUint(opts.Permission, 8, 32)
		if err != nil {
			log.Fatal(err.Error())
		}

		w, err := os.Create(*opts.OutpFile)
		if err != nil {
			log.Fatal(err.Error())
		}

		w.Chmod(os.FileMode(p))

		if opts.Owner != nil {
			u, err := user.Lookup(*opts.Owner)
			if err != nil {
				log.Fatal(err.Error())
			}

			uid, err := strconv.Atoi(u.Uid)
			if err != nil {
				log.Fatal(err.Error())
			}

			gid, err := strconv.Atoi(u.Gid)
			if err != nil {
				log.Fatal(err.Error())
			}

			err = w.Chown(uid, gid)
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		_, err = w.Write(([]byte)(outp))
		if err != nil {
			log.Fatal(err.Error())
		}

		w.Close()
	} else {
		fmt.Printf("%v\n", outp)
	}
}
