# Properties

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gookit/properties?style=flat-square)
[![Unit-Tests](https://github.com/gookit/properties/actions/workflows/go.yml/badge.svg)](https://github.com/gookit/properties/actions/workflows/go.yml)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/gookit/properties)](https://github.com/gookit/properties)
[![GoDoc](https://godoc.org/github.com/gookit/properties?status.svg)](https://pkg.go.dev/github.com/gookit/properties/v3)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/properties)](https://goreportcard.com/report/github.com/gookit/properties)

`properties` - Java Properties format contents parse, marshal and unmarshal library.

- Generic properties contents parser, marshal and unmarshal
- Support `Marshal` and `Unmarshal` like `json` package
- Support comments start withs `!`, `#`, `//`, `/* multi line comments */`
- Support multi line string value, withs `\\`, `'''multi line string''''`, `"""multi line string"""`
- Support value refer parse by var. format: `${some.other.key}`
- Support ENV var parse. format: `{$APP_ENV}`, `{$APP_ENV | default}`

> **[中文说明](README.zh-CN.md)**

## Install

```shell
go get github.com/gookit/properties
```

## Usage

```go
// ...
```

