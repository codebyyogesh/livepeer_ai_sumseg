/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	livepeergo "github.com/livepeer/livepeer-go"
	"github.com/livepeer/livepeer-go/models/components"
	"github.com/spf13/cobra"
)

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stream called")
		if cfg == nil {
			log.Fatal("Configuration not loaded properly")
		}
		lpClient := livepeergo.New(
			livepeergo.WithSecurity(cfg.LP_AI_API_Key),
		)

		ctx := context.Background()
		res, err := lpClient.Stream.Create(ctx, components.NewStreamPayload{
			Name: "test_stream",
		})
		if err != nil {
			log.Fatal(err)
		}
		if res.Stream != nil {
			log.Printf("Stream created successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(streamCmd)
}
