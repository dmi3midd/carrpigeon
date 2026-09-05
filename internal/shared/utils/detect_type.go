package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	ErrUnsupportedFileType = errors.New("unsupported file type, only .html, .htm and .txt are allowed")
	ErrBinaryFile          = errors.New("binary files are not allowed")
)

// DetectTemplateType checks the file extension and content to determine if it is HTML.
// Returns isHTML (bool) and an error if the file is invalid or binary.
func DetectTemplateType(file multipart.File, filename string) (bool, error) {
	op := "utils.DetectTemplateType"

	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".html" && ext != ".htm" && ext != ".txt" {
		return false, fmt.Errorf("%s: %w (got %s)", op, ErrUnsupportedFileType, ext)
	}

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && !errors.Is(err, io.EOF) {
		return false, fmt.Errorf("%s: failed to read file header: %w", op, err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return false, fmt.Errorf("%s: failed to seek to start: %w", op, err)
	}

	contentType := http.DetectContentType(buffer[:n])
	// Ensure it's text, not binary
	if !strings.HasPrefix(contentType, "text/") && !strings.HasPrefix(contentType, "application/octet-stream") {
		return false, fmt.Errorf("%s: %w (detected %s)", op, ErrBinaryFile, contentType)
	}

	isHTML := ext == ".html" || ext == ".htm"
	return isHTML, nil
}

func DetectType(file multipart.File) (string, error) {
	op := "utils.DetectType"

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("%s: failed to read file header: %w", op, err)
	}

	// Повертаємо курсор на початок обов'язково
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("%s: failed to seek to start: %w", op, err)
	}

	return http.DetectContentType(buffer[:n]), nil
}
