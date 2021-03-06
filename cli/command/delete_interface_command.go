package command

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/cli/prompt"
	"github.com/imc-trading/dock2box/client"
)

func NewDeleteInterfaceCommand() cli.Command {
	return cli.Command{
		Name:  "interface",
		Usage: "Delete interface",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			deleteInterfaceCommandFunc(c)
		},
	}
}

func deleteInterfaceCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify an interface id")
	}
	id := c.Args()[0]

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	if !prompt.Bool("Are you sure you wan't to remove interface id: "+id, true) {
		os.Exit(1)
	}

	s, err := clnt.Interface.Delete(id)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%v\n", string(s.JSON()))
}
