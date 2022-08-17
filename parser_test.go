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

 ##### comments1
top.sub.key0 = a string value
top.sub.key1 = "a string value"
/* comments 1.1 */
top.sub.key2 = 234

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

top.sub2.mline2 = this is\
multi line2\
value
`

	p := properties.NewParser()
	err := p.Parse(text)
	assert.NoErr(t, err)
	fmt.Println("\ndata map:")
	dump.NoLoc(p.Data)

	fmt.Println("\nstring map:")
	dump.NoLoc(p.SMap())

	fmt.Println("\ncomments:")
	dump.NoLoc(p.Comments())
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
	dump.P(smp)
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
