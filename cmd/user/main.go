// Copyright 2019 Axetroy. All rights reserved. MIT license.
package main

import (
	"log"
	"os"

	App "github.com/axetroy/terminal"
	"github.com/axetroy/terminal/core/helper/daemon"
	"github.com/axetroy/terminal/core/server/user_server"
	"github.com/axetroy/terminal/core/util"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Usage = "user server controller"
	app.Author = App.Author
	app.Email = App.Email
	app.Version = App.Version
	cli.AppHelpTemplate = App.CliTemplate

	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start user server",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "daemon, d",
					Usage: "running in daemon mode",
				},
			},
			Action: func(c *cli.Context) error {
				// 判断当其是否是子进程，当父进程return之后，子进程会被系统1号进程接管
				return daemon.Start(user_server.Serve, c.Bool("daemon"))
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
