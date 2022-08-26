package abi

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ability-sh/abi-lib/json"
	"gopkg.in/yaml.v2"
)

func SetConfig(fs_id string, fs_json string, fs_yaml string, fs_f string) {

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

	if fs_json != "" {

		var info interface{} = nil

		err = json.Unmarshal([]byte(fs_json), &info)

		if err != nil {
			Panicln(err)
		}

		_, err = registry.Send("/store/container/set.json", map[string]interface{}{"id": fs_id, "info": info})

		if err != nil {
			Panicln(err)
		}

		Println("SUCCESS")

		return
	}

	if fs_yaml != "" {

		var info interface{} = nil

		err = yaml.Unmarshal([]byte(fs_yaml), &info)

		if err != nil {
			Panicln(err)
		}

		_, err = registry.Send("/store/container/set.json", map[string]interface{}{"id": fs_id, "info": info})

		if err != nil {
			Panicln(err)
		}

		Println("SUCCESS")

		return
	}

	if fs_f != "" {

		var info interface{} = nil

		if strings.HasSuffix(fs_f, ".json") {

			b, err := ioutil.ReadFile(fs_f)

			if err != nil {
				Panicln(err)
			}

			err = json.Unmarshal(b, &info)

			if err != nil {
				Panicln(err)
			}

		} else if strings.HasSuffix(fs_f, ".yaml") {

			b, err := ioutil.ReadFile(fs_f)

			if err != nil {
				Panicln(err)
			}

			err = yaml.Unmarshal(b, &info)

			if err != nil {
				Panicln(err)
			}

		} else {
			Panicln("Only supports json, yaml files")
		}

		_, err = registry.Send("/store/container/set.json", map[string]interface{}{"id": fs_id, "info": info})

		if err != nil {
			Panicln(err)
		}

		Println("SUCCESS")

		return
	}

}
