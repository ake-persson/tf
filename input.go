package main

import (
	"bytes"
	"encoding/json"
	"errors"
	//	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Data format represents which data serialization is used YAML, JSON or TOML.
type DataFmt int

// Constants for data format.
const (
	YAML DataFmt = iota
	TOML
	JSON
)

// Unmarshal YAML/JSON/TOML serialized data.
func UnmarshalData(cont []byte, f DataFmt) (map[string]interface{}, error) {
	v := make(map[string]interface{})

	switch f {
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
func LoadFile(fn string, data map[string]interface{}) (map[string]interface{}, error) {
	var f DataFmt

	switch filepath.Ext(fn) {
	case ".yaml":
		f = YAML
	case ".json":
		f = JSON
	case ".toml":
		f = TOML
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

	log.Infof("Template input file: %s", fn)
	t := template.Must(template.New("template").Funcs(fns).Parse(string(c)))

	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	log.Infof("Input file result: %s\n%s", fn, string(buf.Bytes()))

	v, err2 := UnmarshalData(buf.Bytes(), f)
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
