#!/bin/bash

checkcmd="git status --porcelain"
printcmd="git status --short"

if [[ "$#" -gt 0 ]]; then
    checkcmd="$checkcmd -- $@"
    printcmd="$printcmd -- $@"
fi

if [[ ! -z "$($checkcmd)" ]]; then
    echo "Working directory is dirty. The following modifications were found:"
    $printcmd
    exit 1
fi
