package properties

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
	"github.com/mitchellh/mapstructure"
)

// special chars consts
const (
	MultiLineValMarkS = "'''"
	MultiLineValMarkD = `"""`
	MultiLineValMarkQ = "\\"
	MultiLineCmtEnd   = "*/"
	VarRefStartChars  = "${"
)

// Parser for parse properties contents
type Parser struct {
	maputil.Data
	// last parse error
	err error
	// lex *lexer
	// text string
	opts *Options
	// key path map
	smap maputil.SMap
	// comments map
	comments map[string]string
}

// NewParser instance
func NewParser(optFns ...OpFunc) *Parser {
	p := &Parser{
		opts: newDefaultOption(),
		smap: make(maputil.SMap),
		Data: make(maputil.Data),
		// comments map
		comments: make(map[string]string),
	}

	return p.WithOptions(optFns...)
}

// WithOptions for the parser
func (p *Parser) WithOptions(optFns ...OpFunc) *Parser {
	for _, fn := range optFns {
		fn(p.opts)
	}
	return p
}

// Unmarshal parse properties text and decode to struct
func (p *Parser) Unmarshal(v []byte, ptr any) error {
	if err := p.ParseBytes(v); err != nil {
		return err
	}
	return p.MapStruct("", ptr)
}

// Parse text contents
func (p *Parser) Parse(text string) error {
	if text = strings.TrimSpace(text); text == "" {
		return errors.New("cannot input empty contents to parse")
	}
	return p.ParseFrom(strings.NewReader(text))
}

// ParseBytes text contents
func (p *Parser) ParseBytes(bs []byte) error {
	if len(bs) == 0 {
		return errors.New("cannot input empty contents to parse")
	}
	return p.ParseFrom(bytes.NewReader(bs))
}

// ParseFrom contents
func (p *Parser) ParseFrom(r io.Reader) error {
	p.err = nil
	s := bufio.NewScanner(r)

	var tok rune
	var line int
	var key, val, comments string

	// TODO
	// var ti tokenItem

	for s.Scan() { // split by '\n'
		if p.err != nil {
			break
		}

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
				tok = TokInvalid
			} else {
				comments += "\n"
			}
			continue
		}

		// multi line value
		if tok == TokMLValMarkS {
			if strings.HasSuffix(str, MultiLineValMarkS) { // end
				tok = TokInvalid
				val += str[:ln-3]
				p.setValue(key, val, comments)
				comments = "" // reset
			} else {
				val += str + "\n"
			}
			continue
		}

		// multi line value
		if tok == TokMLValMarkD {
			if strings.HasSuffix(str, MultiLineValMarkD) { // end
				tok = TokInvalid
				val += str[:ln-3]
				p.setValue(key, val, comments)
				comments = "" // reset
			} else {
				val += str + "\n"
			}
			continue
		}

		// multi line value
		if tok == TokMLValMarkQ {
			if strings.HasSuffix(str, MultiLineValMarkQ) { // go on
				val += str[:ln-1]
			} else {
				tok = TokInvalid
				val += str
				p.setValue(key, val, comments)
				comments = "" // reset
			}
			continue
		}

		// a line comments
		if str[0] == '#' || str[0] == '!' {
			tok = TokOLComments
			comments += raw
			continue
		}

		if str[0] == '/' {
			if ln < 2 {
				p.err = errorx.Rawf("invalid contents %q, at line#%d", str, line)
				continue
			}

			if str[1] == '/' {
				tok = TokOLComments
				comments += raw
				continue
			}

			// multi line comments start
			if str[1] == '*' {
				tok = TokMLComments
				comments += raw

				if strings.HasSuffix(str, MultiLineCmtEnd) {
					tok = TokInvalid
				} else {
					comments += "\n"
				}
				continue
			}
		}

		tok = TokValueLine

		// TODO ...
		// switch tok {
		// case TokOLComments:
		// case TokMLComments:
		//
		// case TokValueLine:
		//
		// }

		nodes := strutil.SplitNTrimmed(str, "=", 2)
		if len(nodes) != 2 {
			p.err = errorx.Rawf("invalid contents %q(should be KEY=VALUE), at line#%d", str, line)
			continue
		}

		key, val = nodes[0], nodes[1]
		if len(key) == 0 {
			p.err = errorx.Rawf("key cannot be empty: %q, at line#%d", str, line)
			continue
		}

		if p.opts.Debug {
			fmt.Printf("value line: %s = %s, tok=%d(%s)\n", key, val, tok, TokString(tok))
		}

		vln := len(val)
		// multi line value ended by \
		if vln > 0 && strings.HasSuffix(val, MultiLineValMarkQ) {
			tok = TokMLValMarkQ
			val = val[:vln-1]
			continue
		}

		if vln > 2 {
			// multi line value start
			hasPfx := strutil.HasOnePrefix(val, []string{MultiLineValMarkD, MultiLineValMarkS})
			if hasPfx && tok == TokValueLine {
				tok = TokMLValMarkS
				if val[0] == '"' {
					tok = TokMLValMarkD
				}
				val = val[3:] + "\n"
				continue
			}

			// clear quotes
			if val[0] == '"' || val[0] == '\'' {
				val = strutil.Unquote(val)
			} else if p.opts.InlineComment {
				// split inline comments
				var comment string
				val, comment = splitInlineComment(val)
				if len(comment) > 0 {
					if len(comments) > 0 {
						comments += "\n" + comment
					} else {
						comments += comment
					}
				}
			}
		}

		p.setValue(key, val, comments)
		comments = "" // reset
	}

	return p.err
}

