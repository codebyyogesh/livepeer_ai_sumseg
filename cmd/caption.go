/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/transcribe"
	"github.com/aws/aws-sdk-go-v2/service/transcribe/types"
	"github.com/spf13/cobra"
)

// Path in S3
var s3InputVideoPath = "videos/process.mp4"
var s3OutputTranscriptionPath = "transcriptions/"

// Global variable to hold the transcription result
var transcriptionResult struct {
	Results struct {
		Transcripts []struct {
			Transcript string `json:"transcript"`
		} `json:"transcripts"`
	} `json:"results"`
	Subtitles              string // To hold subtitle content
	Summary                string // To hold summary content
	transcriptionProcessed bool
	lastProcessedVideoFile string
}

func loadAWSConfig() (aws.Config, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
				cfg.AWS_ACCESS_KEY_ID_Key,
				cfg.AWS_SECRET_ACCESS_KEY_Key,
				"")),
		))
	if err != nil {
		return aws.Config{}, fmt.Errorf("unable to load SDK config: %v", err)
	}
	return cfg, nil

}

/*
	func uploadToS3(s3Client *s3.Client, bucketName, filePath string) (string, error) {
		// Load AWS credentials and create an S3 client as before...
		// Create an S3 service client

		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to open file %q: %v", filePath, err)
		}
		defer file.Close()

		// Get the file name
		fileName := filepath.Base(filePath)

		// Construct the S3 key (path within the bucket)
		s3Key := fmt.Sprintf("videos/%s", fileName) // This specifies that the file is going into the "videos" folder

		// Upload the file to S3
		_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(s3Key),
			Body:   file,
		})
		if err != nil {
			return "", fmt.Errorf("failed to upload file to S3: %v", err)
		}

		// Construct the S3 file URL
		fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, s3Key)

		return fileURL, nil
	}
*/

// streamToS3 streams the content from the provided URL directly to the given S3 bucket
func streamToS3(s3Client *s3.Client, ipfsURL, bucketName string, contentLength int64) error {
	// Send a GET request to the URL to stream the content
	resp, err := http.Get(ipfsURL)
	if err != nil {
		return fmt.Errorf("failed to stream from IPFS: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is 200 OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to stream file, status: %s", resp.Status)
	}

	// Create the upload input for streaming the file
	// All files will be stored as process.mp4 under videos folder so that we do not keep uploading all files into s3
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s3InputVideoPath),

		Body:          resp.Body,                // Streaming the content directly to S3
		ContentLength: aws.Int64(contentLength), // Optional: Set the content length for the contentLength,
	}
	fmt.Printf("Uploading %s to %s/%s\n", ipfsURL, bucketName, s3InputVideoPath)
	// Upload the file to S3
	_, err = s3Client.PutObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %v", err)
	}

	return nil
}

func readTranscriptionResult(s3Client *s3.Client, bucketName, jsonFileKey, srtFilekey string) error {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s3OutputTranscriptionPath + jsonFileKey),
	}

	result, err := s3Client.GetObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to get json object from S3: %v", err)
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return fmt.Errorf("failed to read object body: %v", err)
	}

	err = json.Unmarshal(body, &transcriptionResult)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	input = &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s3OutputTranscriptionPath + srtFilekey),
	}

	result, err = s3Client.GetObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to get srt object from S3: %v", err)
	}
	defer result.Body.Close()

	srtBody, err := io.ReadAll(result.Body)
	if err != nil {
		return fmt.Errorf("failed to read object body: %v", err)
	}
	transcriptionResult.Subtitles = string(srtBody) // Store subtitles in global variable

	return nil
}

