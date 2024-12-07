package url

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"
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
