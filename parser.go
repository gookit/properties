package properties

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil/textscan"
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
	ts := textscan.NewScanner(r)
	ts.AddMatchers(
		&textscan.CommentsMatcher{
			InlineChars: []byte{'#', '!'},
		},
		&textscan.KeyValueMatcher{
			InlineComment: p.opts.InlineComment,
			MergeComments: true,
		},
	)

	// scan and parsing
	for ts.Scan() {
		tok := ts.Token()

		// collect value
		if tok.Kind() == textscan.TokValue {
			p.setValue(tok.(*textscan.ValueToken))
		}
	}

	p.err = ts.Err()
	return p.err
}

// collect set value
func (p *Parser) setValue(tok *textscan.ValueToken) {
	var value string
	if tok.Mark() == textscan.MultiLineValMarkQ {
		value = strings.Join(tok.Values(), "")
	} else {
		value = tok.Value()
	}

	key := tok.Key()
	if tok.HasComment() {
		p.comments[key] = tok.Comment()
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

// ErrNotFound error
var ErrNotFound = errors.New("this key does not exists")

// Decode the parsed data to struct ptr
func (p *Parser) Decode(ptr any) error {
	return p.MapStruct("", ptr)
}

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
