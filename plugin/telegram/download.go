package telegram

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// downloadFileId download file with fileID, return the filepath and blob.
func (b *Bot) downloadFileID(ctx context.Context, fileID string) (string, []byte, error) {
	file, err := b.GetFile(ctx, fileID)
	if err != nil {
		return "", nil, err
	}
	blob, err := b.downloadFilepath(ctx, file.FilePath)
	if err != nil {
		return "", nil, err
	}

	return file.FilePath, blob, nil
}

// downloadFilepath download file with filepath, you can get filepath by calling GetFile.
func (b *Bot) downloadFilepath(ctx context.Context, filePath string) ([]byte, error) {
	apiURL, err := b.apiURL(ctx)
	if err != nil {
		return nil, err
	}

	idx := strings.LastIndex(apiURL, "/bot")
	if idx < 0 {
		return nil, ErrInvalidToken
	}

	fileURL := apiURL[:idx] + "/file" + apiURL[idx:]

	resp, err := http.Get(fileURL + "/" + filePath)
	if err != nil {
		return nil, fmt.Errorf("fail to http.Get: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fail to io.ReadAll: %s", err)
	}

	return body, nil
}
