package caption

import (
	"fmt"

	lpsumsegconfig "github.com/codebyyogesh/livepeer_ai_sumseg.git/cmd/config"
	"github.com/spf13/cobra"
)

func handleSummaryCommand(videoURL string, env *lpsumsegconfig.Config) error {
	tcParams := newTranscribeParams()

	err := processTranscription(tcParams, videoURL, env)
	if err != nil {
		return err
	}

	fmt.Println("Summary Content:")
	fmt.Println(transcriptionResult.Summary)

	return nil
}

var SummaryCmd = &cobra.Command{
	Use:   "summary [videoURL]",
	Short: "Generate summary of video",
	RunE: func(cmd *cobra.Command, args []string) error {
		env, ok := cmd.Context().Value(lpsumsegconfig.ConfigKey("config")).(*lpsumsegconfig.Config) // Retrieve config from context

		if !ok {
			return fmt.Errorf("asset:failed to retrieve config from context")
		}

		return handleSummaryCommand(args[0], env)
	},
}
