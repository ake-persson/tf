#!/bin/bash

error()
{
    echo $1 >&2
    exit 1
}

[ -d .git ] || error "You need to be in the Git repo. root to use this script."
[ -d .githooks ] || error "There is no .githooks directory."
[ -d .git/hooks ] && rm -rf .git/hooks
(
    cd .git
    ln -s ../.githooks hooks
)
