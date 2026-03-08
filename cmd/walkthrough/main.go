package main

import (
	"fmt"
	"os"

	"github.com/tahrioui/code-walkthrough/adapter"
)

func main() {
	if err := adapter.NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
