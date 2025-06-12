package url

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
)

// File represents a file to be uploaded.
type File struct {
	FileName string
	Reader   io.Reader
}

// ProgressCallback is called during file upload to report progress.
type (
	ProgressCallback func(uploaded, total int64)
	// ProgressReader wraps an io.Reader to track read progress.
	ProgressReader struct {
		reader   io.Reader
		total    int64
		uploaded int64
		callback ProgressCallback
		mu       sync.Mutex
	}
)

// NewProgressReader creates a new ProgressReader.
func NewProgressReader(reader io.Reader, total int64, callback ProgressCallback) *ProgressReader {
	return &ProgressReader{
		reader:   reader,
		total:    total,
		callback: callback,
	}
}

// Read implements io.Reader interface with progress tracking.
func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	if n > 0 {
		pr.mu.Lock()
		pr.uploaded += int64(n)
		uploaded := pr.uploaded
		pr.mu.Unlock()
		if pr.callback != nil {
			pr.callback(uploaded, pr.total)
		}
	}
	return n, err
}

// Progress returns the current upload progress.
func (pr *ProgressReader) Progress() (uploaded, total int64) {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	return pr.uploaded, pr.total
}

// Reset resets the progress counter.
func (pr *ProgressReader) Reset() {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	pr.uploaded = 0
}

// FileUpload represents a file upload with optional progress tracking.
type FileUpload struct {
	FieldName string
	FileName  string
	Reader    io.Reader
	Size      int64
	Progress  ProgressCallback
}

// NewFile creates a new File from a file path.
func NewFile(filePath string) (*File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return &File{FileName: filePath, Reader: file}, nil
}

// NewFileUpload creates a FileUpload from a file path.
func NewFileUpload(fieldName, filePath string) (*FileUpload, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}
	return &FileUpload{
		FieldName: fieldName,
		FileName:  filepath.Base(filePath),
		Reader:    file,
		Size:      stat.Size(),
	}, nil
}

// NewFileUploadFromReader creates a FileUpload from an io.Reader.
func NewFileUploadFromReader(fieldName, fileName string, reader io.Reader, size int64) *FileUpload {
	return &FileUpload{
		FieldName: fieldName,
		FileName:  fileName,
		Reader:    reader,
		Size:      size,
	}
}

// WithProgress sets a progress callback for the file upload.
func (fu *FileUpload) WithProgress(callback ProgressCallback) *FileUpload {
	fu.Progress = callback
	return fu
}

// AddFileToForm adds a file to a multipart form writer.
func AddFileToForm(writer *multipart.Writer, fieldName string, file *File) error {
	// 使用 filepath.Base 提取实际文件名
	part, err := writer.CreateFormFile(fieldName, filepath.Base(file.FileName))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file.Reader)
	return err
}

// AddFileUploadToForm adds a FileUpload to a multipart form writer with optional progress tracking.
func AddFileUploadToForm(writer *multipart.Writer, upload *FileUpload) error {
	part, err := writer.CreateFormFile(upload.FieldName, upload.FileName)
	if err != nil {
		return err
	}
	reader := upload.Reader
	if upload.Progress != nil && upload.Size > 0 {
		reader = NewProgressReader(upload.Reader, upload.Size, upload.Progress)
	}
	_, err = io.Copy(part, reader)
	return err
}

// MultipartEncoder encodes multiple files and form fields into multipart format.
type MultipartEncoder struct {
	files  []*FileUpload
	fields map[string]string
}

// NewMultipartEncoder creates a new MultipartEncoder.
func NewMultipartEncoder() *MultipartEncoder {
	return &MultipartEncoder{
		files:  make([]*FileUpload, 0),
		fields: make(map[string]string),
	}
}

// AddFile adds a file to the encoder.
func (me *MultipartEncoder) AddFile(upload *FileUpload) *MultipartEncoder {
	me.files = append(me.files, upload)
	return me
}

// AddField adds a form field to the encoder.
func (me *MultipartEncoder) AddField(name, value string) *MultipartEncoder {
	me.fields[name] = value
	return me
}

// Encode writes the multipart data to the writer and returns the content type.
func (me *MultipartEncoder) Encode(w io.Writer) (contentType string, err error) {
	writer := multipart.NewWriter(w)
	defer writer.Close()
	// Add form fields first
	for name, value := range me.fields {
		if err := writer.WriteField(name, value); err != nil {
			return "", err
		}
	}
	// Add files
	for _, upload := range me.files {
		if err := AddFileUploadToForm(writer, upload); err != nil {
			return "", err
		}
	}
	return writer.FormDataContentType(), nil
}
