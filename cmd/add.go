package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add values passed to function",
	Long:  `Demo application to demonstrate cobra featues`,
	Run: func(cmd *cobra.Command, args []string) {
		sum := 0
		for _, args := range args {
			num, err := strconv.Atoi(args)

			if err != nil {
				fmt.Println(err)
			}
			sum = sum + num
		}
		fmt.Println("result of addition is", sum)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
