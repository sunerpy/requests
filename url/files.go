package url

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type File struct {
	FileName string
	Reader   io.Reader
}

func NewFile(filePath string) (*File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return &File{FileName: filePath, Reader: file}, nil
}

func AddFileToForm(writer *multipart.Writer, fieldName string, file *File) error {
	// 使用 filepath.Base 提取实际文件名
	part, err := writer.CreateFormFile(fieldName, filepath.Base(file.FileName))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file.Reader)
	return err
}