func checkTranscriptionJobStatus(transcribeClient *transcribe.Client, jobName string) error {
	for {
		input := &transcribe.GetTranscriptionJobInput{
			TranscriptionJobName: aws.String(jobName),
		}

		result, err := transcribeClient.GetTranscriptionJob(context.TODO(), input)
		if err != nil {
			return fmt.Errorf("failed to get transcription job status: %v", err)
		}
		status := result.TranscriptionJob.TranscriptionJobStatus
		//status := *result.TranscriptionJob.TranscriptionJobStatus
		fmt.Printf("Current status of job %s: %s\n", jobName, status)

		if status == "COMPLETED" {
			fmt.Printf("Transcription completed! Results available at: %s\n", *result.TranscriptionJob.Transcript.TranscriptFileUri)
			break
		} else if status == "FAILED" {
			fmt.Printf("Transcription job failed. Reason: %s\n", *result.TranscriptionJob.FailureReason)
			break
		}

		time.Sleep(5 * time.Second) // Wait before checking again
	}
	return nil
}

// isTranscriptionJobNotFound checks if the error is a "job not found" error
func deleteTranscriptionJob(transcribeClient *transcribe.Client, jobName string) error {
	_, err := transcribeClient.DeleteTranscriptionJob(context.TODO(), &transcribe.DeleteTranscriptionJobInput{
		TranscriptionJobName: aws.String(jobName),
	})
	if err != nil {
		return fmt.Errorf("deleteTranscriptionJob: failed to delete job %s: %v", jobName, err)
	}
	return nil
}

func createTranscriptionJob(transcribeClient *transcribe.Client, inputBucketName, videoFileName, outputBucketName, jobName string) (string, error) {

	// Check if transcription job already exists, if yes delete it

	getResult, err := transcribeClient.GetTranscriptionJob(context.TODO(), &transcribe.GetTranscriptionJobInput{TranscriptionJobName: aws.String(jobName)})

	// As there is no error, such a job already exists , so we must first delete it and then create a new one
	if err == nil && getResult.TranscriptionJob != nil {
		// Delete the existing job
		err = deleteTranscriptionJob(transcribeClient, *getResult.TranscriptionJob.TranscriptionJobName)
		if err != nil {
			log.Fatalf("createTranscriptionJob: Error deleting transcription job: %v", err)
			return "", fmt.Errorf("createTranscriptionJob: Error deleting transcription job: %v", err)
		}
		fmt.Printf("createTranscriptionJob: Deleted existing Transcription job: %s\n", jobName)
	} else if err != nil {
		// If the error is anything other than job not found, handle it
		/*return "", fmt.Errorf("createTranscriptionJob: failed to get transcription job: %v", err)*/
		log.Printf("createTranscriptionJob: failed to get transcription job: %v", err)
	}

	// if err is not nil , means such a job does not exist and we can continue creating it directly
	// Prepare the input parameters for the transcription job
	input := &transcribe.StartTranscriptionJobInput{
		TranscriptionJobName: aws.String(jobName),
		LanguageCode:         "en-US", // Change as needed
		Media:                &types.Media{MediaFileUri: aws.String(fmt.Sprintf("s3://%s/%s", inputBucketName, videoFileName))},
		Subtitles: &types.Subtitles{
			Formats: []types.SubtitleFormat{
				types.SubtitleFormatSrt,
			},
			OutputStartIndex: aws.Int32(1),
		},
		OutputBucketName: aws.String(outputBucketName),
		OutputKey:        aws.String("transcriptions/"),
	}
	// Start a transcription job
	startResult, err := transcribeClient.StartTranscriptionJob(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("createTranscriptionJob: failed to start transcription job: %v", err)
	}
	fmt.Printf("createTranscriptionJob: Started transcription job: %s\n", *startResult.TranscriptionJob.TranscriptionJobName)
	return *startResult.TranscriptionJob.TranscriptionJobName, nil // Return job name or handle as needed
}

// getInputFileSize fetches the content length from IPFS by sending an HTTP HEAD request
func getInputFileSize(fileURL string) (int64, error) {
	// Send a HEAD request to get the content length
	resp, err := http.Head(fileURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch file size from URL: %v", err)
	}
	defer resp.Body.Close()

	// Check if the Content-Length header is present
	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		return 0, fmt.Errorf("content length not available in the response")
	}

	// Convert the content length to an integer
	size, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse content length: %v", err)
	}

	return size, nil
}

