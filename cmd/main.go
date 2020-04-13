package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sedind/cli"
)

func main() {

	app := cli.New()
	app.Description = "A sample application"
	app.UseHelpFlag = true
	app.UseVersionFlag = true
	app.Commands = cli.Commands{
		&cli.Command{
			Name: "init",
			Action: func(ctx *cli.Context) error {
				fmt.Println("init Action")
				return nil
			},
			Commands: cli.Commands{
				&cli.Command{
					Name: "project",
					Action: func(ctx *cli.Context) error {
						fmt.Println("Project Action")
						return nil
					},
					Commands: cli.Commands{
						&cli.Command{
							Name: "test1",
							Action: func(ctx *cli.Context) error {
								fmt.Println("test1 Action")
								return nil
							},
							Commands: cli.Commands{
								&cli.Command{
									Name:        "test2",
									Description: "Test 2 description",
									Flags: cli.Flags{
										cli.NewFlag("flag1", []string{"f1"}, "Flag 1 used for nothing", false, "nothing"),
									},
									Action: func(ctx *cli.Context) error {
										val := ctx.Value("flag1")
										fmt.Printf("flag1 value: %v\n", val)
										fmt.Println("test2 Action")
										return nil
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
