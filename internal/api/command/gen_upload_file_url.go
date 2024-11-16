package command

import (
	"clipfy/internal/api/service"
)

type GenUploadFileURLCommand struct {
	storage *service.StorageService
}

type GenUploadFileURLCommandInput struct {
	FileName string `json:"file_name"`
}

type GenUploadFileURLCommandOutput struct {
	FileURL   string `json:"file_url"`
	UploadURL string `json:"upload_url"`
}

func NewUploadFileCommand(storage *service.StorageService) *GenUploadFileURLCommand {
	return &GenUploadFileURLCommand{
		storage: storage,
	}
}

func (u *GenUploadFileURLCommand) Execute(input *GenUploadFileURLCommandInput) *GenUploadFileURLCommandOutput {
	output, err := u.storage.GeneratePresignedUploadURL(&service.GeneratePresignedUploadURLInput{
		FileName: input.FileName,
	})
	if err != nil {
		panic(err)
	}

	return &GenUploadFileURLCommandOutput{
		FileURL:   output.FileURL,
		UploadURL: output.UploadURL,
	}
}
