#!/bin/sh

pwd
GOOS=linux go build main.go
zip function.zip main  
