package helpers

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

func GetFileFromContext(c *gin.Context, fieldName string) ([]byte, error) {
	fileHeader, err := c.FormFile(fieldName)
	if err != nil {
		return nil, fmt.Errorf("no file uploaded: %w", err)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read file bytes: %w", err)
	}

	return fileBytes, nil
}
