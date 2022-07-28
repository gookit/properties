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
 ##### comments1
top.sub.key0 = a string value
top.sub.key1 = "a string value"
/* comments 1.1 */
top.sub.key2 = 234
/*
multi line
comments
*/
top.sub.key2-other = has-char
# comments 2
top.sub.key3 = false
top.sub.key4[0] = abc # comments at end1
top.sub.key4[1] = def // comments at end2
## --- comments 3 ---
top.sub.key5[0].f1 = ab
top.sub.key5[1].f2 = de
invalid line
top.sub2.mline1 = """multi line
value
"""
`

	p := properties.NewParser()
	err := p.Parse(text)
	assert.NoErr(t, err)

	fmt.Println("string map:")
	dump.NoLoc(p.SMap())

	fmt.Println("comments:")
	dump.NoLoc(p.Comments())
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
