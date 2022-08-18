# Properties

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gookit/properties?style=flat-square)
[![Unit-Tests](https://github.com/gookit/properties/actions/workflows/go.yml/badge.svg)](https://github.com/gookit/properties/actions/workflows/go.yml)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/gookit/properties)](https://github.com/gookit/properties)
[![GoDoc](https://godoc.org/github.com/gookit/properties?status.svg)](https://pkg.go.dev/github.com/gookit/properties/v3)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/properties)](https://goreportcard.com/report/github.com/gookit/properties)

`properties` - Java Properties format contents parse, marshal and unmarshal library.

- Generic properties contents parser, marshal and unmarshal
- Support `Marshal` and `Unmarshal` like `json` package
- Support comments start withs `!`, `#`
    - enhanced: allow `//`, `/* multi line comments */`
- Support multi line string value, withs `\\`
    - enhanced: allow `'''multi line string''''`, `"""multi line string"""`
- Support value refer parse by var. format: `${some.other.key}`
- Support ENV var parse. format: `${APP_ENV}`, `${APP_ENV | default}`

> **[EN README](README.md)**

## 安装

```shell
go get github.com/gookit/properties
```

## 使用

```go
package main

import (
	"fmt"

	"github.com/gookit/properties"
)

func Example() {
	text := `
# properties string
name = inhere
age = 200
`

	p, err := properties.Parse(text)
	if err != nil {
		panic(err)
	}

	type MyConf struct {
		Name string `properties:"name"`
		Age  int    `properties:"age"`
	}

	cfg := &MyConf{}
	err = p.MapStruct("", cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println(*cfg)

	// Output:
	// {inhere 200}
}
```

## Marshal/Unmarshal

- `Marshal(v interface{}) ([]byte, error)`
- `Unmarshal(v []byte, ptr interface{}) error`

## 配置管理

If you want more support for file formats and multi file loads, recommended use [gookit/config](https://github.com/gookit/config)

- Support multi formats: `JSON`(default), `INI`, `Properties`, `YAML`, `TOML`, `HCL`
- Support multi file loads, will auto merge loaded data

## Gookit packages

- [gookit/ini](https://github.com/gookit/ini) INI parse by golang. INI config data management library.
- [gookit/rux](https://github.com/gookit/rux) Simple and fast request router, web framework for golang
- [gookit/gcli](https://github.com/gookit/gcli) build CLI application, tool library, running CLI commands
- [gookit/event](https://github.com/gookit/event) Lightweight event manager and dispatcher implements by Go
- [gookit/cache](https://github.com/gookit/cache) Generic cache use and cache manager for golang. support File, Memory, Redis, Memcached.
- [gookit/config](https://github.com/gookit/config) Go config management. support JSON, YAML, TOML, INI, HCL, ENV and Flags
- [gookit/color](https://github.com/gookit/color) A command-line color library with true color support, universal API methods and Windows support
- [gookit/filter](https://github.com/gookit/filter) Provide filtering, sanitizing, and conversion of golang data
- [gookit/validate](https://github.com/gookit/validate) Use for data validation and filtering. support Map, Struct, Form data
- [gookit/goutil](https://github.com/gookit/goutil) Some utils for the Go: string, array/slice, map, format, cli, env, filesystem, test and more
- More, please see https://github.com/gookit

## See also

- Ini parse [gookit/ini/parser](https://github.com/gookit/ini/tree/master/parser)
- Properties parse [gookit/properties](https://github.com/gookit/properties)

## License

**MIT**
