package input

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"

	"github.com/mickep76/tf/template"
)

// DataFmt represents which data serialization is used YAML, JSON or TOML.
type DataFmt int

// Constants for data format.
const (
	YAML DataFmt = iota
	TOML
	JSON
)

// UnmarshalData unmarshal YAML/JSON/TOML serialized data.
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

// LoadFile loads a file with serialized data.
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
	buf, err := template.Compile(string(c), data)
	if err != nil {
		log.Fatal(err.Error())
	}

	/*
		buf := new(bytes.Buffer)
		err = t.Execute(buf, data)
		if err != nil {
			return nil, err
		}
	*/
	log.Infof("Input file result: %s\n%s", fn, string(buf.Bytes()))

	v, err2 := UnmarshalData(buf.Bytes(), f)
	if err2 != nil {
		return nil, err2
	}

	return v, nil
}

// GetOSEnv gets OS Environment variables.
func GetOSEnv() map[string]interface{} {
	v := make(map[string]interface{})

	for _, e := range os.Environ() {
		a := strings.Split(e, "=")
		v[a[0]] = a[1]
	}

	return v
}

// GetHTTP requests a HTTP url.
func GetHTTP(url string, header string, f DataFmt) (map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	h := strings.Split(header, ":")

	req.Header.Add(h[0], h[1])

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	body, err2 := ioutil.ReadAll(r.Body)
	if err2 != nil {
		return nil, err2
	}

	v, err := UnmarshalData(body, f)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// GetMySQL queries MySQL.
func GetMySQL(user string, pass string, host string, port int64, db string, qry string) ([]interface{}, error) {
	log.Infof("Connecting to MySQL to database %s on host %s", host, db)
	log.Infof("Connect DSN: %s:%s@tcp(%s:%v)/%s", user, "xxxxxxxx", host, port, db)
	dbo, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", user, pass, host, port, db))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer dbo.Close()

	err = dbo.Ping()
	if err != nil {
		return nil, err
	}

	log.Infof("Execute query: %s", qry)
	rows, err := dbo.Query(qry)
	if err != nil {
		return nil, err
	}

	log.Infof("Get result from query")
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var data []interface{}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		var value string
		res := make(map[string]interface{})
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			res[columns[i]] = value
		}
		data = append(data, res)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}
