// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package main

import (
	"log"
	"os"

	App "github.com/axetroy/wsm"
	app2 "github.com/axetroy/wsm/internal/app"
	"github.com/axetroy/wsm/internal/library/daemon"
	"github.com/axetroy/wsm/internal/library/util"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Usage = "user server controller"
	app.Authors = []*cli.Author{
		{
			Name:  App.Author,
			Email: App.Email,
		},
	}
	app.Version = App.Version
	cli.AppHelpTemplate = App.CliTemplate

	app.Commands = []*cli.Command{
		{
			Name:  "start",
			Usage: "start user server",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "daemon",
					Aliases: []string{"d"},
					Usage:   "running in daemon mode",
				},
			},
			Action: func(c *cli.Context) error {
				// 判断当其是否是子进程，当父进程return之后，子进程会被系统1号进程接管
				return daemon.Start(app2.Serve, c.Bool("daemon"))
			},
		},
		{
			Name:  "stop",
			Usage: "stop user server",
			Action: func(c *cli.Context) error {
				return daemon.Stop()
			},
		},
		{
			Name:  "env",
			Usage: "print runtime environment",
			Action: func(c *cli.Context) error {
				util.PrintEnv()
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
