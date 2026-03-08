package main

import (
	"fmt"
	"os"

	"github.com/4thel00z/code-walkthrough/adapter"
	"github.com/4thel00z/code-walkthrough/skilldata"
)

func main() {
	if err := adapter.NewRootCmd(skilldata.SkillMD, skilldata.SchemaJSON).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
