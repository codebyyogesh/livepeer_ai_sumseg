/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"

	livepeeraigo "github.com/livepeer/livepeer-ai-go"
	"github.com/livepeer/livepeer-ai-go/models/components"
	"github.com/spf13/cobra"
)

func ptr(s string) *string {
	return &s
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A small server",
	Run: func(cmd *cobra.Command, args []string) {
		// Use the API key loaded in root.go
		if cfg == nil {
			log.Fatal("Configuration not loaded properly")
		}
		s := livepeeraigo.New(
			// index 0 is for https://dream-gateway.livepeer.cloud
			livepeeraigo.WithServerIndex(0), // for https://livepeer.studio/api/beta/generate

			livepeeraigo.WithSecurity(cfg.LP_AI_API_Key),
		)
		ctx := context.Background()
		res, err := s.Generate.TextToImage(ctx, components.TextToImageParams{
			// SG161222/RealVisXL_V4.0_Lightning:
			// ByteDance/SDXL-Lightning
			ModelID: ptr("ByteDance/SDXL-Lightning"),
			Prompt:  "A puppy dog sleeping on a sofa",
		},
		/*			operations.WithRetries(
					retry.Config{
						Strategy: "backoff",
						Backoff: &retry.BackoffStrategy{
							InitialInterval: 1,
							MaxInterval:     50,
							Exponent:        1.1,
							MaxElapsedTime:  100,
						},
						RetryConnectionErrors: false,
					}),
		*/
		)
		if err != nil {
			log.Printf("serve error: %+v", err)
		}
		if res.ImageResponse != nil {
			// handle response
			log.Println("response:", res.ImageResponse)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
