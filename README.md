# livepeer_ai_sumseg(Video Summary & Segmentation CLI Tool)

A command-line interface (CLI) tool to transcribe, summarize, subtitles and segment videos on Livepeer, built in Go.

## Features (Version v1.0)

- LivePeer AI apis, AWS and Hugging face AI
- Upload an asset using livepeer Api.
- Transcribe (captions) a video using ai.
- Summarize videos using ai.
- Subtitles for a video using ai.
- CLI-based, easy to integrate into existing workflows.
- Currently videos upto 2 mins are tested manually.

## Getting started

## Prerequisites

- Go installed (v1.21 or higher).
- Livepeer API key (Sign up at [Livepeer Studio](https://livepeer.com)).
- Hugging Face API key (for AI-driven text summarization).
- AWS access key and AWS secret key access (aws IAM - manage access to aws resources )

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/codebyyogesh/livepeer_ai_sumseg.git
   cd livepeer_ai_sumseg
   ```
2. Create a .env file ( or use an example .env.local file and rename it to .env). Ensure that you set up the following environment variables:

   | Variable Name | Description |

   -----------------------------------------------------------------------------|

   | `LP_AI_API_KEY` | Livepeer API key from [livepeer.studio](https://livepeer.studio) |

   | `HF_TEXT_SUMMARY_API_KEY` | Hugging Face token key from [Hugging Face](https://huggingface.co/) |

   | `AWS_ACCESS_KEY_ID` | AWS access key from Amazon AWS IAM (Identity Management https://aws.amazon.com/iam/) |

   | `AWS_SECRET_ACCESS_KEY` | AWS secret access key from Amazon AWS IAM |

   | `AWS_REGION` | The AWS region where your resources are located |

## Example `.env` File

```
LP_AI_API_KEY=""

HF_TEXT_SUMMARY_API_KEY=""

AWS_ACCESS_KEY_ID=""

AWS_SECRET_ACCESS_KEY=""

AWS_REGION=""

```

3.  Use the example aws_config_example.json and rename it to aws_config.json. Ensure that you set up the following aws variables:

| Variable Name | Description

----------------------------------------------------------------------------------- |

| `input_bucket_name` | Bucket name in aws s3 for input videos |

| `output_bucket_name` | Bucket name in aws s3 for output videos|

| `s3_input_video_path` | Folder name in aws s3 for storing input videos|

| `s3_output_transcription_path` | Folder name in aws s3 for storing output or processed videos|

## Example `aws_config` File

```
{
    "input_bucket_name": "livepeer",
    "output_bucket_name": "livepeer",
    "s3_input_video_path": "input/process.mp4",
    "s3_output_transcription_path": "transcribe/"
}
```

**PS**

**Bucket Creation**: You need to create an S3 bucket named livepeer. This bucket will be used for both input and output.

**Input Folder**: Inside the livepeer bucket, create a folder named input. Your input videos will always be copied as process.mp4 in this folder. Note that the CLI tool processes only one video at a time.

**Output Folder**: Create another folder named transcribe inside the livepeer bucket. This folder will be used to store the transcription output. Ensure that the folder name ends with a slash (/) in your configuration.

## Build and Run

From the project root folder, just use the make file to build

```
make build
```

This will build the binary (lpaisumseg) in the bin directory. To run the tool, you can run it as

```
$bin/lpaisumseg -h   --> for help
$bin/lpaisumseg caption -h --> help for caption cmd
```

To run a cmd (e.g)

```
$bin/lpaisumseg caption <VideoURL>  -->generates caption for the video
```

More details on the commands see wiki
https://github.com/codebyyogesh/livepeer_ai_sumseg/wiki

## Issues in the cli tool and livepeer AI apis.

1. Currently cli uploading asset uses polling mechanism to check the updated storage status to ipfs, in future this will be moved to livepeers webhook apis.

2. Once livepeer ai apis provide support for subtitles, summarize, segments for a video, the aws and hugging face apis can be replaced with livepeer ai apis.

3. The livepeer golang ai api image to video does not generate videos of larger length. Currently even experimenting with various params of the api like fps, num_inference_steps etc, yields videos of length only 3-4 secs.

## Upcoming features (v2.0 and above)

- Support Docker
- Support for segments in a video
- Use livepeer ai apis for captions, summarize and subtitles if supported
- Upload asset from a url
- Move to webhooks for getting status updates of storage.
- Support transcode
- Delete or remove an asset
- Multi Stream support
- Access control
- Metrics
- See if the cli features can also be available as a backend abstracted Apis.
