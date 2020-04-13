package cli

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

// GenericFlag is a flag with type interface
type GenericFlag struct {
	Name        string
	Aliases     []string
	Description string
	Required    bool
	Value       interface{}
	Destination interface{}
}

// IsRequired returns whether or not the flag is required
func (f *GenericFlag) IsRequired() bool {
	return f.Required
}

// Names returns the names of the flag
func (f *GenericFlag) Names() []string {
	var ret []string

	for _, part := range append([]string{f.Name}, f.Aliases...) {
		ret = append(ret, commaWhitespace.ReplaceAllString(part, ""))
	}

	return ret
}

// Apply populates the flag given the flag set
func (f *GenericFlag) Apply(set *flag.FlagSet) error {
	var err error
	for _, name := range f.Names() {
		if err = f.apply(set, name); err != nil {
			break
		}
	}

	return err
}

func (f *GenericFlag) apply(set *flag.FlagSet, name string) error {
	var err error
	switch v := f.Value.(type) {
	case bool:
		if dest, ok := f.Destination.(*bool); ok {
			set.BoolVar(dest, name, v, f.Description)
		} else {
			set.Bool(name, v, f.Description)
		}
	case time.Duration:
		if dest, ok := f.Destination.(*time.Duration); ok {
			set.DurationVar(dest, name, v, f.Description)
		} else {
			set.Duration(name, v, f.Description)
		}
	case float64:
		if dest, ok := f.Destination.(*float64); ok {
			set.Float64Var(dest, name, v, f.Description)
		} else {
			set.Float64(name, v, f.Description)
		}
	case int:
		if dest, ok := f.Destination.(*int); ok {
			set.IntVar(dest, name, v, f.Description)
		} else {
			set.Int(name, v, f.Description)
		}
	case int64:
		if dest, ok := f.Destination.(*int64); ok {
			set.Int64Var(dest, name, v, f.Description)
		} else {
			set.Int64(name, v, f.Description)
		}
	case string:
		if dest, ok := f.Destination.(*string); ok {
			set.StringVar(dest, name, v, f.Description)
		} else {
			set.String(name, v, f.Description)
		}
	case uint:
		if dest, ok := f.Destination.(*uint); ok {
			set.UintVar(dest, name, v, f.Description)
		} else {
			set.Uint(name, v, f.Description)
		}
	case uint64:
		if dest, ok := f.Destination.(*uint64); ok {
			set.Uint64Var(dest, name, v, f.Description)
		} else {
			set.Uint64(name, v, f.Description)
		}
	default:
		if v == nil {
			err = fmt.Errorf("flag `%s` must have value provided", f.Name)
		} else {
			err = fmt.Errorf("flag `%s` has unsupported type: %T", f.Name, v)
		}

	}

	return err

}

// Describe describes flag usage
func (f *GenericFlag) Describe() string {
	defaultVal := ""
	switch v := f.Value.(type) {
	case bool:
		defaultVal = fmt.Sprintf("%t", v)
	case time.Duration:
		defaultVal = fmt.Sprintf("%s", v)
	default:
		defaultVal = fmt.Sprintf("%v", v)
	}
	return fmt.Sprintf("%s\t%s\t%s (%s)", f.Name, strings.Join(f.Aliases, ","), f.Description, defaultVal)
}

// NewFlag creates Generic flag
func NewFlag(name string, aliases []string, description string, required bool, value interface{}) *GenericFlag {
	return &GenericFlag{
		Name:        name,
		Aliases:     aliases,
		Description: description,
		Required:    required,
		Value:       value,
	}
}

// NewFlagVar creates Generic flag with destination variable
func NewFlagVar(name string, aliases []string, description string, required bool, value interface{}, destination interface{}) *GenericFlag {
	return &GenericFlag{
		Name:        name,
		Aliases:     aliases,
		Description: description,
		Required:    required,
		Value:       value,
		Destination: destination,
	}
}
