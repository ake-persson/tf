package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	etcd "github.com/coreos/go-etcd/etcd"
	flags "github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
	"io/ioutil"
//	"log"
	log "github.com/Sirupsen/logrus"
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

func map_print(y map[string]interface{}, dir string, pad string) {
	for key, val := range y {
		if reflect.ValueOf(val).Kind() == reflect.Map {
			fmt.Printf("%v[%v]\n", pad, key)
			map_print(val.(map[string]interface{}), dir+"/"+key, pad+"    ")
		} else {
			fmt.Printf("%v%v: %v (%v)\n", pad, key, val, dir)
		}
	}
}

// Data format type and constants
type DataFmt int

const (
	YAML DataFmt = iota
	TOML
	JSON
)

// Unmarshal YAML/JSON/TOML data
func UnmarshalData(cont []byte, fmt DataFmt) (map[string]interface{}, error) {
    v := make(map[string]interface{})

	switch fmt {
	case YAML:
		log.Info("Unmarshaling YAML data")
		err := yaml.Unmarshal(cont, &v)
		if err != nil {
			return nil, err
		}
	case TOML:
		log.Info("Unmarshaling TOML data")
		err := toml.Unmarshal(cont, &v)
		if err != nil {
			return nil, err
		}
	case JSON:
		log.Info("Unmarshaling JSON data")
		err := json.Unmarshal(cont, &v)
		if err != nil {
			return nil, err
		}
	default:
		log.Error("Unsupported data format")
		// Placeholder: return error
	}

	return v, nil
}

// Load YAML/JSON/TOML file
func LoadFile(fn string, fmt DataFmt) (map[string]interface{}, error) {
	_, err := os.Stat(fn)
	if os.IsNotExist(err) {
		log.Errorf("File doesn't exist: %s", fn)
		return nil, err
	}

	log.Infof("Reading file: %s", fn)
	cont, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Errorf("Failed to read file: %s", fn)
		return nil, err
	}

	v, err2 := UnmarshalData(cont, fmt)
	if err2 != nil {
		return nil, err2
	}

	return v, nil
}

// Get OS Environment variables
func GetOSEnv() map[string]interface{} {
    r := make(map[string]interface{})

    for _, e := range os.Environ() {
        v := strings.Split(e, "=")
        r[v[0]] = v[1]
    }

	return r
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
		Config     string `short:"c" long:"config" description:"TOML Config file"`
		Input      string `short:"i" long:"input" description:"YAML input"`
		InpFile    string `short:"f" long:"input-file" description:"YAML input file"`
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
	check(err)

	// Print version
	if opts.Version == true {
		fmt.Println(version)
		os.Exit(0)
	}

	// Get input
	var y map[string]interface{}
	if opts.Input != "" {
		if opts.InpFile != "" {
			log.Fatal("Can't specify both --input (-i) and --input-file (-f)\n")
		}
		y, err = UnmarshalData([]byte(opts.Input), YAML)
		check(err)
	} else if opts.InpFile != "" {
		y, err = LoadFile(opts.InpFile, YAML)
		check(err)
	} else {
 	   y = make(map[string]interface{})
	}

	// Get environment
	env := GetOSEnv()
	y["Env"] = env

	// Load config file
	if opts.Config != "" {
		cfg, err := ioutil.ReadFile(opts.Config)
		check(err)

		var t map[string]interface{}
		err = toml.Unmarshal(cfg, &t)
		check(err)
		y["Cfg"] = t

		if reflect.ValueOf(t["inputs"]).Kind() == reflect.Map {
			var i map[string]interface{}
			i = t["inputs"].(map[string]interface{})

			for key, val := range i {
				fmt.Printf("%v, %v\n", key, val)
			}
		}
	}

	vars := make(map[string]interface{})
	if opts.EtcdNode != "" {
		node := []string{fmt.Sprintf("http://%v:%v", opts.EtcdNode, opts.EtcdPort)}
		client := etcd.NewClient(node)
		res, _ := client.Get(opts.EtcdKey, true, true)
		etcdNestedMap(res.Node, vars)
		y["Etcd"] = vars
	}

	// s, err := yaml.Marshal(&y)
	// fmt.Printf("%s\n", string(s))

	if opts.Verbose {
		map_print(y, "", "")
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
