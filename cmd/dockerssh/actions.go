package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
	"github.com/xwjdsh/dockerssh"
)

var (
	nameRegex = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_.-]+$")
)

func create(c *cli.Context) error {
	err := dockerssh.ArgsCountCheck(c.NArg(), 1, 1)
	if err != nil {
		fmt.Printf("%s %v\n", promptui.IconBad, err)
		return nil
	}
	options := &dockerssh.Options{
		Name:    c.Args().First(),
		Verbose: c.Bool("verbose"),
	}
	if !nameRegex.MatchString(options.Name) {
		fmt.Printf("%s %s %s\n", promptui.IconBad, "Name:", "Invalid name, only [a-zA-Z0-9][a-zA-Z0-9_.-] are allowed")
		return nil
	}
	fmt.Printf("%s %s %s\n", promptui.IconGood, "Name:", options.Name)

	options.Port, err = getValue("Port", func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("Invalid local port")
		}
		return nil
	}, "3000")
	if err != nil {
		return err
	}
	options.Volume, err = getValue("Volume", nil, fmt.Sprintf("./%s", options.Name))
	if err != nil {
		return err
	}
	options.Volume, err = filepath.Abs(options.Volume)
	if err != nil {
		return err
	}
	os.MkdirAll(options.Volume, os.ModePerm)

	err = dockerssh.Create(options)
	if err == nil {
		fmt.Printf("%s %s\n", promptui.IconGood, "Success")
	} else {
		fmt.Printf("%s %s\n", promptui.IconBad, err.Error())
	}
	return nil
}

func list(c *cli.Context) error {
	//clients, err := kcm.List()
	//if err != nil {
	//fmt.Printf("%s %s\n", bold(promptui.IconBad), bold("Failed:"), err.Error())
	//return nil
	//}
	//maxLen := 5
	//for _, client := range clients {
	//if l := len(client["name"]); l > maxLen {
	//maxLen = l
	//}
	//}
	//format := "%-" + strconv.Itoa(maxLen+5) + "s\t%-10s\t%s\n"
	//fmt.Printf(format, "Name", "State", "Service")
	//for _, client := range clients {
	//fmt.Printf(format, client["name"], client["state"], client["service"])
	//}
	return nil
}

func destroy(c *cli.Context) error {
	return nil
}
