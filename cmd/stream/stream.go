/*
Copyright © 2024 Yogesh Kulkarni <yogeshcodes@zohomail.in>
*/
package stream

import (
	"context"
	"errors"
	"fmt"
	"log"

	lpsumsegconfig "github.com/codebyyogesh/livepeer_ai_sumseg.git/cmd/config"
	livepeergo "github.com/livepeer/livepeer-go"
	"github.com/livepeer/livepeer-go/models/components"
	"github.com/livepeer/livepeer-go/models/sdkerrors"
	"github.com/spf13/cobra"
)

// streamCmd represents the stream command
var StreamCmd = &cobra.Command{
	Use:   "stream",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, ok := cmd.Context().Value(lpsumsegconfig.ConfigKey("config")).(*lpsumsegconfig.Config) // Retrieve config from context

		if !ok {
			return fmt.Errorf("asset:failed to retrieve config from context")
		}

		lpClient := livepeergo.New(
			livepeergo.WithSecurity(cfg.LP_AI_API_Key),
		)

		ctx := context.Background()
		res, err := lpClient.Stream.Create(ctx, components.NewStreamPayload{
			Name: "test_stream",
		})
		if err != nil {
			var sdkErr *sdkerrors.Error
			if errors.As(err, &sdkErr) {
				return fmt.Errorf("failed to create stream: %s", sdkErr.Error())
			}
			log.Fatalf("stream:Unexpected error: %s", err)

		}
		if res.Stream != nil {
			log.Printf("Stream created successfully")
		}

		return nil
	},
}
