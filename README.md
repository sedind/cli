cli
===

[![GoDoc](https://godoc.org/github.com/sedind/cli?status.svg)](https://godoc.org/github.com/sedind/cli)
[![Go Report Card](https://goreportcard.com/badge/sedind/cli)](https://goreportcard.com/report/sedind/cli)

cli is a simple package for building command line apps in Go. The goal is to enable developers to write fast and distributable command line applications in an expressive way.

## Installation

Using this package requires a working Go environment. [See the install instructions for Go](http://golang.org/doc/install.html).

Go Modules are required when using this package. [See the go blog guide on using Go Modules](https://blog.golang.org/using-go-modules).


```
$ GO111MODULE=on go get github.com/sedind/cli
```

```go
...
import (
  "github.com/sedind/cli"
)
...
```

## Usage Documentation

### GOPATH

Make sure your `PATH` includes the `$GOPATH/bin` directory so your commands can
be easily used:
```
export PATH=$PATH:$GOPATH/bin
```

### Supported platforms

cli is tested against multiple versions of Go on Linux, and against the latest
released version of Go on OS X and Windows. 