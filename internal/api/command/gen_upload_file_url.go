package command

import (
	"clipfy/internal/common"
	"fmt"
)

type GenUploadFileURLCommand struct {
	storage *common.StorageService
}

type GenUploadFileURLCommandInput struct {
	FileName string `json:"file_name"`
	UserID   string
}

type GenUploadFileURLCommandOutput struct {
	FileURL   string `json:"file_url"`
	UploadURL string `json:"upload_url"`
}

func NewUploadFileCommand(storage *common.StorageService) *GenUploadFileURLCommand {
	return &GenUploadFileURLCommand{
		storage: storage,
	}
}

func (u *GenUploadFileURLCommand) Execute(input *GenUploadFileURLCommandInput) *GenUploadFileURLCommandOutput {
	output, err := u.storage.GeneratePresignedUploadURL(&common.GeneratePresignedUploadURLInput{
		FileName: fmt.Sprintf("%s/%s", input.UserID, input.FileName),
	})
	if err != nil {
		panic(err)
	}

	return &GenUploadFileURLCommandOutput{
		FileURL:   output.FileURL,
		UploadURL: output.UploadURL,
	}
}
