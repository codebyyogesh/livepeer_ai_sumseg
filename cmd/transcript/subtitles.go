package caption

import (
	"fmt"

	lpsumsegconfig "github.com/codebyyogesh/livepeer_ai_sumseg.git/cmd/config"
	"github.com/spf13/cobra"
)

func handleSubtitlesCommand(videoURL string, env *lpsumsegconfig.Config) error {
	tcParams := newTranscribeParams()

	err := processTranscription(tcParams, videoURL, env)
	if err != nil {
		return err
	}

	fmt.Println("Subtitle Content:")
	fmt.Println(transcriptionResult.Subtitles)

	return nil
}

var SubtitlesCmd = &cobra.Command{
	Use:   "subtitles [videoURL]",
	Short: "Generate subtitles for video",
	RunE: func(cmd *cobra.Command, args []string) error {
		env, ok := cmd.Context().Value(lpsumsegconfig.ConfigKey("config")).(*lpsumsegconfig.Config) // Retrieve config from context

		if !ok {
			return fmt.Errorf("asset:failed to retrieve config from context")
		}

		return handleSubtitlesCommand(args[0], env)
	},
}
