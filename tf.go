package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	etcd "github.com/coreos/go-etcd/etcd"
	flags "github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/user"
	"path/filepath"
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

// Data format represents which data serialization is used YAML, JSON or TOML.
type DataFmt int

// Constants for data format.
const (
	YAML DataFmt = iota
	TOML
	JSON
)

// Unmarshal YAML/JSON/TOML serialized data.
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
		return nil, errors.New("Unsupported data format")
	}

	return v, nil
}

// Load file with serialized data.
func LoadFile(fn string) (map[string]interface{}, error) {
	var fmt DataFmt

	switch filepath.Ext(fn) {
	case ".yaml":
		fmt = YAML
	case ".json":
		fmt = JSON
	case ".toml":
		fmt = TOML
	default:
		log.Error("Unsupported data format, needs to be .yaml, .json or .toml")
		return nil, errors.New("Unsupported data format")
	}

	_, err := os.Stat(fn)
	if os.IsNotExist(err) {
		log.Errorf("File doesn't exist: %s", fn)
		return nil, err
	}

	log.Infof("Reading file: %s", fn)
	c, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Errorf("Failed to read file: %s", fn)
		return nil, err
	}

	v, err2 := UnmarshalData(c, fmt)
	if err2 != nil {
		return nil, err2
	}

	return v, nil
}

// Get OS Environment variables.
func GetOSEnv() map[string]interface{} {
	v := make(map[string]interface{})

	for _, e := range os.Environ() {
		a := strings.Split(e, "=")
		v[a[0]] = a[1]
	}

	return v
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
		y, err = UnmarshalData([]byte(opts.Input), fmt)
		check(err)
	} else if opts.InpFile != "" {
		y, err = LoadFile(opts.InpFile)
		check(err)
	} else {
		y = make(map[string]interface{})
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
			for key, val := range c["inputs"].(map[string]interface{}) {
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
