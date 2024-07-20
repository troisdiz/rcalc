package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"troisdizaines.com/rcalc/rcalc"
)

func NewRunCommand(rootOptions *RootOptions) *cobra.Command {

	var runCmd = &cobra.Command{
		Use:   "run <program> --args [args]",
		Short: "loads and run a Rcalc program",
		Long:  `Run long desc`,
		Run: func(cmd *cobra.Command, args []string) {
			var rCalcDir string
			if rootOptions.configFolder == "" {
				dir, err := homedir.Dir()
				if err != nil {
					fmt.Println("Cannot get home directory")
					os.Exit(-1)
				}
				rCalcDir = path.Join(dir, ".rcalc")
			} else {
				rCalcDir = rootOptions.configFolder
			}

			rcalc.Run(rCalcDir, true, rootOptions.debugMode)
		},
	}
	return runCmd
}
