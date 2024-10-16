package caption

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Kardbord/hfapigo/v3"
	lpsumsegconfig "github.com/codebyyogesh/livepeer_ai_sumseg.git/cmd/config"
	"github.com/spf13/cobra"
)

func writeSummaryToFile(summary string) error {
	// Get the path of the executable
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Get the directory of the executable
	exeDir := filepath.Dir(exePath)

	// Define the path to the bin directory
	binDir := filepath.Join(exeDir, "..", "bin") // Go up one level and then into bin

	// Check if bin directory exists; if not, use current working directory
	if _, err := os.Stat(binDir); os.IsNotExist(err) {
		// If bin directory doesn't exist, use current working directory
		cwd, _ := os.Getwd()
		binDir = filepath.Join(cwd, "bin")
	}

	// Create the bin directory if it doesn't exist
	if err := os.MkdirAll(binDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create bin directory: %v", err)
	}

	// Define the full path for the summary file
	summaryFilePath := filepath.Join(binDir, "summary.txt")

	// Create or open the summary file for writing
	file, err := os.Create(summaryFilePath)
	if err != nil {
		return fmt.Errorf("failed to create summary file: %v", err)
	}
	defer file.Close() // Ensure the file is closed after writing

	// Write summary to the file
	if _, err := file.WriteString(summary); err != nil {
		return fmt.Errorf("failed to write summary to file: %v", err)
	}

	fmt.Printf("Summary written to %s\n", summaryFilePath)
	return nil
}

func summarizeTranscript(inputs []string, env *lpsumsegconfig.Config) ([]string, error) {
	hfapigo.SetAPIKey(env.HF_TEXT_SUMMARY_API_Key)
	fmt.Printf("Summarizing transcript\n")
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

	writeSummaryToFile(transcriptionResult.Summary)

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
