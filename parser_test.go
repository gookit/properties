package properties_test

import (
	"fmt"
	"testing"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/properties"
)

func TestParser_Parse(t *testing.T) {
	text := `
name = inhere
age = 345
only-key = 
env-key = ${SHELL | bash}

 ##### comments1
top.sub.key0 = a string value
top.sub.key1 = "a quote value1"
top.sub.key2 = 'a quote value2'
/* comments 1.1 */
top.sub.key3 = 234

# inline list
top2.inline.list.ids = [234, 345, 456]

# use var refer
top2.sub.var-refer = ${top.sub.key0}

/*
multi line
comments
*/
top2.sub.key2-other = has-char

# comments 2
top.sub.key3 = false

# slice list
top.sub.key4[0] = abc
top.sub.key4[1] = def

## --- comments 3 ---
top.sub.key5[0].f1 = ab
top.sub.key5[1].f2 = de

top.sub2.mline1 = """multi line
value
"""

top.sub2.mline2 = this is \
multi line2 \
value
`

	p := properties.NewParser(
		properties.WithDebug,
		properties.ParseEnv,
		properties.ParseInlineSlice,
	)
	err := p.Parse(text)
	assert.NoErr(t, err)
	fmt.Println("\ndata map:")
	dump.NoLoc(p.Data)

	fmt.Println("\nstring map:")
	dump.NoLoc(p.SMap())

	fmt.Println("\ncomments:")
	dump.NoLoc(p.Comments())

	assert.Eq(t, 345, p.Int("age"))
	assert.Eq(t, "inhere", p.Str("name"))
	assert.Eq(t, "a quote value1", p.Str("top.sub.key1"))
	assert.Eq(t, []string{"234", "345", "456"}, p.Strings("top2.inline.list.ids"))
	assert.Eq(t, "[234, 345, 456]", p.SMap().Get("top2.inline.list.ids"))
}

func TestParser_WithOptions_parseVar(t *testing.T) {
	text := `key0 = val1
top.sub.key0 = a string value
top2.sub.var-refer = ${top.sub.key0}
`

	p := properties.NewParser()
	err := p.Parse(text)
	assert.NoErr(t, err)

	smp := p.SMap()
	assert.Eq(t, "a string value", smp.Str("top.sub.key0"))
	assert.Eq(t, "a string value", smp.Str("top2.sub.var-refer"))
}

func TestParser_Parse_multiLineValS(t *testing.T) {
	text := `key0 = val1
top.sub2.mline1 = '''multi line
value
'''
key1 = val2
`

	p := properties.NewParser()
	err := p.Parse(text)
	assert.NoErr(t, err)
	smp := p.SMap()
	assert.NotEmpty(t, smp)
	assert.ContainsKeys(t, smp, []string{"key0", "key1", "top.sub2.mline1"})
	assert.Eq(t, "multi line\nvalue\n", smp.Str("top.sub2.mline1"))

	// start and end mark at new line
	text = `key0 = val1
top.sub2.mline1 = '''
multi line
value
'''
key1 = val2
`

	p = properties.NewParser()
	assert.NoErr(t, p.Parse(text))
	smp = p.SMap()
	assert.NotEmpty(t, smp)
	assert.ContainsKey(t, smp, "top.sub2.mline1")
	assert.Eq(t, "\nmulti line\nvalue\n", smp.Str("top.sub2.mline1"))

	// value at end line
	text = `key0 = val1
top.sub2.mline1 = '''multi line
value'''
key1 = val2
`

	p = properties.NewParser()
	assert.NoErr(t, p.Parse(text))
	smp = p.SMap()
	assert.NotEmpty(t, smp)
	assert.ContainsKey(t, smp, "top.sub2.mline1")
	assert.Eq(t, "multi line\nvalue", smp.Str("top.sub2.mline1"))

	// empty value
	text = `key0 = val1
top.sub2.mline1 = '''
'''
key1 = val2
`

	p = properties.NewParser()
	assert.NoErr(t, p.Parse(text))
	smp = p.SMap()
	assert.NotEmpty(t, smp)
	assert.ContainsKey(t, smp, "top.sub2.mline1")
	assert.Eq(t, "\n", smp.Str("top.sub2.mline1"))
}

func TestParser_Parse_multiLineValD(t *testing.T) {
	text := `key0 = val1
top.sub2.mline1 = """multi line
value
"""
key1 = val2
`

	p := properties.NewParser()
	err := p.Parse(text)
	assert.NoErr(t, err)
	smp := p.SMap()
	assert.NotEmpty(t, smp)
	assert.ContainsKeys(t, smp, []string{"key0", "key1", "top.sub2.mline1"})
	assert.Eq(t, "multi line\nvalue\n", smp.Str("top.sub2.mline1"))

	// start and end mark at new line
	text = `key0 = val1
top.sub2.mline1 = """
multi line
value
"""
key1 = val2
`

	p = properties.NewParser()
	assert.NoErr(t, p.Parse(text))
	smp = p.SMap()
	assert.NotEmpty(t, smp)
	assert.ContainsKey(t, smp, "top.sub2.mline1")
	assert.Eq(t, "\nmulti line\nvalue\n", smp.Str("top.sub2.mline1"))

	// value at end line
	text = `key0 = val1
top.sub2.mline1 = """multi line
value"""
key1 = val2
`

	p = properties.NewParser()
	assert.NoErr(t, p.Parse(text))
	smp = p.SMap()
	assert.NotEmpty(t, smp)
	assert.ContainsKey(t, smp, "top.sub2.mline1")
	assert.Eq(t, "multi line\nvalue", smp.Str("top.sub2.mline1"))

	// empty value
	text = `
key0 = val1
top.sub2.mline1 = """
"""
key1 = val2
`

	p = properties.NewParser()
	assert.NoErr(t, p.Parse(text))
	smp = p.SMap()
	assert.NotEmpty(t, smp)
	assert.ContainsKey(t, smp, "top.sub2.mline1")
	assert.Eq(t, "\n", smp.Str("top.sub2.mline1"))
}

func TestParser_Parse_multiLineValQ(t *testing.T) {
	text := `key0 = val1
top.sub2.mline1 = multi line \
value
key1 = val2
`

	p := properties.NewParser()
	err := p.Parse(text)
	assert.NoErr(t, err)
	smp := p.SMap()
	assert.NotEmpty(t, smp)
	assert.ContainsKeys(t, smp, []string{"key0", "key1", "top.sub2.mline1"})
	assert.Eq(t, "multi line value", smp.Str("top.sub2.mline1"))

}
