package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"troisdizaines.com/rcalc/rcalc"
)

type RootOptions struct {
	debugMode    bool
	configFolder string
}

func NewRootCommand() *cobra.Command {

	rootOptions := &RootOptions{}

	var rootCmd = &cobra.Command{
		Use:   "rcalc",
		Short: "Rcalc is a RPN command line calculator",
		Long: `Rcalc is a RPN command line calculator
It includes a programming language`,
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

			rcalc.RunRepl(rCalcDir, true, rootOptions.debugMode)
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&rootOptions.debugMode, "debugMode", "d", false, "Sets logs verbosity to debug")
	rootCmd.PersistentFlags().StringVarP(&rootOptions.configFolder, "configFolder", "c", "", "Sets the config folder")

	rootCmd.AddCommand(NewRunCommand(rootOptions))

	return rootCmd
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
