package abi

import (
	"fmt"
	"os"

	"github.com/ability-sh/abi-lib/json"
	"gopkg.in/yaml.v2"
)

func Println(args ...interface{}) {
	fmt.Println(args...)
}

func Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func Panicln(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(0)
}

func Panicf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(0)
}

func PrintJSON(data interface{}) {
	b, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(b))
}

func PrintYAML(data interface{}) {
	b, _ := yaml.Marshal(data)
	fmt.Println(string(b))
}
