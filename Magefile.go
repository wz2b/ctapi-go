//go:build mage

package main

import (
	"fmt"
	sh "github.com/magefile/mage/sh"
	"strings"
)

func build(pkg string, os string, arch string) error {
	source := fmt.Sprintf("./cmd/%s", pkg)
	dest := fmt.Sprintf("out/%s-%s-%s", pkg, os, arch)
	if os == "windows" {
		dest = dest + ".exe"
	}

	err := sh.RunWith(map[string]string{"GOOS": os, "GOARCH": arch}, "go", "build", "-v", "-o", dest, source)

	if err != nil {
		fmt.Printf("%s\n", err)
	}
	return err
}

func Build() error {
	var err error
	cmds := []string{
		"lstags",
		"subscriber",
		"alarmlog",
		"cdb",
	}

	arches := []string{
		"windows/amd64",
		"windows/386",
	}

	if err = sh.Run("go", "mod", "download"); err != nil {
		return err
	}

	for _, command_package := range cmds {
		fmt.Printf("BUILDING %s\n", command_package)
		for _, arch := range arches {
			p := strings.Split(arch, "/")
			err = build(command_package, p[0], p[1])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

var Default = Build
