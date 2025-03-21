package cmd

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/yyewolf/rodent/api"
	"github.com/yyewolf/rodent/mischief"
)

var (
	host     string
	port     string
	browsers string

	concurrency          int
	browserRetakeTimeout int
	pageStabilityTimeout int
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the Mischief API Server.",
	Long:  `Start the Mischief API Server.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		mischief, err := mischief.New(
			mischief.WithConcurrency(concurrency),
			mischief.WithBrowserRetakeTimeout(time.Duration(browserRetakeTimeout)*time.Second),
			mischief.WithPageStabilityTimeout(time.Duration(pageStabilityTimeout)*time.Second),
			mischief.WithLogger(logger),
		)
		if err != nil {
			panic(err)
		}

		apiServer, err := api.New(
			api.WithHost(host),
			api.WithPort(port),
			api.WithMischief(mischief),
			api.WithLogger(logger),
		)
		if err != nil {
			panic(err)
		}

		apiServer.Start()

		// Wait for ctrl+c
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
		<-signalChannel

		err = apiServer.Shutdown(cmd.Context())
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
	apiCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to run the API server on.")
	apiCmd.Flags().StringVarP(&host, "host", "H", "localhost", "Host to run the API server on.")
	apiCmd.Flags().StringVarP(&browsers, "browsers", "b", "http://localhost:9222", "URLs to connect to external browsers. Use commas to separate multiple URLs.")
	apiCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of browsers to use to take screenshots concurrently.")
	apiCmd.Flags().IntVarP(&browserRetakeTimeout, "browser-retake-timeout", "r", 5, "Timeout used when taking a browser from the pool.")
	apiCmd.Flags().IntVarP(&pageStabilityTimeout, "page-stability-timeout", "s", 3, "Timeout used when waiting for the page to be stable.")
}
