package abi

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Unapprove(fs_appid string, fs_containerId string) {

	registry := GetRegistry()

	err := registry.Auth()

	if err != nil {
		Panicln(err)
	}

	rd := bufio.NewReader(os.Stdin)

	if fs_appid == "" {
		fmt.Printf("Please enter a App ID: ")
		fs_appid, err = rd.ReadString('\n')
		fs_appid = strings.TrimSpace(fs_appid)
		if err != nil {
			Panicln(err)
		}
	}

	if fs_containerId == "" {
		fmt.Printf("Please enter a Container ID: ")
		fs_containerId, err = rd.ReadString('\n')
		fs_containerId = strings.TrimSpace(fs_containerId)
		if err != nil {
			Panicln(err)
		}
	}

	_, err = registry.Send("/store/app/unapprove.json", map[string]interface{}{"id": fs_appid, "containerId": fs_containerId})

	if err != nil {
		Panicln(err)
	}

	Println("SUCCESS")

}
