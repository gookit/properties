package properties

import (
	"bytes"

	"github.com/gookit/goutil/strutil"
	"github.com/mitchellh/mapstructure"
)

// Encoder struct
type Encoder struct {
	// TagName for encode a struct. default: properties
	TagName string
}

// NewEncoder instance.
func NewEncoder() *Encoder {
	return &Encoder{
		TagName: DefaultTagName,
	}
}

// Marshal data(struct, map) to properties text
func (e *Encoder) Marshal(v interface{}) ([]byte, error) {
	return e.Encode(v)
}

// Encode data(struct, map) to properties text
func (e *Encoder) Encode(v interface{}) ([]byte, error) {
	mp, ok := v.(map[string]interface{})

	// try convert v to map[string]interface{}
	if !ok {
		mp = make(map[string]interface{})
		cfg := &mapstructure.DecoderConfig{
			TagName: e.TagName,
			Result:  &mp,
		}

		decoder, err := mapstructure.NewDecoder(cfg)
		if err != nil {
			return nil, err
		}

		if err := decoder.Decode(v); err != nil {
			return nil, err
		}
	}

	return e.encode(mp)
}

func (e *Encoder) encode(mp map[string]interface{}) ([]byte, error) {
	var path string
	var buf bytes.Buffer

	// TODO sort keys

	// TODO...
	for name, val := range mp {
		path = name
		buf.WriteString(path)
		buf.WriteByte('=')
		buf.WriteString(strutil.QuietString(val))
		buf.WriteByte('\n')
	}

	return buf.Bytes(), nil
}
