package command

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/codegangsta/cli"

	"github.com/imc-trading/dock2box/client"
)

func NewQueryInterfaceCommand() cli.Command {
	return cli.Command{
		Name:  "interfaces",
		Usage: "Query interfaces",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) {
			queryInterfaceCommandFunc(c)
		},
	}
}

func queryInterfaceCommandFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatal("You need to specify a comma-separated list of key/value pairs, percent sign (%) can be used as a wildcard")
	}
	cond := map[string]string{}
	for _, e := range strings.Split(c.Args()[0], ",") {
		l := strings.Split(e, "=")
		cond[l[0]] = l[1]
	}

	clnt := client.New(c.GlobalString("server"))
	if c.GlobalBool("debug") {
		clnt.SetDebug()
	}

	h, err := clnt.Interface.Query(cond)
	if err != nil {
		log.Fatal(err.Error())
	}

	b, _ := json.MarshalIndent(h, "", "  ")
	fmt.Printf("%v\n", string(b))
}
