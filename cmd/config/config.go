package lpsumsegconfig

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

// cmd/root.go
type ConfigKey string

type Config struct {
	LP_AI_API_Key             string
	HF_TEXT_SUMMARY_API_Key   string
	AWS_ACCESS_KEY_ID_Key     string
	AWS_SECRET_ACCESS_KEY_Key string
	AWS_REGION                string
}
type AWSConfig struct {
	InputBucketName           string `mapstructure:"input_bucket_name"`
	OutputBucketName          string `mapstructure:"output_bucket_name"`
	S3InputVideoPath          string `mapstructure:"s3_input_video_path"`
	S3OutputTranscriptionPath string `mapstructure:"s3_output_transcription_path"`
}

func LoadAWSConfig() (*AWSConfig, error) {
	viper.SetConfigName("aws_config") // Name of the config file (without extension)
	viper.SetConfigType("json")       // Specify the type of the config file
	viper.AddConfigPath(".")          // Look for the config in the root of project

	if err := viper.ReadInConfig(); err != nil { // Read the config file
		fmt.Printf("Error reading aws config file: %s\n", err)
		return nil, err
	}
	var awsconfig AWSConfig

	if err := viper.Unmarshal(&awsconfig); err != nil { // Unmarshal into the Config struct
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &awsconfig, nil
}
func LoadConfig() *Config {
	lpApiKey := os.Getenv("LP_AI_API_KEY")
	hfTextSummaryApiKey := os.Getenv("HF_TEXT_SUMMARY_API_KEY")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")

	if lpApiKey == "" || hfTextSummaryApiKey == "" || awsAccessKey == "" || awsSecretKey == "" || awsRegion == "" {
		log.Fatal("LP_AI_API_KEY or HF_VIDEO_TO_TEXT_API_KEY or AWS_ACCESS_KEY_ID or AWS_SECRET_ACCESS_KEY or AWS_REGION is not set in the environment")

	}

	return &Config{
		LP_AI_API_Key:             lpApiKey,
		HF_TEXT_SUMMARY_API_Key:   hfTextSummaryApiKey,
		AWS_ACCESS_KEY_ID_Key:     awsAccessKey,
		AWS_SECRET_ACCESS_KEY_Key: awsSecretKey,
		AWS_REGION:                awsRegion,
	}
}
