package main

import (
	etcd "github.com/coreos/go-etcd/etcd"
	"strings"
)

// Create a nested data structure from Etcd node.
func EtcdMap(root *etcd.Node) map[string]interface{} {
    v := make(map[string]interface{})

	for _, n := range root.Nodes {
		keys := strings.Split(n.Key, "/")
		k := keys[len(keys)-1]
		if n.Dir {
			v[k] = make(map[string]interface{})
			v[k] = EtcdMap(n)
		} else {
			v[k] = n.Value
		}
	}
	return v
}
