package properties

import (
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gookit/goutil/strutil"
	"github.com/mitchellh/mapstructure"
)

// ValDecodeHookFunc returns a mapstructure.DecodeHookFunc that parse time string
func ValDecodeHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		str := data.(string)
		ln := len(str)
		if ln < 2 {
			return str, nil
		}

		// start char is number(1-9) and end char is a-z.
		if str[0] > '0' && str[0] < '9' && str[ln-1] > 'a' {
			// parse time string. eg: 10s
			if t.Kind() == reflect.Int64 {
				dur, err := time.ParseDuration(str)
				if err == nil {
					return dur, nil
				}
			}
		}
		return str, nil
	}
}

// eg: ${some.other.key} -> some.other.key
var refRegex = regexp.MustCompile(`^[a-z][a-z\d.]+$`)

func parseVarRefName(val string) (string, bool) {
	if !strings.HasPrefix(val, VarRefStartChars) || !strings.HasSuffix(val, "}") {
		return "", false
	}

	refName := val[2 : len(val)-1]
	if refRegex.MatchString(refName) {
		return refName, true
	}
	return "", false
}

func parseInlineSlice(s string, ln int) (ss []string, ok bool) {
	// eg: [34, 56]
	if s[0] == '[' && s[ln-1] == ']' {
		return strutil.Split(s[1:ln-1], ","), true
	}
	return
}

func splitInlineComment(val string) (string, string) {
	if pos := strings.IndexRune(val, '#'); pos > -1 {
		return strings.TrimRight(val[0:pos], " "), val[pos:]
	}

	if pos := strings.Index(val, "//"); pos > -1 {
		return strings.TrimRight(val[0:pos], " "), val[pos:]
	}
	return val, ""
}
