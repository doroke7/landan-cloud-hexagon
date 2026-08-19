package main

import (
	"fmt"
	"reflect"
)

type Tool struct {
	Name string
}

func main() {

	oTool := Tool{}

	oType := reflect.TypeOf(oTool)

	fmt.Println(oType.Name())
}
