package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

type SupabaseConfig struct {
	ProjectRef string
	ApiKey     string
}

type SupabaseHttp struct {
	config SupabaseConfig
	client *http.Client
}

func NewSupabaseHttpRepo(config SupabaseConfig) ObjectStorage {
	return &SupabaseHttp{
		config: config,
		client: &http.Client{},
	}
}

func (r *SupabaseHttp) Upload(ctx context.Context, fileName string, fileData []byte, bucketName string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	contentType := http.DetectContentType(fileData)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileName))
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, bytes.NewReader(fileData))
	if err != nil {
		return "", err
	}

	writer.Close()
	url := fmt.Sprintf("https://%s.supabase.co/storage/v1/object/%s/%s", r.config.ProjectRef, bucketName, fileName)

	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+r.config.ApiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-upsert", "true")

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload  with status %d: %s", resp.StatusCode, string(respBody))
	}

	publicURL := fmt.Sprintf("https://%s.supabase.co/storage/v1/object/public/%s/%s", r.config.ProjectRef, bucketName, fileName)
	return publicURL, nil
}

// Update updates an existing file in Supabase storage (uses PUT with upsert)
func (r *SupabaseHttp) Update(ctx context.Context, fileName string, fileData []byte, bucketName string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	contentType := http.DetectContentType(fileData)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileName))
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, bytes.NewReader(fileData))
	if err != nil {
		return "", err
	}

	writer.Close()
	url := fmt.Sprintf("https://%s.supabase.co/storage/v1/object/%s/%s", r.config.ProjectRef, bucketName, fileName)

	req, err := http.NewRequestWithContext(ctx, "PUT", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+r.config.ApiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-upsert", "true")

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("update failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	publicURL := fmt.Sprintf("https://%s.supabase.co/storage/v1/object/public/%s/%s", r.config.ProjectRef, bucketName, fileName)
	return publicURL, nil
}

// Delete removes a file from Supabase storage
func (r *SupabaseHttp) Delete(ctx context.Context, fileName string, bucketName string) error {
	url := fmt.Sprintf("https://%s.supabase.co/storage/v1/object/%s/%s", r.config.ProjectRef, bucketName, fileName)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+r.config.ApiKey)

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
