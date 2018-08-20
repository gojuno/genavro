# genavro

[![License MIT](https://img.shields.io/badge/License-MIT-blue.svg)](http://opensource.org/licenses/MIT) [![GoDoc](https://godoc.org/github.com/gojuno/genavro?status.svg)](http://godoc.org/github.com/gojuno/genavro) [![Go Report Card](https://goreportcard.com/badge/github.com/gojuno/genavro)](https://goreportcard.com/report/github.com/gojuno/genavro) [![Build Status](https://travis-ci.org/gojuno/genavro.svg?branch=master)](http://travis-ci.org/gojuno/genavro)

Generates avrò protocols from golang structs.

#### Get the package using:

```
$ go get -u -v github.com/gojuno/astparser
```

#### Usage

Build tool with help of make file
```bash
make build
```

And generate avro protocols
```bash
bin/genavro -in <go_structs_dir> -o <output_dir> -n <avro_protocol_namespace> 
```

There are two additional flags:
 
 * `e` expect regexp to exclude some files from the passed dir. 
 * `i` expects regexps include only specific files from passed dir.