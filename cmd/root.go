package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"troisdizaines.com/rcalc/rcalc"
)

var debugMode *bool
var configFolder *string

func NewRootCommand() *cobra.Command {

	var debugMode bool
	var configFolder string

	var rootCmd = &cobra.Command{
		Use:   "rcalc",
		Short: "Rcalc is a RPN command line calculator",
		Long: `Rcalc is a RPN command line calculator
It includes a programming language`,
		Run: func(cmd *cobra.Command, args []string) {
			var rCalcDir string
			if configFolder == "" {
				dir, err := homedir.Dir()
				if err != nil {
					fmt.Println("Cannot get home directory")
					os.Exit(-1)
				}
				rCalcDir = path.Join(dir, ".rcalc")
			} else {
				rCalcDir = configFolder
			}

			rcalc.Run(rCalcDir, true, *debugMode)
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debugMode", "d", false, "Sets logs verbosity to debug")
	rootCmd.PersistentFlags().StringVarP(&configFolder, "configFolder", "c", "", "Sets the config folder")
	return rootCmd
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
