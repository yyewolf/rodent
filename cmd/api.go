/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yyewolf/rodent/api"
	"github.com/yyewolf/rodent/mischief"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the Mischief API Server.",
	Long:  `Start the Mischief API Server.`,
	Run: func(cmd *cobra.Command, args []string) {
		concurrency, err := cmd.Flags().GetInt("concurrency")
		if err != nil {
			panic(err)
		}

		mischief, err := mischief.New(
			mischief.WithConcurrency(concurrency),
		)
		if err != nil {
			panic(err)
		}

		apiServer, err := api.New(
			api.WithHost(cmd.Flag("host").Value.String()),
			api.WithPort(cmd.Flag("port").Value.String()),
			api.WithMischief(mischief),
		)
		if err != nil {
			panic(err)
		}

		err = apiServer.Start()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	apiCmd.Flags().StringP("port", "p", "8080", "Port to run the API server on.")
	apiCmd.Flags().StringP("host", "H", "localhost", "Host to run the API server on.")
	apiCmd.Flags().StringP("browsers", "b", "http://localhost:9222", "URLs to connect to external browsers. Use commas to separate multiple URLs.")
	apiCmd.Flags().IntP("concurrency", "c", 1, "Number of browsers to use to take screenshots concurrently.")
	apiCmd.Flags().IntP("browser-retake-timeout", "r", 5, "Timeout used when taking a browser from the pool.")
	apiCmd.Flags().IntP("page-stability-timeout", "s", 3, "Timeout used when waiting for the page to be stable.")
}
