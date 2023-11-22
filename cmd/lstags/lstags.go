package main

import (
	"fmt"
	"github.com/wz2b/ctapi-go/ctapi"
)

var dllPath = "C:\\Program Files (x86)\\AVEVA Plant SCADA\\Bin\\Bin (x64)"

func main() {
	api, err := ctapi.Init("CtApi.dll", dllPath)
	if err != nil {
		panic(err)
	}

	hConnection, err := api.CtOpen()
	if err != nil {
		panic(err)
	} else {
		println(hConnection)
	}
	defer api.CtClose(hConnection)

	const table = "Tag"
	const property = "TAG"
	for hObj := range api.FindAll(hConnection, table) {
		v, e := api.GetPropertyAsString(hObj, property)
		if e != nil {
			fmt.Printf("There was an error getting the string property: %v\n", e)

		}
		s, e := api.GetPropertyAsString(hObj, "TYPE")
		if e != nil {
			fmt.Printf("There was an error getting the string property: %v\n", e)
		}
		//s := "?"

		fmt.Printf("%20s %s\n", v, s)
	}
}
