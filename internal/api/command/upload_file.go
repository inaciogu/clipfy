package command

import (
	"clipfy/internal/api/service"
	"fmt"
	"io"
)

type UploadFileCommand struct {
	storage *service.StorageService
}

type UploadFileCommandInput struct {
	FileName      string    `json:"file_name"`
	File          io.Reader `json:"file"`
	ContentLength int64     `json:"content_length"`
}

func NewUploadFileCommand(storage *service.StorageService) *UploadFileCommand {
	return &UploadFileCommand{
		storage: storage,
	}
}

func (u *UploadFileCommand) Execute(input *UploadFileCommandInput) {
	uploadOutput, err := u.storage.UploadFile(&service.UploadFileInput{
		File:     input.File,
		FileName: input.FileName,
	})
	if err != nil {
		fmt.Println("Error uploading file")
		panic(err)
	}

	fmt.Println("File uploaded successfully")

	// save file url in database
	fmt.Println(uploadOutput.FileURL)
}
