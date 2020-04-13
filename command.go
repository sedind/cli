package cli

import (
	"flag"
	"fmt"
	"io"
)

// Command object for cli app
type Command struct {
	// The name of the command
	Name string
	// A list of aliases for the command
	Aliases []string
	// A longer explanation of how the command works
	Description string
	// An action to execute before any subcommands are run, but after the context is ready
	// If a non-nil error is returned, no subcommands are run
	Before BeforeFunc
	// An action to execute after any subcommands are run, but after the subcommand has finished
	// It is run even if ActionFn() panics
	After AfterFunc
	// The action to execute when no subcommands are specified
	Action ActionFunc
	// List of flags to parse
	Flags []Flag
	// List of child commands
	Commands []*Command

	flagSet *flag.FlagSet
}

// Commands collection
type Commands []*Command

// Len -
func (c Commands) Len() int {
	return len(c)
}

// Less -
func (c Commands) Less(i, j int) bool {
	return lexicographicLess(c[i].Name, c[j].Name)
}

// Swap -
func (c Commands) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// FlagSet returns Command flagSet
func (c *Command) FlagSet() (*flag.FlagSet, error) {
	if c.flagSet == nil {
		set := flag.NewFlagSet(c.Name, flag.ContinueOnError)

		for _, f := range c.Flags {
			if err := f.Apply(set); err != nil {
				return nil, err
			}
		}
		c.flagSet = set
	}

	return c.flagSet, nil
}

// Names returns the names including short names and aliases.
func (c *Command) Names() []string {
	return append([]string{c.Name}, c.Aliases...)
}

// HasName returns true if Command.Name matches given name
func (c *Command) HasName(name string) bool {
	for _, n := range c.Names() {
		if n == name {
			return true
		}
	}
	return false
}

// Run invokes the command given the context, parses ctx.Args() to generate command-specific flags
func (c *Command) Run(ctx *Context) (err error) {
	fmt.Printf("Executing command %s\n", c.Name)
	return nil
}

// Usage prints command Usage instructions
func (c *Command) Usage(w io.Writer) {
	_, _ = fmt.Fprintf(w, "\n%s\t%s\n", c.Name, c.Description)

	_, _ = fmt.Fprint(w, "  Usage:  ")
	Usage(w, c.Name, c.Flags)
	for _, cmd := range c.Commands {
		cmd.Usage(w)
	}
}
