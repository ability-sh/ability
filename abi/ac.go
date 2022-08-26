package abi

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-lib/http"
)

type acRegistry struct {
	baseURL string
	token   string
	key     string
}

func NewACRegistry(baseURL string) Registry {
	if baseURL == "" {
		baseURL = "https://ac.ability.sh"
	}
	m := md5.New()
	m.Write([]byte(baseURL))
	return &acRegistry{baseURL: baseURL, key: hex.EncodeToString(m.Sum(nil))}
}

func (ac *acRegistry) login() error {

	Println("Please input your email: ")

	var email string = ""

	fmt.Scanln(&email)

	res, err := http.
		NewHTTPRequest("POST").
		SetURL(fmt.Sprintf("%s/store/mail/send.json", ac.baseURL), nil).
		SetUrlencodeBody(map[string]string{"email": email}).
		Send()

	if err != nil {
		return err
	}

	if res.Code() != 200 {
		return fmt.Errorf("%d %s", res.Code(), string(res.Body()))
	}

	body, err := res.PraseBody()

	if err != nil {
		return err
	}

	if dynamic.IntValue(dynamic.Get(body, "errno"), 0) == 200 {
		Printf("Captcha code has been sent to email %s\n", email)
		Println("Please enter captcha code: ")
		var code string = ""
		fmt.Scanln(&code)

		res, err = http.
			NewHTTPRequest("POST").
			SetURL(fmt.Sprintf("%s/store/login.json", ac.baseURL), nil).
			SetUrlencodeBody(map[string]string{"email": email, "code": code}).
			Send()

		if err != nil {
			return err
		}

		if res.Code() != 200 {
			return fmt.Errorf("%d %s", res.Code(), string(res.Body()))
		}

		body, err = res.PraseBody()

		if err != nil {
			return err
		}

		if dynamic.IntValue(dynamic.Get(body, "errno"), 0) == 200 {

			ac.token = dynamic.StringValue(dynamic.GetWithKeys(body, []string{"data", "token"}), "")
			ac.setToken(ac.token)

		} else {
			return fmt.Errorf("%s", dynamic.StringValue(dynamic.Get(body, "errmsg"), "Please check network settings"))
		}
	} else {
		return fmt.Errorf("%s", dynamic.StringValue(dynamic.Get(body, "errmsg"), "Please check network settings"))
	}

	return nil
}

func (ac *acRegistry) getToken() string {
	b, _ := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".ability-sh", ac.key))
	if b != nil {
		return string(b)
	}
	return ""
}

func (ac *acRegistry) setToken(v string) error {
	dir := filepath.Join(os.Getenv("HOME"), ".ability-sh")
	os.MkdirAll(dir, os.ModePerm)
	return ioutil.WriteFile(filepath.Join(dir, ac.key), []byte(v), os.ModePerm)
}

func (ac *acRegistry) SetToken(token string) {
	ac.token = token
}

func (ac *acRegistry) Logout() {
	os.RemoveAll(filepath.Join(os.Getenv("HOME"), ".ability-sh", ac.key))
}

func (ac *acRegistry) Auth() error {

	if ac.token == "" {
		ac.token = ac.getToken()
	}

	if ac.token == "" {
		err := ac.login()
		if err != nil {
			return err
		}
	} else {

		res, err := http.
			NewHTTPRequest("POST").
			SetURL(fmt.Sprintf("%s/store/user/get.json", ac.baseURL), nil).
			SetUrlencodeBody(map[string]string{"token": ac.token}).
			Send()

		if err != nil {
			return err
		}

		if res.Code() != 200 {
			return fmt.Errorf("%d %s", res.Code(), string(res.Body()))
		}

		body, err := res.PraseBody()

		if err != nil {
			return err
		}

		if dynamic.IntValue(dynamic.Get(body, "errno"), 0) != 200 {
			err = ac.login()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ac *acRegistry) Send(path string, inputData interface{}) (interface{}, error) {

	dynamic.Set(inputData, "token", ac.token)

	res, err := http.
		NewHTTPRequest("POST").
		SetURL(fmt.Sprintf("%s%s", ac.baseURL, path), nil).
		SetJSONBody(inputData).
		Send()

	if err != nil {
		return nil, err
	}

	if res.Code() != 200 {
		return nil, fmt.Errorf("%d %s", res.Code(), string(res.Body()))
	}

	body, err := res.PraseBody()

	if err != nil {
		return nil, err
	}

	if dynamic.IntValue(dynamic.Get(body, "errno"), 0) == 200 {
		return dynamic.Get(body, "data"), nil
	} else {
		return nil, fmt.Errorf("%s", dynamic.StringValue(dynamic.Get(body, "errmsg"), "Please check network settings"))
	}
}
