/*
Copyright © 2024 Yogesh Kulkarni <yogeshcodes@zohomail.in>
*/
package caption

import (
	"fmt"

	lpsumsegconfig "github.com/codebyyogesh/livepeer_ai_sumseg.git/cmd/config"
	"github.com/spf13/cobra"
)

func handleCaptionCommand(videoURL string, env *lpsumsegconfig.Config) error {

	tcParams, err := newTranscribeParams(env)
	if err != nil {
		return err
	}

	err = processTranscription(tcParams, videoURL)
	if err != nil {
		return err
	}
	fmt.Println("Transcript Content:")
	fmt.Println(transcriptionResult.Results.Transcripts)

	return nil
}

var CaptionCmd = &cobra.Command{
	Use:   "caption [videoURL]",
	Short: "Generate caption for video",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, ok := cmd.Context().Value(lpsumsegconfig.ConfigKey("config")).(*lpsumsegconfig.Config) // Retrieve config from context

		if !ok {
			return fmt.Errorf("asset:failed to retrieve config from context")
		}

		return handleCaptionCommand(args[0], env) // Pass the video URL to the handler
	},
}
