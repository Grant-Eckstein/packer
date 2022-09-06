# Go Packer Command Line Tool[![Go Report Card](https://goreportcard.com/badge/github.com/Grant-Eckstein/packer)](https://goreportcard.com/report/github.com/Grant-Eckstein/packer) [![GoDoc](https://godoc.org/github.com/Grant-Eckstein/packer?status.svg)](https://godoc.org/github.com/Grant-Eckstein/packer)
*An executable packing tool*

## Overview
Packer is my take on obfustating an existing executable. Since Packer is written in Golang, it functions on all major systems, assuming that your payload also supports the target system. 

## Example usage
This assumes you are building for a Mac. 
```bash
$ git clone https://github.com/Grant-Eckstein/packer.git
$ cd packer
$ go get ...
$ go build
$ ./packer pack INSERT_EXE_HERE -o darwin -a amd64
$ ./packed
```