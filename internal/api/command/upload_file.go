package command

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

type UploadFileCommand struct {
	awsConfig aws.Config
}

type UploadFileCommandInput struct {
	FileName      string `json:"file_name"`
	ContentType   string `json:"content_type"`
	ContentLength int64  `json:"content_length"`
}

func (u *UploadFileCommand) Execute(input *UploadFileCommandInput) {
}
