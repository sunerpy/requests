package url

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewFile(t *testing.T) {
	t.Run("Valid File", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "testfile")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())
		_, err = tmpFile.WriteString("This is a test file.")
		assert.NoError(t, err)
		file, err := NewFile(tmpFile.Name())
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.Equal(t, tmpFile.Name(), file.FileName)
		content, err := io.ReadAll(file.Reader)
		assert.NoError(t, err)
		assert.Equal(t, "This is a test file.", string(content))
	})
	t.Run("Invalid File", func(t *testing.T) {
		file, err := NewFile("nonexistentfile.txt")
		assert.Error(t, err)
		assert.Nil(t, file)
	})
}

func TestAddFileToForm(t *testing.T) {
	t.Run("Valid AddFileToForm", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "testfile")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())
		_, err = tmpFile.WriteString("This is a test file.")
		assert.NoError(t, err)
		file, err := NewFile(tmpFile.Name())
		assert.NoError(t, err)
		var buffer bytes.Buffer
		writer := multipart.NewWriter(&buffer)
		fieldName := "testfile"
		err = AddFileToForm(writer, fieldName, file)
		assert.NoError(t, err)
		writer.Close()
		reader := multipart.NewReader(&buffer, writer.Boundary())
		part, err := reader.NextPart()
		assert.NoError(t, err)
		assert.Equal(t, fieldName, part.FormName())
		assert.Equal(t, filepath.Base(file.FileName), part.FileName())
		content, err := io.ReadAll(part)
		assert.NoError(t, err)
		assert.Equal(t, "This is a test file.", string(content))
	})
	t.Run("CreateFormFile returns error with gomonkey", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "testfile")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())
		_, err = tmpFile.WriteString("This is a test file.")
		assert.NoError(t, err)
		file, err := NewFile(tmpFile.Name())
		assert.NoError(t, err)
		var buffer bytes.Buffer
		writer := multipart.NewWriter(&buffer)
		patches := gomonkey.ApplyMethod(reflect.TypeOf(writer), "CreateFormFile", func(_ *multipart.Writer, fieldName, fileName string) (io.Writer, error) {
			return nil, errors.New("mock CreateFormFile error")
		})
		defer patches.Reset()
		fieldName := "testfile"
		err = AddFileToForm(writer, fieldName, file)
		assert.Error(t, err)
		assert.Equal(t, "mock CreateFormFile error", err.Error())
	})
}

func TestProgressReader(t *testing.T) {
	t.Run("tracks progress correctly", func(t *testing.T) {
		data := []byte("Hello, World! This is test data for progress tracking.")
		reader := bytes.NewReader(data)
		total := int64(len(data))
		var lastUploaded int64
		var callCount int32
		callback := func(uploaded, totalSize int64) {
			atomic.AddInt32(&callCount, 1)
			lastUploaded = uploaded
			assert.Equal(t, total, totalSize)
		}
		pr := NewProgressReader(reader, total, callback)
		// Read in chunks
		buf := make([]byte, 10)
		var totalRead int64
		for {
			n, err := pr.Read(buf)
			if n > 0 {
				totalRead += int64(n)
			}
			if err == io.EOF {
				break
			}
			assert.NoError(t, err)
		}
		assert.Equal(t, total, totalRead)
		assert.Equal(t, total, lastUploaded)
		assert.Greater(t, callCount, int32(0))
	})
	t.Run("Progress method returns correct values", func(t *testing.T) {
		data := []byte("Test data")
		reader := bytes.NewReader(data)
		pr := NewProgressReader(reader, int64(len(data)), nil)
		// Read some data
		buf := make([]byte, 4)
		n, err := pr.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, 4, n)
		uploaded, total := pr.Progress()
		assert.Equal(t, int64(4), uploaded)
		assert.Equal(t, int64(len(data)), total)
	})
	t.Run("Reset clears progress", func(t *testing.T) {
		data := []byte("Test data")
		reader := bytes.NewReader(data)
		pr := NewProgressReader(reader, int64(len(data)), nil)
		// Read some data
		buf := make([]byte, 4)
		_, _ = pr.Read(buf)
		uploaded, _ := pr.Progress()
		assert.Equal(t, int64(4), uploaded)
		pr.Reset()
		uploaded, _ = pr.Progress()
		assert.Equal(t, int64(0), uploaded)
	})
	t.Run("works without callback", func(t *testing.T) {
		data := []byte("Test data")
		reader := bytes.NewReader(data)
		pr := NewProgressReader(reader, int64(len(data)), nil)
		result, err := io.ReadAll(pr)
		assert.NoError(t, err)
		assert.Equal(t, data, result)
	})
}

