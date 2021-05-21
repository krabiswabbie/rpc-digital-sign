package main

import (
	"os"

	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli"
)

const addr = "/tmp/ping.sock"

func main() {
	log.Root().SetHandler(log.StreamHandler(os.Stderr, log.LogfmtFormat()))

	app := cli.NewApp()
	app.Usage = "ping-pong ether app"
	app.Commands = []cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "launch server",
			Action: func(c *cli.Context) error {
				return newServer(addr)
			},
		},
		{
			Name:    "client",
			Aliases: []string{"c"},
			Usage:   "launch client",
			Action: func(c *cli.Context) error {
				return newClient(addr)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		println(err.Error())
	}
}
