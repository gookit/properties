// Package properties provide Java Properties format contents parse, marshal and unmarshal library.
package properties

// Parse properties text contents
func Parse(text string, optFns ...OpFunc) (*Parser, error) {
	p := NewParser(optFns...)
	return p, p.Parse(text)
}

// Marshal data(struct, map) to properties text
func Marshal(v any) ([]byte, error) {
	return NewEncoder().Encode(v)
}

// Unmarshal parse properties text and decode to struct
func Unmarshal(v []byte, ptr any) error {
	return Decode(v, ptr)
}

// Encode data(struct, map) to properties text
func Encode(v any) ([]byte, error) {
	return NewEncoder().Encode(v)
}

// Decode parse properties text and decode to struct
func Decode(v []byte, ptr any) error {
	p := NewParser()
	if err := p.ParseBytes(v); err != nil {
		return err
	}

	return p.Decode(ptr)
}
