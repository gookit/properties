package properties

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
)

const (
	MultiLineValMarkS = "'''"
	MultiLineValMarkD = `"""`
	MultiLineCmtEnd   = "*/"
)

// Options for config
type Options struct {
	// ParseEnv parse ENV var name, default True. eg: "${SHELL}"
	ParseEnv      bool
	ParseVar      bool
	InlineComment bool
	// TrimMultiLine trim multi line value
	TrimMultiLine bool
	// TagName for binding data to struct
	TagName string
}

// Parser for parse properties contents
type Parser struct {
	maputil.Data
	err error
	lex *lexer
	// text string
	opts *Options
	smap maputil.SMap
	// comments map
	comments map[string]string
}

// NewParser instance
func NewParser() *Parser {
	return &Parser{
		smap: make(maputil.SMap),
		Data: make(maputil.Data),
		// comments map
		comments: make(map[string]string),
	}
}

// token consts
const (
	TokMLComments = 'C'
	TokMLValMarkS = 'm' // multi line value by single quotes: '''
	TokMLValMarkD = 'M' // multi line value by double quotes: """
)

// Parse text contents
func (p *Parser) Parse(text string) error {
	return p.ParseReader(strings.NewReader(text))
}

// ParseReader contents
func (p *Parser) ParseReader(r io.Reader) error {
	p.err = nil
	s := bufio.NewScanner(r)

	var tok rune
	var line int
	var key, val, comments string

	for s.Scan() { // split by '\n'
		line++

		raw := s.Text()
		str := strings.TrimSpace(raw)
		ln := len(str)
		if ln == 0 {
			continue
		}

		// multi line comments
		if tok == TokMLComments {
			comments += raw
			if strings.HasSuffix(str, MultiLineCmtEnd) {
				tok = 0
			} else {
				comments += "\n"
			}
			continue
		}

		// multi line value
		if tok == TokMLValMarkS {
			// end
			if strings.HasSuffix(str, MultiLineValMarkS) {
				tok = 0
				val += str[:ln-3]
				p.smap[key] = val
			} else {
				val += str + "\n"
			}
			continue
		}

		// multi line value
		if tok == TokMLValMarkD {
			// end
			if strings.HasSuffix(str, MultiLineValMarkD) {
				tok = 0
				val += str[:ln-3]
				p.smap[key] = val
			} else {
				val += str + "\n"
			}
			continue
		}

		if str[0] == '#' {
			comments += raw
			continue
		}

		if str[0] == '/' {
			if ln < 2 {
				p.err = errorx.Rawf("invalid string %q, at line#%d", str, line)
				continue
			}

			if str[1] == '/' {
				comments += raw
				continue
			}

			// multi line comments start
			if str[1] == '*' {
				tok = TokMLComments
				comments += raw

				if strings.HasSuffix(str, MultiLineCmtEnd) {
					tok = 0
				} else {
					comments += "\n"
				}
				continue
			}
		}

		nodes := strutil.SplitNTrimmed(str, "=", 2)
		if len(nodes) != 2 {
			p.err = errorx.Rawf("invalid format(key=val): %q, at line#%d", str, line)
			continue
		}

		key, val = nodes[0], nodes[1]
		if len(key) == 0 {
			p.err = errorx.Rawf("key is empty: %q, at line#%d", str, line)
			continue
		}

		fmt.Println("split: ", key, "=", val, ", tok=", tok)

		vln := len(val)
		if vln > 2 {
			// multi line value start
			hasPfx := strutil.HasOnePrefix(val, []string{"'''", `"""`})
			if hasPfx && tok == 0 {
				tok = TokMLValMarkS
				if val[0] == '"' {
					tok = TokMLValMarkD
				}
				val = val[3:] + "\n"

				// TODO end at inline: """value"""
				continue
			}

			// clear quotes
			if strings.HasPrefix(val, "'") {
				if pos := strings.IndexRune(val[1:], '\''); pos > -1 {
					val = val[1 : pos+1]
				}
			} else if strings.HasPrefix(val, `"`) {
				if pos := strings.IndexRune(val[1:], '"'); pos > -1 {
					val = val[1 : pos+1]
				}
			}

			// TODO inline comments
		}

		if len(comments) > 0 {
			p.comments[key] = comments
			comments = "" // reset
		}

		p.smap[key] = val
	}

	return nil
}

// Err last error
func (p *Parser) Err() error {
	return p.err
}

// SMap data
func (p *Parser) SMap() maputil.SMap {
	return p.smap
}

// Comments data
func (p *Parser) Comments() map[string]string {
	return p.comments
}
