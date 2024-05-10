package main

import (
	"github.com/goyek/x/boot"

	"github.com/wasilibs/tools/tasks"
)

func main() {
	tasks.Define(tasks.Params{
		LibraryName: "yamllint",
		LibraryRepo: "adrienverge/yamllint",
		GoReleaser:  true,
	})
	boot.Main()
}
