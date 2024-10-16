/*
Copyright © 2024 Yogesh Kulkarni <yogeshcodes@zohomail.in>
*/
package caption

import (
	"fmt"
	"os"
	"path/filepath"

	lpsumsegconfig "github.com/codebyyogesh/livepeer_ai_sumseg.git/cmd/config"
	"github.com/spf13/cobra"
)

// Function to write captions to a file
func writeCaptionsToFile(captions []string) error {
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

	// Define the full path for the caption file
	captionFilePath := filepath.Join(binDir, "caption.txt")

	// Create or open the caption file for writing
	file, err := os.Create(captionFilePath)
	if err != nil {
		return fmt.Errorf("failed to create caption file: %v", err)
	}
	defer file.Close() // Ensure the file is closed after writing

	// Write each caption to the file
	for _, caption := range captions {
		if _, err := file.WriteString(caption + "\n"); err != nil {
			return fmt.Errorf("failed to write caption to file: %v", err)
		}
	}

	fmt.Printf("Captions written to %s\n", captionFilePath)
	return nil
}

// Function to write captions to a file

func handleCaptionCommand(videoURL string, env *lpsumsegconfig.Config) error {

	tcParams, err := newTranscribeParams(env)
	if err != nil {
		return err
	}

	err = processTranscription(tcParams, videoURL)
	if err != nil {
		return err
	}
	// Extract the single transcript into a slice of strings
	var captions []string
	if len(transcriptionResult.Results.Transcripts) > 0 {
		captions = append(captions, transcriptionResult.Results.Transcripts[0].Transcript)
	}

	// Call WriteCaptionsToFile with the extracted caption
	if err := writeCaptionsToFile(captions); err != nil {
		fmt.Println("Error writing captions:", err)
	}
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
