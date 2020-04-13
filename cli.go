// Package cli provides tools for creating and organizing command line
// Go applications. cli is designed to be easy to understand and write, the most simple
// cli application can be written as follows:
//   func main() {
//     app := &cli.App{
//			 Name: "greet",
//			 Description: "say a greeting",
//			 Action: func(c *cli.Context) error {
//				 fmt.Println("Greetings")
//				 return nil
//			 },
//		 }
//
//     app.Run(os.Args)
//   }
package cli

// BeforeFunc is an action to execute before any subcommands are run, but after
// the context is ready if a non-nil error is returned, no subcommands are run
type BeforeFunc func(*Context) error

// AfterFunc is an action to execute after any subcommands are run, but after the
// subcommand has finished it is run even if Action() panics
type AfterFunc func(*Context) error

// ActionFunc is the action to execute when no subcommands are specified
type ActionFunc func(*Context) error
