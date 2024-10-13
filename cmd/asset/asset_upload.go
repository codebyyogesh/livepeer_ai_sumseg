package asset

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	lpsumsegconfig "github.com/codebyyogesh/livepeer_ai_sumseg.git/cmd/config"
	livepeergo "github.com/livepeer/livepeer-go"
	"github.com/livepeer/livepeer-go/models/components"
	"github.com/spf13/cobra"
)

var AssetUploadCmd = &cobra.Command{
	Use:   "assetupload [filename] [videoURL]",
	Short: "Upload an asset (mp4 file) to livepeer with filename",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, ok := cmd.Context().Value(lpsumsegconfig.ConfigKey("config")).(*lpsumsegconfig.Config) // Retrieve config from context

		if !ok {
			return fmt.Errorf("asset:failed to retrieve config from context")
		}

		s := livepeergo.New(
			livepeergo.WithSecurity(cfg.LP_AI_API_Key),
		)

		ctx := context.Background()

		// Create a boolean to enable IPFS
		ipfsEnabled := true
		fmt.Printf("Uploading %s to livepeer\n", filepath.Base(args[0]))
		// Create an asset payload without webhook
		assetPayload := components.NewAssetPayload{
			Name:      filepath.Base(filepath.Base(args[0])),
			StaticMp4: livepeergo.Bool(true),
			// Initialize IPFS storage
			Storage: &components.NewAssetPayloadStorage{Ipfs: &components.NewAssetPayloadIpfs{
				Boolean: &ipfsEnabled, // Set this to true to enable IPFS
				// You can set NewAssetPayloadIpfs1 and Type if needed
			}},
			PlaybackPolicy: &components.PlaybackPolicy{
				Type: components.TypePublic, // Set to public instead of webhook
			},
			Profiles: []components.TranscodeProfile{
				{
					Width:   livepeergo.Int64(1280),
					Name:    livepeergo.String("720p"),
					Height:  livepeergo.Int64(720),
					Bitrate: 3000000,
					Quality: livepeergo.Int64(23),
					Fps:     livepeergo.Int64(30),
					FpsDen:  livepeergo.Int64(1),
					Gop:     livepeergo.String("2"),
					Profile: components.TranscodeProfileProfileH264Baseline.ToPointer(),
					Encoder: components.TranscodeProfileEncoderH264.ToPointer(),
				},
			},
		}
		res, err := s.Asset.Create(ctx, assetPayload)
		if err != nil {
			log.Fatal(err)
		}
		// Step 3: Once asset is created, Get Upload URL
		uploadURL := res.Data.URL // This should be provided in the response
		if uploadURL == "" {
			return fmt.Errorf("no upload URL returned")
		}

		fmt.Printf("Asset ID: %s\n", res.Data.Asset.ID)
		fmt.Printf("Asset Playback: %s\n", *res.Data.Asset.PlaybackID)

		// Step 4: Open the video file to be uploaded
		videoFile, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("failed to open video file: %v", err)
		}
		defer videoFile.Close()

		// Step 5: Upload the video file using PUT request to the asset created
		req, err := http.NewRequest("PUT", uploadURL, videoFile)
		if err != nil {
			return fmt.Errorf("failed to create request for upload: %v", err)
		}

		req.Header.Set("Content-Type", "video/mp4") // Set content type for MP4

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to upload video file: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body) // Read response body for debugging
			return fmt.Errorf("upload failed with status %s: %s", resp.Status, body)
		}
		log.Println("Asset uploaded successfully.")
		// Optional: Check asset status or retrieve IPFS link from asset metadata
		log.Printf("IPFS link will be available once processing is complete.")

		// from asset ID get the IPFS storage details
		assetID := res.Data.Asset.ID

		if err != nil {
			log.Fatal(err)
		}
		// get the IPFS link

		// Poll the status of the asset until the IPFS storage is completed

		for {
			assetRes, err := s.Asset.Get(ctx, assetID)
			if err != nil {
				log.Fatalf("Error fetching asset status: %v", err)
			}

			if assetRes.Asset.Storage == nil || assetRes.Asset.Storage.Status == nil {
				log.Println("Storage or Status is not available yet. Retrying...")
				time.Sleep(10 * time.Second) // Wait before polling again
				continue
			}

			storageStatus := assetRes.Asset.Storage.Status
			log.Printf("Current storage status phase: %s", storageStatus.Phase)
			if storageStatus.Phase == "ready" {
				log.Println("Store to IPFS completed!")
				break
			} else if storageStatus.Phase == "failed" {
				log.Fatal("Storage process failed.")
			}
			log.Println("Waiting for storage to complete...")
			// Wait before polling again (e.g., 10 seconds)
			time.Sleep(10 * time.Second)
		}

		// Once the storage is complete, get the IPFS metadata
		finalAssetRes, err := s.Asset.Get(ctx, assetID)
		if err != nil {
			log.Fatal("Failed to retrieve asset after storage completion: ", err)
		}

		// Print IPFS CID and URL
		if finalAssetRes.Asset.Storage != nil && finalAssetRes.Asset.Storage.Ipfs != nil && finalAssetRes.Asset.Storage.Ipfs.NftMetadata != nil {
			assetDownloadURL := finalAssetRes.Asset.GetDownloadURL()
			fmt.Printf("Asset URL link: %+v\n", *assetDownloadURL)

		} else {
			log.Println("No IPFS metadata found.")
		}
		return nil
	},
}
