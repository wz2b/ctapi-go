package main

import (
	"fmt"
	"golang.org/x/sys/windows"
	"golink/pkg/ctapi"
	"time"
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
	}
	defer api.CtClose(hConnection)

	list, err := api.NewList(hConnection)
	if err != nil {
		panic(err)
	}
	defer list.Free()

	var hTags []windows.Handle
	for i := 1; i <= 73; i = i + 1 {
		tagName := fmt.Sprintf("GAS%02d", i)
		fmt.Printf("Subscribe %s\n", tagName)
		hTag, err := list.Add(tagName)
		if err != nil {
			panic(err)
		}
		hTags = append(hTags, hTag)
	}

	fmt.Printf("list.Read() at %s\n", time.Now())
	err = list.Read()
	if err != nil {
		panic(err)
	}

	fmt.Printf("list dump() at %s\n", time.Now())

	for _, t := range hTags {
		value, err := list.GetFloatValue(t)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%.1f\n", value)
	}

}
