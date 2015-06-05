#!/bin/bash

set -eu
set -x

TMPDIR="$(pwd)/.gobuild"
export GOPATH=${TMPDIR}
SRCDIR="${GOPATH}/src/github.com/mickep76/tf"

[ -d ${TMPDIR} ] && rm -rf ${TMPDIR}
mkdir -p ${GOPATH}/{src,pkg,bin}
mkdir -p ${SRCDIR}
cp tf.go ${SRCDIR}
(
    echo ${GOPATH}
    cd ${SRCDIR}
    go get .
    go install .
)
