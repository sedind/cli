package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Usage prints usage documentation for command
func Usage(w io.Writer, cmdName string, flags []Flag) {
	base := filepath.Base(os.Args[0])
	pre := fmt.Sprintf("%s ", base)
	if cmdName == base {
		pre = ""
	}
	fmt.Fprintf(w, "%s%s [flags] [args...]\n", pre, cmdName)
	for _, f := range flags {
		fmt.Fprintf(w, "\t-%s\n", f.Describe())
	}

}