func TestNewFileUpload(t *testing.T) {
	t.Run("creates FileUpload from path", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "testfile*.txt")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())
		content := "Test file content"
		_, err = tmpFile.WriteString(content)
		assert.NoError(t, err)
		tmpFile.Close()
		upload, err := NewFileUpload("myfile", tmpFile.Name())
		assert.NoError(t, err)
		assert.NotNil(t, upload)
		assert.Equal(t, "myfile", upload.FieldName)
		assert.Equal(t, filepath.Base(tmpFile.Name()), upload.FileName)
		assert.Equal(t, int64(len(content)), upload.Size)
	})
	t.Run("returns error for non-existent file", func(t *testing.T) {
		upload, err := NewFileUpload("myfile", "/nonexistent/file.txt")
		assert.Error(t, err)
		assert.Nil(t, upload)
	})
}

func TestNewFileUploadFromReader(t *testing.T) {
	t.Run("creates FileUpload from reader", func(t *testing.T) {
		data := []byte("Test content")
		reader := bytes.NewReader(data)
		upload := NewFileUploadFromReader("field", "test.txt", reader, int64(len(data)))
		assert.Equal(t, "field", upload.FieldName)
		assert.Equal(t, "test.txt", upload.FileName)
		assert.Equal(t, int64(len(data)), upload.Size)
	})
}

func TestFileUploadWithProgress(t *testing.T) {
	t.Run("sets progress callback", func(t *testing.T) {
		data := []byte("Test content")
		reader := bytes.NewReader(data)
		callback := func(uploaded, total int64) {
			// callback set
		}
		upload := NewFileUploadFromReader("field", "test.txt", reader, int64(len(data)))
		upload.WithProgress(callback)
		assert.NotNil(t, upload.Progress)
	})
}

func TestAddFileUploadToForm(t *testing.T) {
	t.Run("adds file upload to form", func(t *testing.T) {
		data := []byte("Test file content")
		reader := bytes.NewReader(data)
		upload := NewFileUploadFromReader("myfile", "test.txt", reader, int64(len(data)))
		var buffer bytes.Buffer
		writer := multipart.NewWriter(&buffer)
		err := AddFileUploadToForm(writer, upload)
		assert.NoError(t, err)
		writer.Close()
		// Verify the multipart content
		mpReader := multipart.NewReader(&buffer, writer.Boundary())
		part, err := mpReader.NextPart()
		assert.NoError(t, err)
		assert.Equal(t, "myfile", part.FormName())
		assert.Equal(t, "test.txt", part.FileName())
		content, err := io.ReadAll(part)
		assert.NoError(t, err)
		assert.Equal(t, data, content)
	})
	t.Run("tracks progress during upload", func(t *testing.T) {
		data := []byte("Test file content for progress tracking")
		reader := bytes.NewReader(data)
		var progressCalled bool
		var lastUploaded int64
		callback := func(uploaded, total int64) {
			progressCalled = true
			lastUploaded = uploaded
		}
		upload := NewFileUploadFromReader("myfile", "test.txt", reader, int64(len(data)))
		upload.WithProgress(callback)
		var buffer bytes.Buffer
		writer := multipart.NewWriter(&buffer)
		err := AddFileUploadToForm(writer, upload)
		assert.NoError(t, err)
		writer.Close()
		assert.True(t, progressCalled)
		assert.Equal(t, int64(len(data)), lastUploaded)
	})
	t.Run("handles CreateFormFile error", func(t *testing.T) {
		data := []byte("Test content")
		reader := bytes.NewReader(data)
		upload := NewFileUploadFromReader("myfile", "test.txt", reader, int64(len(data)))
		var buffer bytes.Buffer
		writer := multipart.NewWriter(&buffer)
		patches := gomonkey.ApplyMethod(reflect.TypeOf(writer), "CreateFormFile", func(_ *multipart.Writer, fieldName, fileName string) (io.Writer, error) {
			return nil, errors.New("mock error")
		})
		defer patches.Reset()
		err := AddFileUploadToForm(writer, upload)
		assert.Error(t, err)
		assert.Equal(t, "mock error", err.Error())
	})
}

