package cli

import (
	"context"
	"flag"
)

// Context -
type Context struct {
	context.Context
	parentContext *Context
	App           *App
	flagSet       *flag.FlagSet
}

// newContext creates a new context. For use in when invoking a Command action.
func newContext(app *App, set *flag.FlagSet, parentCtx *Context) *Context {
	c := &Context{App: app, flagSet: set, parentContext: parentCtx}
	if parentCtx != nil {
		c.Context = parentCtx.Context
		if parentCtx.flagSet == nil {
			parentCtx.flagSet = &flag.FlagSet{}
		}
	}
	if c.Context == nil {
		c.Context = context.Background()
	}
	return c
}

// NumFlags returns the number of flags set
func (c *Context) NumFlags() int {
	return c.flagSet.NFlag()
}

// Args returns the command line arguments associated with the context.
func (c *Context) Args() Args {
	ret := args(c.flagSet.Args())
	return &ret
}

// Set sets a context flag to a value.
func (c *Context) Set(name, value string) error {
	return c.flagSet.Set(name, value)
}

// Value returns the value of the flag corresponding to `name`
func (c *Context) Value(name string) interface{} {
	return c.flagSet.Lookup(name).Value.(flag.Getter).Get()
}

// NArg returns the number of the command line arguments.
func (c *Context) NArg() int {
	return c.Args().Len()
}
