package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

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
