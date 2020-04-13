package cli

import (
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"
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
	Flags Flags
	// List of child commands
	Commands Commands

	parent *Command
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
	return c.Flags.FlagSet(c.Name, flag.ContinueOnError)
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
	// sort Flags
	sort.Sort(c.Flags)
	// sort commands
	sort.Sort(c.Commands)

	flagSet, err := c.FlagSet()
	if err != nil {
		return err
	}

	flagSet.Usage = func() {
		c.Usage(ctx.App.Writer)
	}

	args := ctx.Args().Slice()
	// parse arguments
	flagSet.Parse(args[1:])

	// check if app has required flags
	errRf := checkRequiredFlags(c.Flags, flagSet)
	if errRf != nil {
		fmt.Fprintf(ctx.App.ErrWriter, "\n%s\n", errRf.Error())
		c.Usage(ctx.App.Writer)
		return nil
	}

	appCtx := newContext(ctx.App, flagSet, ctx)
	if c.After != nil {
		defer func() {
			if err := c.After(appCtx); err != nil {
				_, _ = fmt.Fprintf(ctx.App.ErrWriter, "%v\n\n", err)
			}
		}()
	}

	if c.Before != nil {
		if err = c.Before(appCtx); err != nil {
			_, _ = fmt.Fprintf(ctx.App.ErrWriter, "%v\n\n", err)
			return err
		}
	}

	cArgs := appCtx.Args()
	if cArgs.Present() {
		name := cArgs.First()
		cmd := c.SubCommand(name)
		if cmd != nil {
			cmd.parent = c
			return cmd.Run(appCtx)
		}
	}

	if c.Action != nil {
		return c.Action(appCtx)
	}

	// no command has been executed, show usage
	c.Usage(ctx.App.Writer)

	return nil
}

// SubCommand returns the named sub command on Commans. Returns nil if the command does not exist
func (c *Command) SubCommand(name string) *Command {
	for _, c := range c.Commands {
		if c.HasName(name) {
			return c
		}
	}

	return nil
}

// CommandPathName returns full command name in command tree
func (c *Command) CommandPathName() string {
	prefix := ""
	if c.parent != nil {
		prefix = c.parent.CommandPathName()
	}
	return fmt.Sprintf("%s %s", prefix, c.Name)
}

// Usage prints command Usage instructions
func (c *Command) Usage(w io.Writer) {
	_, _ = fmt.Fprintf(w, "\n%s\t%s\n", c.Name, c.Description)

	_, _ = fmt.Fprint(w, "  Usage:  ")
	name := strings.TrimPrefix(c.CommandPathName(), " ")
	Usage(w, name, c.Flags)
	for _, cmd := range c.Commands {
		cmd.parent = c
		cmd.Usage(w)
	}
}
