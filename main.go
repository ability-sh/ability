package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/ability-sh/ability/abi"
	"github.com/ability-sh/ability/commander"
)

type nilWriter struct {
}

func (w *nilWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func main() {

	app := commander.NewCommand("ability")

	app.
		SetString("registry", "https://ac.ability.sh", "").
		SetString("token", "", "").
		SetBool("help", "").
		SetAction(func(cmd *commander.Command) bool {

			if cmd.Bool("help") {
				cmd.Usage()
				return true
			}

			abi.SetRegistry(abi.NewACRegistry(cmd.String("registry")))

			token := cmd.String("token")

			if token != "" {
				abi.GetRegistry().SetToken(token)
			}

			return false
		})

	app.
		SubCommand("login").
		SetAction(func(cmd *commander.Command) bool {
			abi.Login()
			return true
		})

	app.
		SubCommand("logout").
		SetAction(func(cmd *commander.Command) bool {
			abi.Logout()
			return true
		})

	app.
		SubCommand("create").
		SetString("json", "", "").
		SetString("yaml", "", "").
		SetString("file", "", "").
		SetAction(func(cmd *commander.Command) bool {
			abi.Create(cmd.String("json"), cmd.String("yaml"), cmd.String("file"))
			return true
		})

	app.
		SubCommand("setsecret").
		SetString("id", "", "Container ID").
		SetAction(func(cmd *commander.Command) bool {
			abi.SetSecret(cmd.String("id"))
			return true
		})

	app.
		SubCommand("setconfig").
		SetString("id", "", "Container ID").
		SetString("json", "", "").
		SetString("yaml", "", "").
		SetString("file", "", "").
		SetAction(func(cmd *commander.Command) bool {
			abi.SetConfig(cmd.String("id"), cmd.String("json"), cmd.String("yaml"), cmd.String("file"))
			return false
		})

	app.
		SubCommand("getconfig").
		SetString("id", "", "Container ID").
		SetString("format", "", "json|yaml").
		SetAction(func(cmd *commander.Command) bool {

			abi.GetConfig(cmd.String("id"), cmd.String("format"))

			return true
		})

	{
		s := app.SubCommand("app")

		s.
			SubCommand("create").
			SetString("json", "", "").
			SetString("yaml", "", "").
			SetString("file", "", "").
			SetAction(func(cmd *commander.Command) bool {

				abi.CreateApp(cmd.String("json"), cmd.String("yaml"), cmd.String("file"))

				return true
			})

		s.
			SubCommand("setconfig").
			SetString("appid", "", "App ID").
			SetString("json", "", "").
			SetString("yaml", "", "").
			SetString("file", "", "").
			SetAction(func(cmd *commander.Command) bool {

				abi.SetAppConfig(cmd.String("appid"), cmd.String("json"), cmd.String("yaml"), cmd.String("file"))

				return true
			})

		s.
			SubCommand("getconfig").
			SetString("appid", "", "App ID").
			SetString("format", "", "json|yaml").
			SetAction(func(cmd *commander.Command) bool {

				abi.GetAppConfig(cmd.String("appid"), cmd.String("format"))

				return true
			})

		s.
			SubCommand("getver").
			SetString("appid", "", "App ID").
			SetString("ver", "", "App Ver").
			SetString("format", "", "json|yaml").
			SetAction(func(cmd *commander.Command) bool {

				abi.GetAppVerConfig(cmd.String("appid"), cmd.String("ver"), cmd.String("format"))

				return true
			})

		s.
			SubCommand("approve").
			SetString("appid", "", "App ID").
			SetString("id", "", "Container ID").
			SetAction(func(cmd *commander.Command) bool {

				abi.Approve(cmd.String("appid"), cmd.String("id"))

				return true
			})

		s.
			SubCommand("unapprove").
			SetString("appid", "", "App ID").
			SetString("id", "", "Container ID").
			SetAction(func(cmd *commander.Command) bool {

				abi.Unapprove(cmd.String("appid"), cmd.String("id"))

				return true
			})

		s.
			SubCommand("publish").
			SetString("file", "", "").
			SetString("ver", "", "").
			SetString("number", "", "").
			SetAction(func(cmd *commander.Command) bool {

				abi.Publish(cmd.String("file"), cmd.String("ver"), cmd.String("number"))

				return true
			})
	}

	{
		s := app.SubCommand("env")
		s.
			SubCommand("os").
			SetAction(func(cmd *commander.Command) bool {

				fmt.Printf("%s", runtime.GOOS)

				return true
			})
		s.
			SubCommand("arch").
			SetAction(func(cmd *commander.Command) bool {

				fmt.Printf("%s", runtime.GOARCH)

				return true
			})
	}

	app.Parse(os.Args[1:])
}
