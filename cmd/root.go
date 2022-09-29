package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"troisdizaines.com/rcalc/rcalc"
)

var rootCmd = &cobra.Command{
	Use:   "rcalc",
	Short: "Rcalc is a RPN command line calculator",
	Long: `Later a more complete long doc
           With a website link`,
	Run: func(cmd *cobra.Command, args []string) {
		rcalc.Run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
