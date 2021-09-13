package main

import (
	"github.com/ory/cli/cmd"
	"github.com/ory/x/profilex"
)

func main() {
	defer profilex.Profile().Stop()
	cmd.Execute()
}
