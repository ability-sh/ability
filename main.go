package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/ability-sh/ability/abi"
	"github.com/ability-sh/ability/commander"
)

func main() {

	app := commander.NewCommand("ability")

	app.SetAction(func(cmd *commander.Command, args []string, usage bool) bool {

		fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

		fs_help := fs.Bool("help", false, "")
		fs_registry := fs.String("registry", "https://ac.ability.sh", "")
		fs_token := fs.String("token", "", "")

		if usage {
			fs.Usage()
			return false
		}

		fs.Parse(args)

		if *fs_help || len(args) == 0 {
			cmd.Run(args, true)
			return true
		}

		abi.SetRegistry(abi.NewACRegistry(*fs_registry))

		if *fs_token != "" {
			abi.GetRegistry().SetToken(*fs_token)
		}

		return false
	})

	app.
		SubCommand("login").
		SetAction(func(cmd *commander.Command, args []string, usage bool) bool {

			fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

			if usage {
				fs.Usage()
				return false
			}

			abi.Login()

			return false
		})

	app.
		SubCommand("logout").
		SetAction(func(cmd *commander.Command, args []string, usage bool) bool {

			fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

			if usage {
				fs.Usage()
				return false
			}

			fs.Parse(args)

			abi.Logout()

			return false
		})

	app.
		SubCommand("create").
		SetAction(func(cmd *commander.Command, args []string, usage bool) bool {

			fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

			fs_json := fs.String("json", "", "")
			fs_yaml := fs.String("yaml", "", "")
			fs_file := fs.String("file", "", "")

			if usage {
				fs.Usage()
				return false
			}

			fs.Parse(args)

			abi.Create(*fs_json, *fs_yaml, *fs_file)

			return false
		})

	app.
		SubCommand("setsecret").
		SetAction(func(cmd *commander.Command, args []string, usage bool) bool {

			fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

			fs_id := fs.String("id", "", "Container ID")

			if usage {
				fs.Usage()
				return false
			}

			fs.Parse(args)

			abi.SetSecret(*fs_id)

			return false
		})

	app.
		SubCommand("setconfig").
		SetAction(func(cmd *commander.Command, args []string, usage bool) bool {

			fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

			fs_id := fs.String("id", "", "Container ID")
			fs_json := fs.String("json", "", "")
			fs_yaml := fs.String("yaml", "", "")
			fs_file := fs.String("file", "", "")

			if usage {
				fs.Usage()
				return false
			}

			fs.Parse(args)

			abi.SetConfig(*fs_id, *fs_json, *fs_yaml, *fs_file)

			return false
		})

	app.
		SubCommand("getconfig").
		SetAction(func(cmd *commander.Command, args []string, usage bool) bool {

			fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

			fs_id := fs.String("id", "", "Container ID")
			fs_format := fs.String("format", "", "json|yaml")

			if usage {
				fs.Usage()
				return false
			}

			fs.Parse(args)

			abi.GetConfig(*fs_id, *fs_format)

			return false
		})

	{
		s := app.SubCommand("app")

		s.
			SubCommand("create").
			SetAction(func(cmd *commander.Command, args []string, usage bool) bool {

				fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

				fs_json := fs.String("json", "", "")
				fs_yaml := fs.String("yaml", "", "")
				fs_file := fs.String("file", "", "")

				if usage {
					fs.Usage()
					return false
				}

				fs.Parse(args)

				abi.CreateApp(*fs_json, *fs_yaml, *fs_file)

				return false
			})

		s.
			SubCommand("setconfig").
			SetAction(func(cmd *commander.Command, args []string, usage bool) bool {

				fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

				fs_appid := fs.String("appid", "", "App ID")
				fs_json := fs.String("json", "", "")
				fs_yaml := fs.String("yaml", "", "")
				fs_file := fs.String("file", "", "")

				if usage {
					fs.Usage()
					return false
				}

				fs.Parse(args)

				abi.SetAppConfig(*fs_appid, *fs_json, *fs_yaml, *fs_file)

				return false
			})

		s.
			SubCommand("getconfig").
			SetAction(func(cmd *commander.Command, args []string, usage bool) bool {

				fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

				fs_appid := fs.String("appid", "", "App ID")
				fs_format := fs.String("format", "", "json|yaml")

				if usage {
					fs.Usage()
					return false
				}

				fs.Parse(args)

				abi.GetAppConfig(*fs_appid, *fs_format)

				return false
			})

		s.
			SubCommand("getver").
			SetAction(func(cmd *commander.Command, args []string, usage bool) bool {

				fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

				fs_appid := fs.String("appid", "", "App ID")
				fs_ver := fs.String("ver", "", "App Ver")
				fs_format := fs.String("format", "", "json|yaml")

				if usage {
					fs.Usage()
					return false
				}

				fs.Parse(args)

				abi.GetAppVerConfig(*fs_appid, *fs_ver, *fs_format)

				return false
			})

		s.
			SubCommand("approve").
			SetAction(func(cmd *commander.Command, args []string, usage bool) bool {

				fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

				fs_appid := fs.String("appid", "", "App ID")
				fs_id := fs.String("id", "", "Container ID")

				if usage {
					fs.Usage()
					return false
				}

				fs.Parse(args)

				abi.Approve(*fs_appid, *fs_id)

				return false
			})

		s.
			SubCommand("unapprove").
			SetAction(func(cmd *commander.Command, args []string, usage bool) bool {

				fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

				fs_appid := fs.String("appid", "", "App ID")
				fs_id := fs.String("id", "", "Container ID")

				if usage {
					fs.Usage()
					return false
				}

				fs.Parse(args)

				abi.Unapprove(*fs_appid, *fs_id)

				return false
			})

		s.
			SubCommand("publish").
			SetAction(func(cmd *commander.Command, args []string, usage bool) bool {

				fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

				fs_file := fs.String("file", "", "")
				fs_ver := fs.String("ver", "", "")

				if usage {
					fs.Usage()
					return false
				}

				fs.Parse(args)

				abi.Publish(*fs_file, *fs_ver)

				return false
			})
	}

	{
		s := app.SubCommand("env")
		s.
			SubCommand("os").
			SetAction(func(cmd *commander.Command, args []string, usage bool) bool {

				fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

				if usage {
					fs.Usage()
					return false
				}

				fs.Parse(args)

				fmt.Printf("%s", runtime.GOOS)

				return false
			})
		s.
			SubCommand("arch").
			SetAction(func(cmd *commander.Command, args []string, usage bool) bool {

				fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

				if usage {
					fs.Usage()
					return false
				}

				fs.Parse(args)

				fmt.Printf("%s", runtime.GOARCH)

				return false
			})
	}

	app.Run(os.Args[1:], false)
}
