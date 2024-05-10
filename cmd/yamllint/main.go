package main

import (
	"os"

	"github.com/wasilibs/go-yamllint/internal/pysite"
	"github.com/wasilibs/go-yamllint/internal/runner"
)

func main() {
	os.Exit(runner.Run("yamllint", os.Args[1:], pysite.Python, os.Stdin, os.Stdout, os.Stderr, "."))
}
