# Properties

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gookit/properties?style=flat-square)
[![Actions Status](https://github.com/gookit/properties/workflows/action-tests/badge.svg)](https://github.com/gookit/properties/actions)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/gookit/properties)](https://github.com/gookit/properties)
[![GoDoc](https://godoc.org/github.com/gookit/properties?status.svg)](https://pkg.go.dev/github.com/gookit/properties/v3)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/properties)](https://goreportcard.com/report/github.com/gookit/properties)

Java `properties` contents parse, marshal and unmarshal library.

- Generic properties contents parser
- Support comments withs `#`, `//`, `/* multi line comments */`
- Support multi line string value, withs `'''multi line string''''`, `"""multi line string"""`
- Support ENV var parse. format: `{$APP_ENV}`, `{$APP_ENV | default}`
- Support value refer parse by var. format: `{$key_name}`

> **[EN README](README.md)**

## Install

```shell
go get github.com/gookit/properties
```

## Usage

```go
// ...
```

