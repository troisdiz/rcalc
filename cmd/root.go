package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"os"
	"path"
	"troisdizaines.com/rcalc/rcalc"
)

var rootCmd = &cobra.Command{
	Use:   "rcalc",
	Short: "Rcalc is a RPN command line calculator",
	Long: `Later a more complete long doc
           With a website link`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := homedir.Dir()
		if err != nil {
			fmt.Println("Cannot get home directory")
			os.Exit(-1)
		}

		rCalcDir := path.Join(dir, ".rcalc")

		rcalc.Run(rCalcDir)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
