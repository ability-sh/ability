package abi

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-lib/http"
	"github.com/ability-sh/abi-lib/json"
	"gopkg.in/yaml.v2"
)

func md5File(f string) (string, error) {
	m := md5.New()
	fd, err := os.Open(f)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	_, err = io.Copy(m, fd)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(m.Sum(nil)), nil
}

func Publish(fs_file string, fs_ver string, fs_number string) {

	registry := GetRegistry()

	err := registry.Auth()

	if err != nil {
		Panicln(err)
	}

	if fs_file == "" {
		fs_file = "app.json"
	}

	var info interface{} = nil

	if strings.HasSuffix(fs_file, ".json") {

		b, err := ioutil.ReadFile(fs_file)

		if err != nil {
			Panicln(err)
		}

		err = json.Unmarshal(b, &info)

		if err != nil {
			Panicln(err)
		}

	} else if strings.HasSuffix(fs_file, ".yaml") {

		b, err := ioutil.ReadFile(fs_file)

		if err != nil {
			Panicln(err)
		}

		err = yaml.Unmarshal(b, &info)

		if err != nil {
			Panicln(err)
		}

	} else {
		Panicln("Only supports json, yaml files", fs_file)
	}

	dir, _ := filepath.Abs(filepath.Dir(fs_file))

	appid := dynamic.StringValue(dynamic.Get(info, "appid"), "")
	ver := dynamic.StringValue(dynamic.Get(info, "ver"), "")
	ability := dynamic.StringValue(dynamic.Get(info, "ability"), "")

	if fs_ver != "" {
		ver = fs_ver
		dynamic.Set(info, "ver", fs_ver)
	}

	if fs_number != "" {
		ver = fmt.Sprintf("%s-%s", strings.Split(ver, "-")[0], fs_number)
		dynamic.Set(info, "ver", ver)
	}

	if appid == "" {
		Panicln(fs_file, "appid not found")
	}

	if ver == "" {
		Panicln(fs_file, "ver not found")
	}

	if ability == "" {
		Panicln(fs_file, "ability not found")
	}

	vs := strings.Split(ability, "|")

	zipFile := func(abi string, dst string) error {

		config := dynamic.Get(info, abi)

		if config == nil {
			return fmt.Errorf("ability %s config not found", abi)
		}

		root := dynamic.StringValue(dynamic.Get(config, "root"), abi)

		fd, err := os.Create(dst)

		if err != nil {
			return err
		}

		defer fd.Close()

		zFile := zip.NewWriter(fd)

		err = filepath.Walk(filepath.Join(dir, root), func(path string, info fs.FileInfo, err error) error {

			if err != nil {
				return err
			}

			fh, err := zip.FileInfoHeader(info)

			fh.Name, _ = filepath.Rel(dir, path)

			if err != nil {
				return err
			}

			if info.IsDir() {
				fh.Name += "/"
			}

			w, err := zFile.CreateHeader(fh)

			if err != nil {
				return err
			}

			if !fh.Mode().IsRegular() {
				return nil
			}

			src, err := os.Open(path)

			if err != nil {
				return err
			}

			defer src.Close()

			_, err = io.Copy(w, src)

			if err != nil {
				return err
			}

			return nil
		})

		defer zFile.Close()

		if err != nil {
			return err
		}

		return nil
	}

	for _, v := range vs {

		config := dynamic.Get(info, v)

		if config == nil {
			Panicln(fs_file, "ability config not found", v)
		}

		dst := filepath.Join(dir, fmt.Sprintf("%s.zip", v))

		err = zipFile(v, dst)

		if err != nil {
			Panicln(err)
		}

		m, err := md5File(dst)

		if err != nil {
			Panicln(err)
		}

		dynamic.Set(config, "md5", m)

		defer os.Remove(dst)
	}

	for _, v := range vs {

		dst := filepath.Join(dir, fmt.Sprintf("%s.zip", v))

		data, err := ioutil.ReadFile(dst)

		if err != nil {
			Panicln(err)
		}

		rs, err := registry.Send("/store/app/ver/up.json", map[string]interface{}{"id": appid, "ver": ver, "ability": v})

		if err != nil {
			Panicln(err)
		}

		u := dynamic.StringValue(dynamic.Get(rs, "url"), "")

		res, err := http.NewHTTPRequest("PUT").SetURL(u, nil).SetBody(data).Send()

		if err != nil {
			Panicln(err)
		}

		if res.Code()/100 != 2 {
			Panicln(string(res.Body()))
		}
	}

	_, err = registry.Send("/store/app/ver/done.json", map[string]interface{}{"id": appid, "ver": ver, "info": info})

	if err != nil {
		Panicln(err)
	}

	PrintJSON(info)

}
