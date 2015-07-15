package main

import (
	"errors"
	"fmt"
)

type CfgDefault struct {
	EtcdNode   *string
	EtcdPort   *int64
	HttpHeader *string
	HttpFormat *string
	MysqlUser  *string
	MysqlPass  *string
	MysqlHost  *string
	MysqlPort  *int64
	MysqlDb    *string
}

type CfgInput struct {
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

// Get defaults from config file.
func GetDefaults(defs map[string]interface{}) (CfgDefault, error) {
	var d CfgDefault
	for k, v := range defs {
		switch k {
		case "etcd_node":
			s := v.(string)
			d.EtcdNode = &s
		case "etcd_port":
			n := v.(int64)
			d.EtcdPort = &n
		case "http_header":
			s := v.(string)
			d.HttpHeader = &s
		case "http_format":
			s := v.(string)
			d.HttpFormat = &s
		case "mysql_user":
			s := v.(string)
			d.MysqlUser = &s
		case "mysql_pass":
			s := v.(string)
			d.MysqlPass = &s
		case "mysql_host":
			s := v.(string)
			d.MysqlHost = &s
		case "mysql_port":
			n := v.(int64)
			d.MysqlPort = &n
		case "mysql_db":
			s := v.(string)
			d.MysqlDb = &s
		default:
			return CfgDefault{}, errors.New(fmt.Sprintf("Invalid configuration key \"%v\" in [defaults]", k))
		}
	}
	return d, nil
}

// Get input from configuration file.
func GetInput(name string, inp map[string]interface{}, d CfgDefault) (CfgInput, error) {
	var i CfgInput

	if d.EtcdNode != nil {
		i.EtcdNode = d.EtcdNode
	}
	if d.EtcdPort != nil {
		i.EtcdPort = d.EtcdPort
	} else {
		n := int64(2379)
		i.EtcdPort = &n
	}
	if d.HttpHeader != nil {
		i.HttpHeader = d.HttpHeader
	} else {
		s := "Accept: application/json"
		i.HttpHeader = &s
	}
	if d.HttpFormat != nil {
		i.HttpFormat = d.HttpFormat
	} else {
		s := "JSON"
		i.HttpFormat = &s
	}
	if d.MysqlUser != nil {
		i.MysqlUser = d.MysqlUser
	}
	if d.MysqlPass != nil {
		i.MysqlPass = d.MysqlPass
	}
	if d.MysqlHost != nil {
		i.MysqlHost = d.MysqlHost
	}
	if d.MysqlPort != nil {
		i.MysqlPort = d.MysqlPort
	} else {
		n := int64(3306)
		i.MysqlPort = &n
	}
	if d.MysqlDb != nil {
		i.MysqlDb = d.MysqlDb
	}

	i.Name = &name
	for k, v := range inp {
		switch k {
		case "name":
			s := v.(string)
			i.Name = &s
		case "type":
			s := v.(string)
			i.Type = &s
		case "path":
			s := v.(string)
			i.Path = &s
		case "etcd_node":
			s := v.(string)
			i.EtcdNode = &s
		case "etcd_port":
			n := v.(int64)
			i.EtcdPort = &n
		case "etcd_dir":
			s := v.(string)
			i.EtcdDir = &s
		case "http_url":
			s := v.(string)
			i.HttpUrl = &s
		case "http_header":
			s := v.(string)
			i.HttpHeader = &s
		case "http_format":
			s := v.(string)
			i.HttpFormat = &s
		case "mysql_user":
			s := v.(string)
			i.MysqlUser = &s
		case "mysql_pass":
			s := v.(string)
			i.MysqlPass = &s
		case "mysql_host":
			s := v.(string)
			i.MysqlHost = &s
		case "mysql_port":
			n := v.(int64)
			i.MysqlPort = &n
		case "mysql_db":
			s := v.(string)
			i.MysqlDb = &s
		case "mysql_qry":
			s := v.(string)
			i.MysqlQry = &s
		default:
			return CfgInput{}, errors.New(fmt.Sprintf("Invalid configuration key \"%v\" in [inputs.%v]", k, name))
		}
	}

	switch *i.Type {
	case "file":
		if i.Path == nil {
			return CfgInput{}, errors.New(fmt.Sprintf("For input [inputs.%v] type \"file\" you need to specify \"path\"", name))
		}
	case "etcd":
		if i.EtcdNode == nil {
			return CfgInput{}, errors.New(fmt.Sprintf("For input [inputs.%v] type \"etcd\" you need to specify \"etcd_node\"", name))
		}
		if i.EtcdPort == nil {
			return CfgInput{}, errors.New(fmt.Sprintf("For input [inputs.%v] type \"etcd\" you need to specify \"etcd_port\"", name))
		}
	case "http":
		if i.HttpUrl == nil {
			return CfgInput{}, errors.New(fmt.Sprintf("For input [inputs.%v] type \"http\" you need to specify \"http_url\"", name))
		}
		if i.HttpHeader == nil {
			return CfgInput{}, errors.New(fmt.Sprintf("For input [inputs.%v] type \"http\" you need to specify \"http_header\"", name))
		}
		if i.HttpFormat == nil {
			return CfgInput{}, errors.New(fmt.Sprintf("For input [inputs.%v] type \"http\" you need to specify \"http_format\"", name))
		}
	case "mysql":
		if i.MysqlUser == nil {
			return CfgInput{}, errors.New(fmt.Sprintf("For input [inputs.%v] type \"mysql\" you need to specify \"mysql_user\"", name))
		}
		if i.MysqlPass == nil {
			return CfgInput{}, errors.New(fmt.Sprintf("For input [inputs.%v] type \"mysql\" you need to specify \"mysql_pass\"", name))
		}
		if i.MysqlHost == nil {
			return CfgInput{}, errors.New(fmt.Sprintf("For input [inputs.%v] type \"mysql\" you need to specify \"mysql_host\"", name))
		}
		if i.MysqlPort == nil {
			return CfgInput{}, errors.New(fmt.Sprintf("For input [inputs.%v] type \"mysql\" you need to specify \"mysql_port\"", name))
		}
		if i.MysqlDb == nil {
			return CfgInput{}, errors.New(fmt.Sprintf("For input [inputs.%v] type \"mysql\" you need to specify \"mysql_db\"", name))
		}
		if i.MysqlQry == nil {
			return CfgInput{}, errors.New(fmt.Sprintf("For input [inputs.%v] type \"mysql\" you need to specify \"mysql_qry\"", name))
		}
	default:
		return CfgInput{}, errors.New(fmt.Sprintf("Unknown type \"%v\" for input [inputs.%v]", *i.Type, *i.Name))
	}

	return i, nil
}
