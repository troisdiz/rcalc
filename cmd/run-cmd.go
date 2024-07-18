package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"troisdizaines.com/rcalc/rcalc"
)

func NewCommandRun() *cobra.Command {

	var runCmd = &cobra.Command{
		Use:   "rcalc run <program> --args [args]",
		Short: "loads and run a Rcalc program",
		Long:  `Run long desc`,
		Run: func(cmd *cobra.Command, args []string) {
			var rCalcDir string
			if *configFolder == "" {
				dir, err := homedir.Dir()
				if err != nil {
					fmt.Println("Cannot get home directory")
					os.Exit(-1)
				}
				rCalcDir = path.Join(dir, ".rcalc")
			} else {
				rCalcDir = *configFolder
			}

			rcalc.Run(rCalcDir, true, *debugMode)
		},
	}
	return runCmd
}

func init() {
	debugMode = rootCmd.PersistentFlags().BoolP("debugMode", "d", false, "Sets logs verbosity to debug")
	configFolder = rootCmd.PersistentFlags().StringP("configFolder", "c", "", "Sets the config folder")
}
