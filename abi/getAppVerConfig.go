package abi

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetAppVerConfig(fs_id string, fs_ver string, fs_format string) {

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

	if fs_ver == "" {
		fmt.Printf("Please enter a App Ver: ")
		fs_ver, err = rd.ReadString('\n')
		fs_ver = strings.TrimSpace(fs_ver)
		if err != nil {
			Panicln(err)
		}
	}

	rs, err := registry.Send("/store/app/ver/info/get.json", map[string]string{"id": fs_id, "ver": fs_ver})

	if err != nil {
		Panicln(err)
	}

	if fs_format == "yaml" {
		PrintYAML(rs)
	} else {
		PrintJSON(rs)
	}

}
