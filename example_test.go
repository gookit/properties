package properties_test

import (
	"fmt"
	"time"

	"github.com/gookit/properties"
)

func Example() {
	text := `
# properties string
name = inhere
age = 200
`

	p, err := properties.Parse(text)
	if err != nil {
		panic(err)
	}

	type MyConf struct {
		Name string `properties:"name"`
		Age  int    `properties:"age"`
	}

	cfg := &MyConf{}
	err = p.MapStruct("", cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println(*cfg)

	// Output:
	// {inhere 200}
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
	// name=inhere
	// age=300
}

func ExampleUnmarshal() {
	text := `
# properties string
name = inhere
age = 200

project.name = properties
project.version = v1.0.1
# parse time string
project.cache-time = 10s

project.repo.name = ${project.name}
project.repo.url = https://github.com/gookit/properties
`

	type Repo struct {
		Name string `properties:"name"`
		URL  string `properties:"url"`
	}

	type Project struct {
		Name      string        `properties:"name"`
		Version   string        `properties:"version"`
		CacheTime time.Duration `properties:"cache-time"`
		Repo      Repo          `properties:"repo"`
	}

	type MyConf struct {
		Name    string  `properties:"name"`
		Age     int     `properties:"age"`
		Project Project `properties:"project"`
	}

	cfg := &MyConf{}
	err := properties.Unmarshal([]byte(text), cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println(*cfg)

	// Output:
	// {inhere 200 {properties v1.0.1 10s {properties https://github.com/gookit/properties}}}
}