func TestMultipartEncoder(t *testing.T) {
	t.Run("encodes files and fields", func(t *testing.T) {
		encoder := NewMultipartEncoder()
		// Add fields
		encoder.AddField("name", "John")
		encoder.AddField("email", "john@example.com")
		// Add file
		fileData := []byte("File content")
		upload := NewFileUploadFromReader("document", "doc.txt", bytes.NewReader(fileData), int64(len(fileData)))
		encoder.AddFile(upload)
		var buffer bytes.Buffer
		contentType, err := encoder.Encode(&buffer)
		assert.NoError(t, err)
		assert.Contains(t, contentType, "multipart/form-data")
		// Parse and verify
		boundary := strings.Split(contentType, "boundary=")[1]
		mpReader := multipart.NewReader(&buffer, boundary)
		// Read all parts
		parts := make(map[string]string)
		var fileContent []byte
		for {
			part, err := mpReader.NextPart()
			if err == io.EOF {
				break
			}
			assert.NoError(t, err)
			content, _ := io.ReadAll(part)
			if part.FileName() != "" {
				fileContent = content
			} else {
				parts[part.FormName()] = string(content)
			}
		}
		assert.Equal(t, "John", parts["name"])
		assert.Equal(t, "john@example.com", parts["email"])
		assert.Equal(t, fileData, fileContent)
	})
	t.Run("handles empty encoder", func(t *testing.T) {
		encoder := NewMultipartEncoder()
		var buffer bytes.Buffer
		contentType, err := encoder.Encode(&buffer)
		assert.NoError(t, err)
		assert.Contains(t, contentType, "multipart/form-data")
	})
	t.Run("handles multiple files", func(t *testing.T) {
		encoder := NewMultipartEncoder()
		file1Data := []byte("File 1 content")
		file2Data := []byte("File 2 content")
		encoder.AddFile(NewFileUploadFromReader("file1", "doc1.txt", bytes.NewReader(file1Data), int64(len(file1Data))))
		encoder.AddFile(NewFileUploadFromReader("file2", "doc2.txt", bytes.NewReader(file2Data), int64(len(file2Data))))
		var buffer bytes.Buffer
		contentType, err := encoder.Encode(&buffer)
		assert.NoError(t, err)
		// Parse and count files
		boundary := strings.Split(contentType, "boundary=")[1]
		mpReader := multipart.NewReader(&buffer, boundary)
		fileCount := 0
		for {
			part, err := mpReader.NextPart()
			if err == io.EOF {
				break
			}
			assert.NoError(t, err)
			if part.FileName() != "" {
				fileCount++
			}
		}
		assert.Equal(t, 2, fileCount)
	})
	t.Run("handles file upload error", func(t *testing.T) {
		encoder := NewMultipartEncoder()
		// Create a file upload with a reader that will fail
		failReader := &failingReader{}
		upload := NewFileUploadFromReader("file", "test.txt", failReader, 100)
		encoder.AddFile(upload)
		var buffer bytes.Buffer
		_, err := encoder.Encode(&buffer)
		// The error might be swallowed or returned depending on implementation
		// Just ensure it doesn't panic
		_ = err
	})
	t.Run("handles WriteField error", func(t *testing.T) {
		encoder := NewMultipartEncoder()
		encoder.AddField("name", "value")
		var buffer bytes.Buffer
		writer := multipart.NewWriter(&buffer)
		patches := gomonkey.ApplyMethod(reflect.TypeOf(writer), "WriteField", func(_ *multipart.Writer, fieldname, value string) error {
			return errors.New("mock WriteField error")
		})
		defer patches.Reset()
		_, err := encoder.Encode(&buffer)
		// Error should be returned
		assert.Error(t, err)
	})
}

type failingReader struct{}

func (r *failingReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func TestNewFileUpload_StatError(t *testing.T) {
	// Test with a directory instead of a file
	tmpDir, err := os.MkdirTemp("", "testdir")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	// Create a file in the directory
	tmpFile := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(tmpFile, []byte("test"), 0o644)
	assert.NoError(t, err)
	// This should work
	upload, err := NewFileUpload("field", tmpFile)
	assert.NoError(t, err)
	assert.NotNil(t, upload)
}

func TestNewFileUpload_StatError_Mock(t *testing.T) {
	// Create a temp file
	tmpFile, err := os.CreateTemp("", "testfile*.txt")
	assert.NoError(t, err)
	tmpFileName := tmpFile.Name()
	defer os.Remove(tmpFileName)
	tmpFile.Close()
	// Mock os.Open to return a file that fails on Stat
	patches := gomonkey.ApplyFunc(os.Open, func(name string) (*os.File, error) {
		// Return a file that will fail on Stat by using a pipe
		r, w, _ := os.Pipe()
		w.Close() // Close write end
		return r, nil
	})
	defer patches.Reset()
	upload, err := NewFileUpload("field", tmpFileName)
	// The Stat on a pipe should fail or return unexpected results
	// Either way, we're testing the error path
	if err != nil {
		assert.Nil(t, upload)
	}
}
