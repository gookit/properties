package properties_test

import (
	"fmt"

	"github.com/gookit/properties"
)

func Example() {
	text := `
# properties string
key = value
`

	p, err := properties.Parse(text)
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

func ExampleMarshal() {
	type MyConf struct {
		Name string `properties:"name"`
		Age  int    `properties:"age"`
	}

	cfg := &MyConf{
		Name: "inhere",
		Age:  300,
	}

	bts, err := properties.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bts))

	// Output:
	// name = inhere
	// age = 200
}

func ExampleUnmarshal() {
	text := `
# properties string
name = inhere
age = 200

project.name = properties
project.version = v1.0.1

project.repo.name = ${project.name}
project.repo.url = https://github.com/gookit/properties
`

	type MyConf struct {
		Name string `properties:"name"`
		Age  int    `properties:"age"`
	}

	cfg := &MyConf{}
	err := properties.Unmarshal([]byte(text), cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg.Name)

	// Output:
	// inhere
}
