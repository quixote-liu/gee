package gee

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var DefaultWriter io.Writer = os.Stdout

func debugPrint(format string, values ...interface{}) {
	if !strings.HasPrefix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(DefaultWriter, "[GEE-debug] "+format, values...)
}