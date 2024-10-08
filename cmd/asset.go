/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"

	livepeergo "github.com/livepeer/livepeer-go"
	"github.com/livepeer/livepeer-go/models/sdkerrors"
	"github.com/spf13/cobra"
)

// Set the playback ID of your uploaded asset
var playbackID = `de93swf0r2g7tlrz` // Replace with your actual playback ID
// assetCmd represents the asset command
var assetCmd = &cobra.Command{
	Use:   "asset",
	Short: "A asset playback",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("asset called")
		if cfg == nil {
			log.Fatal("Configuration not loaded properly")
		}
		lpClient := livepeergo.New(
			livepeergo.WithSecurity(cfg.LP_AI_API_Key),
		)

		ctx := context.Background()
		// Retrieve playback info
		playbackInfo, err := lpClient.Playback.Get(ctx, playbackID)
		if err != nil {
			var sdkErr *sdkerrors.Error
			if errors.As(err, &sdkErr) {
				log.Fatalf("Error retrieving playback info: %s", sdkErr.Error())
			}
			log.Fatalf("Unexpected error: %s", err)
		}

		// Print the playback URL
		fmt.Printf("Playback URL: %+v\n", playbackInfo.PlaybackInfo)
	},
}

func init() {
	rootCmd.AddCommand(assetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// assetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// assetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
