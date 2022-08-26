package abi

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ability-sh/abi-lib/dynamic"
)

func GetAppConfig(fs_id string, fs_format string) {

	registry := GetRegistry()

	err := registry.Auth()

	if err != nil {
		Panicln(err)
	}

	rd := bufio.NewReader(os.Stdin)

	if fs_id == "" {
		fmt.Printf("Please enter a App ID: ")
		fs_id, err = rd.ReadString('\n')
		fs_id = strings.TrimSpace(fs_id)
		if err != nil {
			Panicln(err)
		}
	}

	rs, err := registry.Send("/store/app/get.json", map[string]string{"id": fs_id})

	if err != nil {
		Panicln(err)
	}

	if fs_format == "yaml" {
		PrintYAML(dynamic.Get(rs, "info"))
	} else {
		PrintJSON(dynamic.Get(rs, "info"))
	}

}
