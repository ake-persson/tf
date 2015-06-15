package main

import (
	etcd "github.com/coreos/go-etcd/etcd"
	"strings"
)

// Return Etcd structure as nested map[interface{}]interface{}
func etcdNestedMap(node *etcd.Node, data map[interface{}]interface{}) {
	for _, node := range node.Nodes {
		keys := strings.Split(node.Key, "/")
		key := keys[len(keys)-1]
		if node.Dir {
			data[key] = make(map[interface{}]interface{})
			etcdNestedMap(node, data[key].(map[interface{}]interface{}))
		} else {
			data[key] = node.Value
		}
	}
}
