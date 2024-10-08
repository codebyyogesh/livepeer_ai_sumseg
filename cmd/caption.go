/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/transcribe"
	"github.com/aws/aws-sdk-go-v2/service/transcribe/types"
	"github.com/spf13/cobra"
)

var videoURL = "https://vod-cdn.lp-playback.studio/raw/jxf4iblf6wlsyor6526t4tcmtmqa/catalyst-vod-com/hls/de93swf0r2g7tlrz/360p0.mp4"

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

func readTranscriptionResult(s3Client *s3.Client, bucketName, fileKey string) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	}

	result, err := s3Client.GetObject(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to get object from S3: %v", err)
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read object body: %v", err)
	}

	return string(body), nil
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

func createTranscriptionJob(transcribeClient *transcribe.Client, bucketName, videoFileName, outputBucketName, jobName string) (string, error) {

	// Check if transcription job already exists, if yes delete it

	getResult, err := transcribeClient.GetTranscriptionJob(context.TODO(), &transcribe.GetTranscriptionJobInput{TranscriptionJobName: aws.String(jobName)})

	// As there is no error, such a job already exists , so we must first delete it and then create a new one
	if err == nil {
		out, err := transcribeClient.DeleteTranscriptionJob(context.TODO(), &transcribe.DeleteTranscriptionJobInput{TranscriptionJobName: getResult.TranscriptionJob.TranscriptionJobName})
		if err != nil {
			log.Fatalf("createTranscriptionJob: Error deleting transcription job: %v", err)
		}
		fmt.Printf("createTranscriptionJob: Deleted existing Transcription job: %s\n", out)
	}
	// if err is not nil , means such a job does not exist and we can continue creating it directly
	// Prepare the input parameters for the transcription job
	input := &transcribe.StartTranscriptionJobInput{
		TranscriptionJobName: aws.String(jobName),
		LanguageCode:         "en-US", // Change as needed
		Media:                &types.Media{MediaFileUri: aws.String(fmt.Sprintf("s3://%s/%s", bucketName, videoFileName))},
		OutputBucketName:     aws.String(outputBucketName),
	}
	// Start a transcription job
	startResult, err := transcribeClient.StartTranscriptionJob(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("createTranscriptionJob: failed to start transcription job: %v", err)
	}
	fmt.Printf("createTranscriptionJob: Started transcription job: %s\n", *startResult.TranscriptionJob.TranscriptionJobName)
	return *startResult.TranscriptionJob.TranscriptionJobName, nil // Return job name or handle as needed
}

// captionCmd represents the caption command
var captionCmd = &cobra.Command{
	Use:   "caption",
	Short: "Generate caption for video",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("caption called")
		bucketName := "lpvideouploader"
		cfg, _ := loadAWSConfig()
		s3Client := s3.NewFromConfig(cfg)

		videoFilePath := filepath.Join("cmd", "lp1.mp4")

		// Upload the MP4 file to S3
		fileURL, err := uploadToS3(s3Client, bucketName, videoFilePath)
		if err != nil {
			log.Fatalf("Failed to upload file: %v", err)
		}

		fmt.Printf("File uploaded successfully: %s\n", fileURL)
		// Transcription Job for captions and subtitles
		transcribeClient := transcribe.NewFromConfig(cfg)
		videoFileName := "videos/lp1.mp4"     // Path in S3
		outputBucketName := "lpvideouploader" // Change this to your output bucket
		jobName := "GetCaptionsAndSubtitlesTranscriptionJob"

		jobID, err := createTranscriptionJob(transcribeClient, bucketName, videoFileName, outputBucketName, jobName)
		if err != nil {
			log.Fatalf("Error creating transcription job: %v", err)
		}

		fmt.Printf("Transcription job created successfully with ID: %s\n", jobID)

		err = checkTranscriptionJobStatus(transcribeClient, jobID) // Check the status of the transcription job
		if err != nil {
			log.Fatalf("Error checking transcription job status: %v", err)
		}

		fileKey := fmt.Sprintf("%s.json", jobID) // Assuming the output file is named after the job ID with a .json extension

		transcriptContent, err := readTranscriptionResult(s3Client, outputBucketName, fileKey)
		if err != nil {
			log.Fatalf("Error reading transcription result: %v", err)
		}

		fmt.Println("Transcript Content:")
		fmt.Println(transcriptContent)
	},
}

func init() {
	rootCmd.AddCommand(captionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// captionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// captionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
