package main

import (
	"fmt"
	"golink/pkg/ctapi"
)

var dllPath = "C:\\Program Files (x86)\\AVEVA Plant SCADA\\Bin\\Bin (x64)"

//var dllPath_bytes = append([]byte(dllPath), 0)
//var dllPath_ptr = uintptr(unsafe.Pointer(&dllPath_bytes[0]))

// var dotnetpath = "C:\\Program Files (x86)\\Reference Assemblies\\Microsoft\\Framework\\.NETFramework\\v4.8"
// var dotnetpath_bytes = append([]byte(dotnetpath), 0)
// var dotnetpath_ptr = uintptr(unsafe.Pointer(&dotnetpath_bytes[0]))

func main() {
	api, err := ctapi.Init("CtApi.dll", dllPath)
	if err != nil {
		panic(err)
	}

	handle, err := api.CtOpen()
	if err != nil {
		panic(err)
	} else {
		println(handle)
	}
	defer api.CtClose(handle)

	const table = "Tag"
	const property = "TAG"
	for hObj := range api.FindAll(handle, table) {
		v, e := api.GetStringProperty(hObj, property)
		if e != nil {
			fmt.Printf("There was an error getting the string property: %v\n", e)
		} else {
			fmt.Printf("%v \n", v)
		}
	}
}
