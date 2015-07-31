#!/bin/bash

set -eux

if [ -n "${DOCKER_HOST}" ]; then
    IP=$(echo ${DOCKER_HOST##tcp://} | awk -F: '{ print $1}')
else
    IP='localhost'
fi

docker run -d --name etcd -p 4001:4001 coreos/etcd:v0.4.6

curl -L http://${IP}:4001/v2/keys/hosts/host1.example.com/serialno -XPUT -d value="abc123"
curl -L http://${IP}:4001/v2/keys/hosts/host2.example.com/serialno -XPUT -d value="def456"
curl -L http://${IP}:4001/v2/keys/hosts/host3.example.com/serialno -XPUT -d value="fgh789"

tf --config examples/etcd/tf.toml --template examples/etcd/hosts.tf --input "{ EtcdHost: ${IP} }"

docker stop etcd
docker rm etcd
