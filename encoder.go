package properties

import (
	"bytes"
	"errors"
	"reflect"
	"strings"

	"github.com/gookit/goutil/reflects"
	"github.com/mitchellh/mapstructure"
)

// Encoder struct
type Encoder struct {
	buf bytes.Buffer
	// TagName for encode a struct. default: properties
	TagName string
	// comments map data. TODO
	// key is path name, value is comments
	// comments map[string]string
}

// NewEncoder instance.
func NewEncoder() *Encoder {
	return &Encoder{
		TagName: DefaultTagName,
	}
}

// Marshal data(struct, map) to properties text
func (e *Encoder) Marshal(v any) ([]byte, error) {
	return e.Encode(v)
}

// Encode data(struct, map) to properties text
func (e *Encoder) Encode(v any) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	if err := e.encode(v); err != nil {
		return nil, err
	}
	return e.buf.Bytes(), nil
}

// Encode data(struct, map) to properties text
func (e *Encoder) encode(v any) error {
	rv := reflect.Indirect(reflect.ValueOf(v))

	// convert struct to map[string]any
	if rv.Kind() == reflect.Struct {
		mp := make(map[string]any)
		cfg := &mapstructure.DecoderConfig{
			TagName: e.TagName,
			Result:  &mp,
		}

		decoder, _ := mapstructure.NewDecoder(cfg)
		if err := decoder.Decode(v); err != nil {
			return err
		}

		rv = reflect.ValueOf(mp)
	} else if rv.Kind() != reflect.Map {
		return errors.New("only allow encode map and struct data")
	}

	// TODO collect to map[string]string then sort keys
	reflects.FlatMap(rv, e.writeln)
	return nil
}

func (e *Encoder) writeln(path string, rv reflect.Value) {
	e.buf.WriteString(path)
	e.buf.WriteByte('=')

	val := reflects.String(rv)
	if rv.Kind() == reflect.String && strings.ContainsRune(val, '\n') {
		val = strings.Replace(val, "\n", "\\\n", -1)
	}

	e.buf.WriteString(val)
	e.buf.WriteByte('\n')
}
