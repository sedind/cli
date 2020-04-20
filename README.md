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

## Getting Started

One of the philosophies behind cli is that an API should be playful and full of
discovery. So a cli app can be as little as one line of code in `main()`.

``` go
package main

import (
  "os"

  "github.com/sedind/cli"
)

func main() {
  cli.New().Run(os.Args)
}
```

This app will run and show help text, but is not very useful. Let's give an
action to execute and some help documentation:


``` go
package main

import (
  "fmt"
  "log"
  "os"

  "github.com/sedind/cli"
)

func main() {
  app := cli.New()
  app.Name = "greet"
  app.Description =  "creates greeting message"
  app.Action =  func(c *cli.Context) error {
    fmt.Println("Howdy!")
    return nil
  }
  

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
```

Install our command to the `$GOPATH/bin` directory:

```
$ go install
```

Finally run our new command:

```
$ greet
Howdy!
```

cli also generates neat help text:

```
$ greet help
creates greeting message

Usage:	build [flags] [args...]
	-help	    h	Prints application help (false)
	-version	v	Prints application version (false)

```

### GOPATH

Make sure your `PATH` includes the `$GOPATH/bin` directory so your commands can
be easily used:
```
export PATH=$PATH:$GOPATH/bin
```

### Supported platforms

cli is tested against multiple versions of Go on Linux, and against the latest
released version of Go on OS X and Windows. 