func processTranscription(inputBucketName, videoFileURL, outputBucketName, jobName string) error {
	if transcriptionResult.transcriptionProcessed {
		return nil // Return already processed result
	}
	cfg, err := loadAWSConfig()
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %v", err)
	}
	s3Client := s3.NewFromConfig(cfg)
	transcribeClient := transcribe.NewFromConfig(cfg)

	// Step 1: Fetch the content length (file size) from URL
	contentLength, err := getInputFileSize(videoFileURL)
	if err != nil {
		log.Fatalf("Failed to get file size from URL: %v", err)
	}
	fmt.Printf("Content-Length: %d bytes\n", contentLength)

	// Upload the MP4 file to S3
	//fileURL, err := uploadToS3(s3Client, inputBucketName, videoFileURL)
	// Step 2: Stream the MP4 file from the URL to S3
	err = streamToS3(s3Client, videoFileURL, inputBucketName, contentLength)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	fmt.Printf("File uploaded successfully to S3: \n")

	//Create Transcription Job

	jobID, err := createTranscriptionJob(transcribeClient, inputBucketName, s3InputVideoPath, outputBucketName, jobName)
	if err != nil {
		return fmt.Errorf("error creating transcription job: %v", err)
	}

	fmt.Printf("Transcription job created successfully with ID: %s\n", jobID)

	err = checkTranscriptionJobStatus(transcribeClient, jobID) // Check the status of the transcription job
	if err != nil {
		return fmt.Errorf("error checking transcription job status: %v", err)
	}
	jsonFileKey := fmt.Sprintf("%s.json", jobID) // Assuming the output file is named after the job ID with a .json extension

	srtFileKey := fmt.Sprintf("%s.srt", jobID) // Assuming the output file is named after the job ID with a .json extension

	err = readTranscriptionResult(s3Client, outputBucketName, jsonFileKey, srtFileKey)
	if err != nil {
		return fmt.Errorf("error reading transcription result: %v", err)
	}
	transcriptionResult.transcriptionProcessed = true // Mark as processed
	transcriptionResult.lastProcessedVideoFile = videoFileURL
	return nil
}

func handleCaptionCommand(videoURL string) error {

	inputBucketName := "lpvideouploader"
	outputBucketName := "lpvideouploader"
	jobName := "GetCaptionsAndSubtitlesTranscriptionJob"

	err := processTranscription(inputBucketName, videoURL, outputBucketName, jobName)
	if err != nil {
		return err
	}
	fmt.Println("Transcript Content:")
	fmt.Println(transcriptionResult.Results.Transcripts)

	return nil
}

func handleSubtitlesCommand(videoURL string) error {
	inputBucketName := "lpvideouploader"
	outputBucketName := "lpvideouploader"
	jobName := "GetCaptionsAndSubtitlesTranscriptionJob"

	err := processTranscription(inputBucketName, videoURL, outputBucketName, jobName)
	if err != nil {
		return err
	}

	fmt.Println("Subtitle Content:")
	fmt.Println(transcriptionResult.Subtitles)

	return nil
}

func handleSummaryCommand(videoURL string) error {
	inputBucketName := "lpvideouploader"
	outputBucketName := "lpvideouploader"
	jobName := "GetCaptionsAndSubtitlesTranscriptionJob"

	err := processTranscription(inputBucketName, videoURL, outputBucketName, jobName)
	if err != nil {
		return err
	}

	fmt.Println("Summary Content:")
	fmt.Println(transcriptionResult.Summary)

	return nil
}

var captionCmd = &cobra.Command{
	Use:   "caption [videoURL]",
	Short: "Generate caption for video",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleCaptionCommand(args[0]) // Pass the video URL to the handler
	},
}

var subtitlesCmd = &cobra.Command{
	Use:   "subtitles [videoURL]",
	Short: "Generate subtitles for video",
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleSubtitlesCommand(args[0])
	},
}

var summaryCmd = &cobra.Command{
	Use:   "summary [videoURL]",
	Short: "Generate summary of video",
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleSummaryCommand(args[0])
	},
}

func init() {
	//rootCmd.AddCommand(captionCmd)

	rootCmd.AddCommand(captionCmd)
	rootCmd.AddCommand(subtitlesCmd)
	rootCmd.AddCommand(summaryCmd)
}
