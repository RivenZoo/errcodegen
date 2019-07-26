# errcodegen
A error code generator.

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Build Status](https://travis-ci.org/RivenZoo/errcodegen.svg?branch=master)](https://travis-ci.org/RivenZoo/errcodegen)

## Install

Install binary and put it under your `PATH`.

```shell
$ go install github.com/RivenZoo/errcodegen
```

## Usage

```
Usage:
  errcodegen [flags]

Flags:
      --appcode string        App error code, three digits (default "100")      # set code for your app
      --config string         config file (default is $HOME/.errcodegen.yaml)   # config file path
      --err_def string        Error code definiton file                         # definition file path
      --err_func string       New error function name (default "New")           # function to new error variable
      --err_func_pkg string   New error function package import path (default "errors") # package that err_func belongs to
      --pkg string            Generated module package name (default "errorcodes")  # generated code package name
```

### Use docker

```
$ cd example
$ docker run --rm -it -v $(pwd):/data/app rivenzoe/errcodegen errcodegen --err_def=error_code_def.conf --err_func=Errorf --err_func_pkg=fmt --pkg=example
```

## Example

See [example](https://github.com/RivenZoo/errcodegen/tree/master/example)
