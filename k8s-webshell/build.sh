#!/usr/bin/env bash
GOOS=linux go build -o bin/k8s-webshell-linux
GOOS=darwin go build -o bin/k8s-webshell-mac
GOOS=windows go build -o bin/k8s-webshell.exe