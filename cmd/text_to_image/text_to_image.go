/*
Copyright © 2024 Yogesh Kulkarni <yogeshcodes@zohomail.in>
*/
package text_to_image

import (
	"context"
	"fmt"
	"log"

	lpsumsegconfig "github.com/codebyyogesh/livepeer_ai_sumseg.git/cmd/config"
	livepeeraigo "github.com/livepeer/livepeer-ai-go"
	"github.com/livepeer/livepeer-ai-go/models/components"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var TextToImageCmd = &cobra.Command{
	Use:   "texttoimage [prompt]",
	Short: "Generate text to image",
	RunE: func(cmd *cobra.Command, args []string) error {
		env, ok := cmd.Context().Value(lpsumsegconfig.ConfigKey("config")).(*lpsumsegconfig.Config) // Retrieve config from context

		if !ok {
			return fmt.Errorf("asset:failed to retrieve config from context")
		}

		// Use the API key loaded in root.go
		s := livepeeraigo.New(
			// index 0 is for https://dream-gateway.livepeer.cloud
			livepeeraigo.WithServerIndex(0), // for https://livepeer.studio/api/beta/generate

			livepeeraigo.WithSecurity(env.LP_AI_API_Key),
		)
		ctx := context.Background()
		res, err := s.Generate.TextToImage(ctx, components.TextToImageParams{
			// SG161222/RealVisXL_V4.0_Lightning:
			// ByteDance/SDXL-Lightning
			ModelID: livepeeraigo.String("ByteDance/SDXL-Lightning"),
			Prompt:  args[0],
		},
		)
		if err != nil {
			log.Printf("serve error: %+v", err)
		}
		if res.ImageResponse != nil {
			// handle response
			log.Println("Response:", res.ImageResponse.Images[0].URL)
		}

		return nil
	},
}
