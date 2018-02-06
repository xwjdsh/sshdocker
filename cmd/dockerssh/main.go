package main

import (
	"os"

	"github.com/urfave/cli"
)

var version = "master"

func main() {
	app := cli.NewApp()
	app.Usage = "docker sshd CLI tool"
	app.Version = version
	app.Commands = commands()
	app.Run(os.Args)
}

func commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "create",
			Usage:  "Create a new docker sshd shell",
			Action: create,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "verbose, v",
					Usage: "Show verbose log",
				},
			},

			Aliases: []string{"c"},
		},
		{
			Name:    "list",
			Usage:   "List all docker sshd shell",
			Action:  list,
			Aliases: []string{"l"},
		},
		{
			Name:    "destroy",
			Usage:   "Destroy one or more docker sshd shell",
			Action:  destroy,
			Aliases: []string{"d"},
		},
	}
}
