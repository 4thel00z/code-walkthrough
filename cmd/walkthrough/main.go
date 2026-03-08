package main

import (
	"fmt"
	"os"

	"github.com/tahrioui/code-walkthrough/adapter"
	"github.com/tahrioui/code-walkthrough/skilldata"
)

func main() {
	if err := adapter.NewRootCmd(skilldata.SkillMD, skilldata.SchemaJSON).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
