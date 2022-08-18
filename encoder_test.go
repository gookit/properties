package properties_test

import (
	"testing"

	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/properties"
)

func TestEncode(t *testing.T) {
	e := properties.NewEncoder()
	bs, err := e.Marshal(maputil.Data{
		"name": "inhere",
		"age":  234,
	})

	assert.NoErr(t, err)
	assert.NotEmpty(t, bs)
}
