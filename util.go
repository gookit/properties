package properties

import (
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gookit/goutil/envutil"
	"github.com/mitchellh/mapstructure"
)

var refRegex = regexp.MustCompile(`^[a-z][a-z\d.]+`)

func parseVarRefName(val string) (string, bool) {
	if !strings.HasPrefix(val, VarRefStartChars) || !strings.HasSuffix(val, "}") {
		return "", false
	}

	refName := val[3 : len(val)-2]
	if refRegex.MatchString(refName) {
		return refName, true
	}
	return "", false
}

// ValDecodeHookFunc returns a mapstructure.DecodeHookFunc that parse ENV var, and more custom parse
func ValDecodeHookFunc(parseEnv, parseTime bool) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		str := data.(string)
		if len(str) < 2 {
			return str, nil
		}

		// start char is number(1-9)
		if str[0] > '0' && str[0] < '9' {
			// parse time string. eg: 10s
			if parseTime && t.Kind() == reflect.Int64 {
				dur, err := time.ParseDuration(str)
				if err == nil {
					return dur, nil
				}
			}
		} else if parseEnv { // parse ENV value
			str = envutil.ParseEnvValue(str)
		}

		return str, nil
	}
}
