package lpsumsegconfig

import (
	"log"
	"os"
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
