package main

import (
	"github-activity/cmds"
	"log"
	"os"

	"github.com/hashicorp/cli"
)

func usernameCommandFactory() (cli.Command, error) {
	ui := cmds.Username{
		UI: cli.ColoredUi{
			OutputColor: cli.UiColorGreen,
			InfoColor:   cli.UiColorBlue,
			ErrorColor:  cli.UiColorRed,
			WarnColor:   cli.UiColorYellow,
			Ui:          &cli.BasicUi{Writer: os.Stdout},
		},
	}
	return &ui, nil
}

func main() {
	c := cli.NewCLI("learnCli", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"username": usernameCommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
