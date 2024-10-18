package image_to_video

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	lpsumsegconfig "github.com/codebyyogesh/livepeer_ai_sumseg.git/cmd/config"
	livepeeraigo "github.com/livepeer/livepeer-ai-go"
	"github.com/livepeer/livepeer-ai-go/models/components"
	"github.com/livepeer/livepeer-ai-go/models/operations"

	"github.com/spf13/cobra"
)

func ptr(s string) *string {
	return &s
}

var ImageToVideoCmd = &cobra.Command{
	Use:   "imagetovideo [filename]",
	Short: "Generate image to video",
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
		// Make an HTTP GET request
		client := &http.Client{
			Timeout: 180 * time.Second,
		}
		response, err := client.Get(args[0])
		if err != nil {
			panic(err) // Handle error appropriately in production code
		}
		defer response.Body.Close() // Ensure the response body is closed after reading

		// Check if the request was successful
		if response.StatusCode != http.StatusOK {
			panic(fmt.Errorf("failed to fetch image: %s", response.Status))
		}

		// Read the content of the response
		content, err := io.ReadAll(response.Body)
		if err != nil {
			panic(err) // Handle error appropriately in production code
		}
		request := components.BodyGenImageToVideo{
			Image: components.BodyGenImageToVideoImage{
				FileName: args[0],
				Content:  content,
			},
			ModelID: ptr("stabilityai/stable-video-diffusion-img2vid-xt-1-1"),
		}
		res, err := s.Generate.ImageToVideo(ctx, request,
			operations.WithOperationTimeout(120*time.Second))

		if err != nil {
			log.Fatal(err)
		}
		if res.VideoResponse != nil {
			log.Println("Response:", res.VideoResponse.Images[0].URL)
		}
		return nil
	},
}
