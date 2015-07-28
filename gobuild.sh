#!/bin/bash

set -eux

go get github.com/constabulary/gb/...
${GOPATH}/bin/gb build all