// collect set value
func (p *Parser) setValue(key, value, comments string) {
	if len(comments) > 0 {
		p.comments[key] = comments
	}

	ln := len(value)
	if p.opts.TrimValue && ln > 0 {
		value = strings.TrimSpace(value)
	}

	if p.opts.ParseVar && ln > 3 {
		refName, ok := parseVarRefName(value)
		if ok {
			value = p.smap.Default(refName, value)
		}
	}

	var setVal any
	setVal = value
	p.smap[key] = value

	if p.opts.ParseEnv && ln > 3 {
		setVal = envutil.ParseEnvValue(value)
	}

	if p.opts.InlineSlice && ln > 2 {
		ss, ok := parseInlineSlice(value, ln)
		if ok {
			setVal = ss
		}
	}

	var keys []string
	if strings.ContainsRune(key, '.') {
		keys = strings.Split(key, ".")
	} else {
		keys = []string{key}
	}

	if p.opts.BeforeCollect != nil {
		setVal = p.opts.BeforeCollect(key, setVal)
	}

	// set value by keys
	if len(keys) == 1 {
		p.Data[key] = setVal
	} else if len(p.Data) == 0 {
		p.Data = maputil.MakeByKeys(keys, setVal)
	} else {
		err := p.Data.SetByKeys(keys, setVal)
		if err != nil {
			p.err = err
		}
	}
}

// collect set value
func (p *Parser) setValueByItem(ti tokenItem) {
	if !ti.Valid() {
		return
	}

	if len(ti.comments) > 0 {
		p.comments[ti.path] = strings.Join(ti.comments, "\n")
	}

	valueString := strings.Join(ti.values, "\n")
	p.smap[ti.path] = valueString

	// set value by keys
	if len(ti.keys) == 1 {
		p.Data[ti.path] = valueString
	} else {
		err := p.Data.SetByKeys(ti.keys, valueString)
		if err != nil {
			p.err = err
		}
	}
}

// ErrNotFound error
var ErrNotFound = errors.New("this key does not exists")

// MapStruct mapping data to a struct ptr
func (p *Parser) MapStruct(key string, ptr any) error {
	var data any
	if key == "" { // binding all data
		data = p.Data
	} else { // sub data of the p.Data
		var ok bool
		data, ok = p.Value(key)
		if !ok {
			return ErrNotFound
		}
	}

	decConf := p.opts.makeDecoderConfig()
	decConf.Result = ptr // set result ptr

	decoder, err := mapstructure.NewDecoder(decConf)
	if err == nil {
		err = decoder.Decode(data)
	}
	return err
}

// SMap data
func (p *Parser) SMap() maputil.SMap {
	return p.smap
}

// Comments data
func (p *Parser) Comments() map[string]string {
	return p.comments
}
