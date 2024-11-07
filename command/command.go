package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{}

func AddCommand(cmds ...*cobra.Command) {
	root.AddCommand(cmds...)
}

func Execute(run func(cmd *cobra.Command, args []string)) {
	root.Run = run

	if err := root.Execute(); err != nil {
		panic(fmt.Errorf("cmd execute error:%w", err))
	}
}
