package properties_test

import (
	"testing"
	"time"

	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/properties"
)

func TestOptions_basic(t *testing.T) {
	opt := &properties.Options{}

	fn := properties.WithTagName("config")
	fn(opt)

	assert.Eq(t, "config", opt.TagName)
}

func TestOptions_InlineComment(t *testing.T) {
	text := `
key = value # inline comments
key2 = value2 // inline comments
key3 = value3
`

	p := properties.NewParser()
	err := p.Parse(text)
	assert.NoErr(t, err)
	assert.Eq(t, "value # inline comments", p.Str("key"))

	p = properties.NewParser(func(opts *properties.Options) {
		opts.InlineComment = true
	})
	err = p.Parse(text)
	assert.NoErr(t, err)
	assert.Eq(t, "value", p.Str("key"))
}

func TestParseTime(t *testing.T) {
	text := `
// properties string
name = inhere
age = 200
expire = 3s
`

	type MyConf struct {
		Name   string        `properties:"name"`
		Age    int           `properties:"age"`
		Expire time.Duration `properties:"expire"`
	}

	p, err := properties.Parse(text, properties.ParseTime)
	assert.NoErr(t, err)

	cfg := &MyConf{}
	err = p.MapStruct("", cfg)
	assert.NoErr(t, err)
	assert.Eq(t, "inhere", cfg.Name)
	assert.Eq(t, 3*time.Second, cfg.Expire)

	err = p.MapStruct("not-found", cfg)
	assert.Err(t, err)
}

func TestOptions_BeforeCollect(t *testing.T) {
	text := `
// properties string
name = inhere
age = 200
expire = 3s
`

	type MyConf struct {
		Name   string        `properties:"name"`
		Age    int           `properties:"age"`
		Expire time.Duration `properties:"expire"`
	}

	// tests Unmarshal, BeforeCollect
	p := properties.NewParser(properties.ParseTime, func(opts *properties.Options) {
		opts.BeforeCollect = func(name string, val any) interface{} {
			if name == "name" {
				return strutil.Upper(val.(string))
			}
			return val
		}
	})

	cfg := &MyConf{}
	err := p.Unmarshal([]byte(text), cfg)
	assert.NoErr(t, err)
	assert.Eq(t, "INHERE", cfg.Name)
	assert.Eq(t, 3*time.Second, cfg.Expire)

}
