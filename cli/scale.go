package cli

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func showFlags(c *cli.Context) {
	fmt.Println("cattle connect to: ", c.String("H"))
	fmt.Println("ENVS: ")
	for _, env := range c.StringSlice("e") {
		fmt.Println(env)
	}
	fmt.Println("labels: ")
	for _, lalel := range c.StringSlice("l") {
		fmt.Println(lalel)
	}
	fmt.Println("filters: ")
	for _, filter := range c.StringSlice("f") {
		fmt.Println(filter)
	}
	fmt.Println("numbers: ", c.Int("n"))
}

func scale(c *cli.Context) {
	showFlags(c)
}
