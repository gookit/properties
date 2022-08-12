package properties_test

import "github.com/gookit/properties"

func Example() {
	p, err := properties.Parse(`
# properties string
key = value
`)
	if err != nil {
		panic(err)
	}

	type MyConf struct {
		// ...
	}

	cfg := &MyConf{}
	err = p.MapStruct("", cfg)
	if err != nil {
		panic(err)
	}
}
