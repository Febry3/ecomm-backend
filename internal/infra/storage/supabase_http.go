package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type SupabaseHttp struct {
	projectRef string
	apiKey     string
	client     *http.Client
}

func NewSupabaseHttpRepo(projectRef, apiKey, bucketName string) ObjectStorage {
	return &SupabaseHttp{
		projectRef: projectRef,
		apiKey:     apiKey,
		client:     &http.Client{},
	}
}

func (r *SupabaseHttp) Upload(ctx context.Context, fileName string, fileData []byte, bucketName string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, bytes.NewReader(fileData))
	if err != nil {
		return "", err
	}

	writer.Close()
	url := fmt.Sprintf("https://%s.supabase.co/storage/v1/object/%s/%s", r.projectRef, bucketName, fileName)

	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-upsert", "true")

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	publicURL := fmt.Sprintf("https://%s.supabase.co/storage/v1/object/public/%s/%s", r.projectRef, bucketName, fileName)
	return publicURL, nil
}

// Delete implements ObjectStorage.
func (r *SupabaseHttp) Delete(ctx context.Context, fileName string) error {
	panic("unimplemented")
}
