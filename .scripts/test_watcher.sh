#!/bin/bash

dogo() {
       go fmt
       go generate
       gotest -v ./...
}



dogo

inotifywait --exclude "[^g].$|[^o]$" -m -r -e close_write ./ |
    while read path action file; do
           dogo
    done
