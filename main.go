package main

import (
	"fmt"
	"os"

	"github.com/takumakume/dependency-track-policy-applier/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
