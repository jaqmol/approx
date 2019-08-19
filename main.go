package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Definition ...
type Definition struct {
	Assign  map[string]string
	Env     map[string]string
	Command string
}

func main() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

	var dat map[string]interface{}
	if err := yaml.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
}
