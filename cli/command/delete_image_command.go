package command

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/cli/prompt"
	"github.com/imc-trading/dock2box/client"
)

func NewDeleteImageCommand() cli.Command {
	return cli.Command{
		Name:  "image",
		Usage: "Delete image",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			deleteImageCommandFunc(c)
		},
	}
}

func deleteImageCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a image id")
	}
	id := c.Args()[0]

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	if !prompt.Bool("Are you sure you wan't to remove image id: "+id, true) {
		os.Exit(1)
	}

	s, err := clnt.Image.Delete(id)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%v\n", string(s.JSON()))
}
