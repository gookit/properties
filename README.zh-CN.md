# Properties

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gookit/properties?style=flat-square)
[![Unit-Tests](https://github.com/gookit/properties/actions/workflows/go.yml/badge.svg)](https://github.com/gookit/properties/actions/workflows/go.yml)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/gookit/properties)](https://github.com/gookit/properties)
[![GoDoc](https://godoc.org/github.com/gookit/properties?status.svg)](https://pkg.go.dev/github.com/gookit/properties/v3)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/properties)](https://goreportcard.com/report/github.com/gookit/properties)

`properties` - Java Properties format contents parse, marshal and unmarshal library.

- 通用的Java `properties` 格式内容解析器
- 支持像 `json` 包一样的 `Marshal` 和 `Unmarshal` 
- 支持以 `!`, `#` 开头的行注释
    - 增强: 也额外支持 `//`, `/* multi line comments */`
- 支持多行字符串值，以 `\\` 结尾进行换行
    - 增强: 也额外支持 `'''multi line string''''`, `"""multi line string"""`
- 支持值引用 var 解析。 format: `${some.other.key}`
- 支持 ENV 变量解析。 format: `${APP_ENV}`, `${APP_ENV | default}`

> **[EN README](README.md)**

## 安装

```shell
go get github.com/gookit/properties
```

## 使用

`properties` 格式内容示例:

```properties
name = inhere
age = 345
only-key = 
env-key = ${SHELL | bash}

 ##### comments1
top.sub.key0 = a string value
top.sub.key1 = "a quote value1"
top.sub.key2 = 'a quote value2'
/* comments 1.1 */
top.sub.key3 = 234

! inline list
top2.inline.list.ids = [234, 345, 456]

# use var refer
top2.sub.var-refer = ${top.sub.key0}

/*
multi line
comments
*/
top2.sub.key2-other = has-char

# comments 2
top.sub.key3 = false

# slice list
top.sub.key4[0] = abc
top.sub.key4[1] = def

## --- comments 3 ---
top.sub.key5[0].f1 = ab
top.sub.key5[1].f2 = de

# multi line value
top.sub2.mline1 = """multi line
value
"""

# multi line value2
top.sub2.mline2 = this is \
multi line2 \
value
```

## 解析内容

```go
  p := properties.NewParser(
      properties.ParseEnv,
      properties.ParseInlineSlice,
  )
  p.Parse(text)
  fmt.Println("\ndata map:")
  dump.NoLoc(p.Data)
```

**Output**:

```shell
maputil.Data { #len=6
  "name": string("inhere"), #len=6
  "age": string("345"), #len=3
  "only-key": string(""), #len=0
  "env-key": string("/bin/zsh"), #len=8
  "top": map[string]interface {} { #len=2
    "sub": map[string]interface {} { #len=6
      "key5": []map[string]interface {} [ #len=2
        map[string]interface {} { #len=1
          "f1": string("ab"), #len=2
        },
        map[string]interface {} { #len=1
          "f2": string("de"), #len=2
        },
      ],
      "key0": string("a string value"), #len=14
      "key1": string("a quote value1"), #len=14
      "key2": string("a quote value2"), #len=14
      "key3": string("false"), #len=5
      "key4": []string [ #len=2
        string("abc"), #len=3
        string("def"), #len=3
      ],
    },
    "sub2": map[string]interface {} { #len=2
      "mline2": string("this is multi line2 value"), #len=25
      "mline1": string("multi line
value
"), #len=17
    },
  },
  "top2": map[string]interface {} { #len=2
    "sub": map[string]interface {} { #len=2
      "var-refer": string("a string value"), #len=14
      "key2-other": string("has-char"), #len=8
    },
    "inline": map[string]interface {} { #len=1
      "list": map[string]interface {} { #len=1
        "ids": []string [ #len=3
          string("234"), #len=3
          string("345"), #len=3
          string("456"), #len=3
        ],
      },
    },
  },
},
```

## 解析并绑定到结构体

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
