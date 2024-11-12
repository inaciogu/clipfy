package command

import (
	"clipfy/internal/api/service"
	"io"
)

type UploadFileCommand struct {
	storage *service.StorageService
}

type UploadFileCommandInput struct {
	FileName      string    `json:"file_name"`
	File          io.Reader `json:"file"`
	ContentType   string    `json:"content_type"`
	ContentLength int64     `json:"content_length"`
}

func NewUploadFileCommand(storage *service.StorageService) *UploadFileCommand {
	return &UploadFileCommand{
		storage: storage,
	}
}

func (u *UploadFileCommand) Execute(input *UploadFileCommandInput) {
	err := u.storage.UploadFile(&service.UploadFileInput{
		File:          input.File,
		FileName:      input.FileName,
		ContentLength: input.ContentLength,
	})
	if err != nil {
		panic(err)
	}
}
