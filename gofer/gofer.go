package main

import (
	"fmt"
	"os"
)

import (
	"github.com/chuckpreslar/gofer"
)

const VERSION = "0.0.4"

func main() {
	var definition string

	if 1 >= len(os.Args) {
		definition = ""
	} else {
		definition = os.Args[1]
	}

	if "version" == definition {
		fmt.Fprintf(os.Stdout, "%s\n", VERSION)
		os.Exit(0)
	}

	var arguments []string

	if 1 < len(os.Args) {
		arguments = os.Args[2:]
	}

	gofer.LoadAndPerform(definition, arguments...)
}
