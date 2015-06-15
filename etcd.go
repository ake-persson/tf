package main

import (
	etcd "github.com/coreos/go-etcd/etcd"
	"strings"
)

// Return Etcd structure as nested map[interface{}]interface{}
func etcdNestedMap(node *etcd.Node, data map[string]interface{}) {
	for _, node := range node.Nodes {
		keys := strings.Split(node.Key, "/")
		key := keys[len(keys)-1]
		if node.Dir {
			data[key] = make(map[string]interface{})
			etcdNestedMap(node, data[key].(map[string]interface{}))
		} else {
			data[key] = node.Value
		}
	}
}
