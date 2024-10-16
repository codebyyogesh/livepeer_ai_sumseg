package caption

import (
	"fmt"
	"os"
	"path/filepath"

	lpsumsegconfig "github.com/codebyyogesh/livepeer_ai_sumseg.git/cmd/config"
	"github.com/spf13/cobra"
)

func writeSubtitlesToFile(subtitles string) error {
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

	// Define the full path for the subtitles file
	subtitlesFilePath := filepath.Join(binDir, "subtitles.srt")

	// Create or open the subtitles file for writing
	file, err := os.Create(subtitlesFilePath)
	if err != nil {
		return fmt.Errorf("failed to create substitles file: %v", err)
	}
	defer file.Close() // Ensure the file is closed after writing

	// Write subtitles to the file
	if _, err := file.WriteString(subtitles); err != nil {
		return fmt.Errorf("failed to write subtitles to file: %v", err)
	}

	fmt.Printf("Subtitles written to %s\n", subtitlesFilePath)
	return nil
}

func handleSubtitlesCommand(videoURL string, env *lpsumsegconfig.Config) error {
	tcParams, err := newTranscribeParams(env)
	if err != nil {
		return err
	}

	err = processTranscription(tcParams, videoURL)
	if err != nil {
		return err
	}
	writeSubtitlesToFile(transcriptionResult.Subtitles)
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
