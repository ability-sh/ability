package abi

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ability-sh/abi-lib/dynamic"
)

func SetSecret(fs_id string) {

	registry := GetRegistry()

	err := registry.Auth()

	if err != nil {
		Panicln(err)
	}

	rd := bufio.NewReader(os.Stdin)

	if fs_id == "" {
		fmt.Printf("Please enter a Container ID: ")
		fs_id, err = rd.ReadString('\n')
		fs_id = strings.TrimSpace(fs_id)
		if err != nil {
			Panicln(err)
		}
	}

	rs, err := registry.Send("/store/container/set.json", map[string]interface{}{"id": fs_id, "secret": true})

	if err != nil {
		Panicln(err)
	}

	Println("SECRET:", dynamic.StringValue(dynamic.Get(rs, "secret"), ""))

}
