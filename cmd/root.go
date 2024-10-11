/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var cfg *Config // Global configuration

// Config structure to hold environment variables
type Config struct {
	LP_AI_API_Key             string
	HF_VIDEO_TO_TEXT_API_Key  string
	AWS_ACCESS_KEY_ID_Key     string
	AWS_SECRET_ACCESS_KEY_Key string
	AWS_REGION                string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lp_ai_sumseg",
	Short: "AI CLI tool to get captions, subtitles and summarize video",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Load the API key from the environment
func LoadConfig() *Config {
	lpApiKey := os.Getenv("LP_AI_API_KEY")
	hfVidToTextApiKey := os.Getenv("HF_VIDEO_TO_TEXT_API_KEY")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")

	if lpApiKey == "" || hfVidToTextApiKey == "" || awsAccessKey == "" || awsSecretKey == "" || awsRegion == "" {
		log.Fatal("LP_AI_API_KEY or HF_VIDEO_TO_TEXT_API_KEY or AWS_ACCESS_KEY_ID or AWS_SECRET_ACCESS_KEY or AWS_REGION is not set in the environment")

	}

	return &Config{
		LP_AI_API_Key:             lpApiKey,
		HF_VIDEO_TO_TEXT_API_Key:  hfVidToTextApiKey,
		AWS_ACCESS_KEY_ID_Key:     awsAccessKey,
		AWS_SECRET_ACCESS_KEY_Key: awsSecretKey,
		AWS_REGION:                awsRegion,
	}
}

func init() {
	// Load configuration in the root command
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	cfg = LoadConfig()
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.livepeer_ai_sumseg.git.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
