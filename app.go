package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// App represents cli application
type App struct {
	// The name of the program. Defaults to path.Base(os.Args[0])
	Name string
	// Application version
	Version string
	// Application desctiption
	Description string
	// An action to execute before any subcommands are run, but after the context is ready
	// If a non-nil error is returned, no subcommands are run
	Before BeforeFunc
	// An action to execute after any subcommands are run, but after the subcommand has finished
	// It is run even if ActionFn() panics
	After AfterFunc
	// The action to execute when no commands are specified
	Action ActionFunc
	// Application metadata
	Metadata map[string]interface{}
	//Flags represents a set of defined flags.
	Flags Flags
	// List of commands to execute
	Commands Commands
	// Writer writer to write output to
	Writer io.Writer
	// ErrWriter writes error output
	ErrWriter io.Writer
	// UseHelpFlag determines if built-in help flag will be used
	UseHelpFlag bool
	// UseHelpFlag determines if built-in version flag will be used
	UseVersionFlag bool
}

// New creates a new CLI application
func New() *App {
	return &App{
		Name:        filepath.Base(os.Args[0]),
		Description: "A new cli application",
		Version:     "v0.0.0",
		Metadata:    make(map[string]interface{}),
		Writer:      os.Stdout,
		ErrWriter:   os.Stderr,
	}
}

// Run is the entry point to the cli app. Parses the args slice and routes
// to the proper flag/args combination
func (a *App) Run(args []string) (err error) {
	return a.RunContext(context.Background(), args)
}

// RunContext is like Run except it takes a Context that will be
// passed to its commands and sub-commands.
func (a *App) RunContext(ctx context.Context, args []string) error {
	var helpFlagVal bool
	var versionFlagVal bool
	if a.UseHelpFlag {
		a.Flags = append(a.Flags, &GenericFlag{
			Name:        "help",
			Aliases:     []string{"h"},
			Description: "Prints application help",
			Value:       false,
			Destination: &helpFlagVal,
		})
	}
	if a.UseVersionFlag {
		a.Flags = append(a.Flags, &GenericFlag{
			Name:        "version",
			Aliases:     []string{"v"},
			Description: "Prints application version",
			Value:       false,
			Destination: &versionFlagVal,
		})
	}

	// sort Flags
	sort.Sort(a.Flags)
	// sort commands
	sort.Sort(a.Commands)

	flagSet, err := a.Flags.FlagSet(a.Name, flag.ExitOnError)
	if err != nil {
		return err
	}

	flagSet.Usage = func() {
		a.Usage(a.Writer)
	}

	// parse arguments
	flagSet.Parse(args[1:])

	// check if app has required flags
	errRf := checkRequiredFlags(a.Flags, flagSet)
	if errRf != nil {
		fmt.Fprintf(a.ErrWriter, "\n%s\n", errRf.Error())
		a.Usage(a.Writer)
		return nil
	}

	appCtx := newContext(a, flagSet, &Context{Context: ctx})
	if a.After != nil {
		defer func() {
			if err := a.After(appCtx); err != nil {
				_, _ = fmt.Fprintf(a.ErrWriter, "%v\n\n", err)
			}
		}()
	}

	if a.Before != nil {
		if err = a.Before(appCtx); err != nil {
			_, _ = fmt.Fprintf(a.ErrWriter, "%v\n\n", err)
			return err
		}
	}

	cArgs := appCtx.Args()
	if cArgs.Present() {
		name := cArgs.First()
		cmd := a.Command(name)
		if cmd != nil {
			return cmd.Run(appCtx)
		}
	}

	// check if build in flag was used
	if helpFlagVal {
		a.Usage(a.Writer)
		return nil
	}

	// check if built in version vlag was used
	if versionFlagVal {
		a.PrintVersion(a.Writer)
		return nil
	}

	if a.Action != nil {
		return a.Action(appCtx)
	}

	// no command has been executed, show usage
	a.Usage(a.Writer)

	return nil
}

// Command returns the named command on App. Returns nil if the command does not exist
func (a *App) Command(name string) *Command {
	for _, c := range a.Commands {
		if c.HasName(name) {
			return c
		}
	}

	return nil
}

// Usage prints application usage instructions
func (a *App) Usage(w io.Writer) {
	_, _ = fmt.Fprintf(w, "%s\n", a.Description)
	_, _ = fmt.Fprintf(w, "\nUsage:\t")
	Usage(w, a.Name, a.Flags)
	for _, cmd := range a.Commands {
		cmd.Usage(a.ErrWriter)
	}
}

// PrintVersion to writer
func (a *App) PrintVersion(w io.Writer) {
	_, _ = fmt.Fprintf(w, "%s version: %s\n", a.Name, a.Version)
}
