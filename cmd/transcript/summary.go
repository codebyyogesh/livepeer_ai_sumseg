package caption

import (
	"fmt"
	"time"

	"github.com/Kardbord/hfapigo/v3"
	lpsumsegconfig "github.com/codebyyogesh/livepeer_ai_sumseg.git/cmd/config"
	"github.com/spf13/cobra"
)

func summarizeTranscript(inputs []string, env *lpsumsegconfig.Config) ([]string, error) {
	hfapigo.SetAPIKey(env.HF_TEXT_SUMMARY_API_Key)
	fmt.Printf("\nSending request")

	type ChanRv struct {
		resps []*hfapigo.SummarizationResponse
		err   error
	}
	ch := make(chan ChanRv)

	go func() {
		sresps, err := hfapigo.SendSummarizationRequest(hfapigo.RecommmendedSummarizationModel, &hfapigo.SummarizationRequest{
			Inputs:  inputs,
			Options: *hfapigo.NewOptions().SetWaitForModel(true),
		})

		ch <- ChanRv{sresps, err}
	}()

	for {
		select {
		case chrv := <-ch:
			fmt.Println()
			if chrv.err != nil {
				fmt.Println(chrv.err)
				return nil, chrv.err
			}
			summaries := []string{}
			for _, resp := range chrv.resps {
				summaries = append(summaries, resp.SummaryText)
			}
			return summaries, nil
		default:
			fmt.Print(".")
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func handleSummaryCommand(videoURL string, env *lpsumsegconfig.Config) error {
	tcParams, err := newTranscribeParams(env)
	if err != nil {
		return err
	}

	err = processTranscription(tcParams, videoURL)
	if err != nil {
		return err
	}
	// At this point we have the whole transcription result, we now use AI to summarize it

	// Step 1: Extract transcript strings into a slice
	var transcripts []string
	for _, t := range transcriptionResult.Results.Transcripts {
		transcripts = append(transcripts, t.Transcript)
	}
	summary, err := summarizeTranscript(transcripts, env)
	if err != nil {
		return err
	}
	transcriptionResult.Summary = summary[0]
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
