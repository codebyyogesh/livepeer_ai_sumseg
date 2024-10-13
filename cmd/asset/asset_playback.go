/*
Copyright © 2024 Yogesh Kulkarni <yogeshcodes@zohomail.in>
*/
package asset

import (
	"context"
	"errors"
	"fmt"
	"log"

	lpsumsegconfig "github.com/codebyyogesh/livepeer_ai_sumseg.git/cmd/config"
	livepeergo "github.com/livepeer/livepeer-go"
	"github.com/livepeer/livepeer-go/models/sdkerrors"
	"github.com/spf13/cobra"
)

// Set the playback ID of your uploaded asset
var playbackID = `de93swf0r2g7tlrz` // Replace with your actual playback ID
// assetCmd represents the asset command
var AssetPlaybackCmd = &cobra.Command{
	Use:   "assetplayback",
	Short: "A asset playback",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, ok := cmd.Context().Value(lpsumsegconfig.ConfigKey("config")).(*lpsumsegconfig.Config) // Retrieve config from context

		if !ok {
			return fmt.Errorf("asset:failed to retrieve config from context")
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
				return fmt.Errorf("error retrieving playback info: %s", sdkErr.Error())
			}
			log.Fatalf("Unexpected error: %s", err)
		}

		// Print the playback URL
		fmt.Printf("Playback URL: %+v\n", playbackInfo.PlaybackInfo)
		return nil
	},
}
