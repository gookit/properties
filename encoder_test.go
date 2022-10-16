package properties_test

import (
	"fmt"
	"testing"
	"time"

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

func TestEncode_struct(t *testing.T) {
	type MyConf struct {
		Name   string        `properties:"name"`
		Age    int           `properties:"age"`
		Expire time.Duration `properties:"expire"`
	}

	myc := &MyConf{
		Name:   "inhere",
		Age:    234,
		Expire: time.Second * 3,
	}

	bs, err := properties.Encode(myc)
	assert.NoErr(t, err)
	str := string(bs)
	assert.StrContains(t, str, "name=inhere")
	assert.StrContains(t, str, "age=234")
	assert.StrContains(t, str, "expire=3000000000")
}

func TestEncode_error(t *testing.T) {
	bs, err := properties.Encode([]int{12, 34})
	assert.Nil(t, bs)
	assert.ErrMsg(t, err, "only allow encode map and struct data")
}
