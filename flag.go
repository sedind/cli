package cli

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	slPfx = fmt.Sprintf("sl:::%d:::", time.Now().UTC().UnixNano())

	commaWhitespace = regexp.MustCompile("[, ]+.*")
)

// Flag is an interface for parsing flags in cli.
// For more advanced flag parsing techniques, it is recommended that
// this interface be implemented.
type Flag interface {
	// Apply Flag settings to the given flag set
	Apply(*flag.FlagSet) error
	Names() []string
	IsRequired() bool
	Describe() string
}

// Flags is a slice of Flag.
type Flags []Flag

// FlagSet returns flag set with the specified name and
// error handling property. If the name is not empty, it will be printed
// in the default usage message and in error messages.
func (f Flags) FlagSet(name string, errHandler flag.ErrorHandling) (*flag.FlagSet, error) {
	set := flag.NewFlagSet(name, errHandler)
	for _, flag := range f {
		if err := flag.Apply(set); err != nil {
			return nil, err
		}
	}
	return set, nil
}

// Len returns lenght of a flags slice
func (f Flags) Len() int {
	return len(f)
}

// Less -
func (f Flags) Less(i, j int) bool {
	if len(f[j].Names()) == 0 {
		return false
	} else if len(f[i].Names()) == 0 {
		return true
	}
	return lexicographicLess(f[i].Names()[0], f[j].Names()[0])
}

// Swap -
func (f Flags) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

type errRequiredFlags struct {
	missingFlags []string
}

func (e *errRequiredFlags) Error() string {
	numberOfMissingFlags := len(e.missingFlags)
	if numberOfMissingFlags == 1 {
		return fmt.Sprintf("Required flag %q not set", e.missingFlags[0])
	}
	joinedMissingFlags := strings.Join(e.missingFlags, ", ")
	return fmt.Sprintf("Required flags %q not set", joinedMissingFlags)
}

func (e *errRequiredFlags) getMissingFlags() []string {
	return e.missingFlags
}

func checkRequiredFlags(flags []Flag, set *flag.FlagSet) *errRequiredFlags {
	var missingFlags []string
	seen := make(map[string]bool)
	set.Visit(func(f *flag.Flag) {
		seen[f.Name] = true
	})

	for _, f := range flags {
		if f.IsRequired() {
			name := f.Names()[0]
			if !seen[name] {
				missingFlags = append(missingFlags, name)
			}
		}
	}

	if len(missingFlags) > 0 {
		return &errRequiredFlags{missingFlags: missingFlags}
	}

	return nil
}
