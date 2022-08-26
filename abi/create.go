package abi

import (
	"io/ioutil"
	"strings"

	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-lib/json"
	"gopkg.in/yaml.v2"
)

func Create(fs_json string, fs_yaml string, fs_f string) {

	registry := GetRegistry()

	err := registry.Auth()

	if err != nil {
		Panicln(err)
	}

	var info interface{} = nil

	if fs_json != "" {

		var info interface{} = nil

		err = json.Unmarshal([]byte(fs_json), &info)

		if err != nil {
			Panicln(err)
		}

	} else if fs_yaml != "" {

		var info interface{} = nil

		err = yaml.Unmarshal([]byte(fs_yaml), &info)

		if err != nil {
			Panicln(err)
		}
	} else if fs_f != "" {

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

		}
	}

	if info == nil {
		info = map[string]interface{}{}
	}

	rs, err := registry.Send("/store/container/create.json", map[string]interface{}{"info": info})

	if err != nil {
		Panicln(err)
	}

	Println("Container ID:", dynamic.StringValue(dynamic.Get(rs, "id"), ""))
	Println("Container Secret:", dynamic.StringValue(dynamic.Get(rs, "secret"), ""))

}
