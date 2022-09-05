package properties_test

import (
	"fmt"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/properties"
)

func TestEncode(t *testing.T) {
	e := properties.NewEncoder()
	bs, err := e.Marshal(map[string]any{
		"name": "inhere",
		"age":  234,
		"str1": "a string",
		"str2": "a multi \nline string",
		"top": map[string]any{
			"sub0": "val0",
			"sub1": []string{"val1-0", "val1-1"},
		},
	})

	str := string(bs)
	fmt.Println(str)
	assert.NoErr(t, err)
	assert.NotEmpty(t, bs)
	assert.StrContains(t, str, "name=inhere")
	assert.StrContains(t, str, "top.sub1[0]=val1-0")
	assert.StrContains(t, str, "str2=a multi \\\nline string")

	bs, err = properties.Marshal(nil)
	assert.NoErr(t, err)
	assert.Nil(t, bs)
}
