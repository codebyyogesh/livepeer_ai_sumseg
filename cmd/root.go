/*
Copyright © 2024 Yogesh Kulkarni <yogeshcodes@zohomail.in>
*/
package cmd

import (
	"context"
	"log"
	"os"

	"github.com/codebyyogesh/livepeer_ai_sumseg.git/cmd/asset"
	lpsumsegconfig "github.com/codebyyogesh/livepeer_ai_sumseg.git/cmd/config"
	"github.com/codebyyogesh/livepeer_ai_sumseg.git/cmd/stream"
	caption "github.com/codebyyogesh/livepeer_ai_sumseg.git/cmd/transcript"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// Config structure to hold environment variables

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "lp_ai_sumseg",
	Short: "AI CLI tool to get captions, subtitles and summarize video",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load configuration before any command is executed
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
		cfg := lpsumsegconfig.LoadConfig() // Assuming LoadConfig is a function that loads your config
		if cfg == nil {
			log.Fatal("Failed to load configuration")
		}
		cmd.SetContext(context.WithValue(cmd.Context(), lpsumsegconfig.ConfigKey("config"), cfg)) // Store config in context
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Load the API key from the environment

func init() {
	// Load configuration in the root command

	// Add subcommands
	RootCmd.AddCommand(asset.AssetPlaybackCmd) // Register asset command
	RootCmd.AddCommand(asset.AssetUploadCmd)   // Register asset command
	RootCmd.AddCommand(stream.StreamCmd)       // Register stream command
	RootCmd.AddCommand(caption.CaptionCmd)     // Register caption command
	RootCmd.AddCommand(caption.SubtitlesCmd)   // Register subtitles command
	RootCmd.AddCommand(caption.SummaryCmd)     // Register summary command
	// Add other commands as needed
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
