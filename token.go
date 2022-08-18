package properties

import "strings"

const (
	// TokInvalid invalid token
	TokInvalid rune = 0
	// TokOLComments one line comments start by !,#
	TokOLComments = 'c'
	// TokMLComments multi line comments by /* */
	TokMLComments = 'C'
	// TokILComments inline comments
	TokILComments = 'i'
	// TokValueLine value line
	TokValueLine = 'v'
	// TokMLValMarkS multi line value by single quotes: '''
	TokMLValMarkS = 'm'
	// TokMLValMarkD multi line value by double quotes: """
	TokMLValMarkD = 'M'
	// TokMLValMarkQ multi line value by left slash quote: \
	TokMLValMarkQ = 'q'
)

// TokString name
func TokString(tok rune) string {
	switch tok {
	case TokOLComments:
		return "LINE_COMMENT"
	case TokILComments:
		return "INLINE_COMMENT"
	case TokMLComments:
		return "MLINE_COMMENT"
	case TokValueLine:
		return "VALUE_LINE"
	case TokMLValMarkS, TokMLValMarkD, TokMLValMarkQ:
		return "MLINE_VALUE"
	default:
		return "INVALID"
	}
}

type tokenItem struct {
	// see TokValueLine
	kind rune
	// key path string. eg: top.sub.some-key
	path string
	keys []string

	// token value
	value string
	// for multi line value.
	values []string
	// for multi line comments.
	comments []string
}

func newTokenItem(path, value string, kind rune) *tokenItem {
	tk := &tokenItem{
		kind:  kind,
		value: value,
	}

	tk.setPath(path)
	return tk
}

func (ti *tokenItem) setPath(path string) {
	// TODO check path valid
	ti.path = path

	if strings.ContainsRune(path, '.') {
		ti.keys = strings.Split(path, ".")
	}
}

// Valid of the token data.
func (ti *tokenItem) addValue(val string) {
	ti.values = append(ti.values, val)
}

// Valid of the token data.
func (ti *tokenItem) Valid() bool {
	return ti.kind != 0
}
