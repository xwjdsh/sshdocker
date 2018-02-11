package main

import (
	"os"

	"github.com/urfave/cli"
)

var version = "master"

func main() {
	app := cli.NewApp()
	app.Usage = "sshd docker shell CLI tool"
	app.Version = version
	app.Commands = commands()
	app.Run(os.Args)
}

func commands() []cli.Command {
	return []cli.Command{
		{
			Name:    "create",
			Usage:   "Create a new sshd docker shell",
			Action:  create,
			Aliases: []string{"c"},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "verbose, v", Usage: "Show verbose log",
				},
			},
		},
		{
			Name:    "ls",
			Usage:   "List all sshd docker shells",
			Action:  list,
			Aliases: []string{"l"},
		},
		{
			Name:    "destroy",
			Usage:   "Destroy one or more sshd docker shells",
			Action:  destroy,
			Aliases: []string{"d"},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "volume, v", Usage: "Remove volume",
				},
			},
		},
	}
}
