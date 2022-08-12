package properties

import (
	"bytes"

	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/strutil"
)

// Parse text contents
func Parse(text string) (*Parser, error) {
	p := NewParser()
	return p, p.Parse(text)
}

func Marshal(v interface{}) ([]byte, error) {
	return Encode(v)
}

func Unmarshal(v []byte, ptr interface{}) error {
	return Decode(v, ptr)
}

func Encode(v interface{}) ([]byte, error) {
	mp, ok := v.(map[string]interface{})
	if !ok {
		return nil, errorx.Raw("only support encode map[string]any to Properties")
	}

	return encode(mp)
}

// Decode input string to struct ptr
func Decode(v []byte, ptr interface{}) error {
	p := NewParser()
	if err := p.ParseBytes(v); err != nil {
		return err
	}

	return p.MapStruct("", ptr)
}

func encode(mp map[string]interface{}) ([]byte, error) {
	var path string
	var buf bytes.Buffer

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
