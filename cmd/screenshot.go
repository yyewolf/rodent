package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yyewolf/rodent/mischief"
)

// screenshotCmd represents the screenshot command
var screenshotCmd = &cobra.Command{
	Use:   "screenshot [URL]",
	Short: "Take a screenshot of an URL.",
	Long:  `Take a screenshot of an URL.`,
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		// Create a single rodent mischief (which can just be named a rodent then)
		rodent, err := mischief.New()
		if err != nil {
			panic(err)
		}

		// Take a screenshot of the URL
		screenshot, err := rodent.TakeScreenshot(url)
		if err != nil {
			panic(err)
		}

		// Save the screenshot to a file
		output, _ := cmd.Flags().GetString("output")
		err = os.WriteFile(output, screenshot, 0644)
		if err != nil {
			panic(err)
		}
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(screenshotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// screenshotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	screenshotCmd.Flags().StringP("output", "o", "screenshot.png", "Output file for the screenshot")
}
