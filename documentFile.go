package larki

import (
	"context"
	"errors"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
	"go.uber.org/ratelimit"
	"io"
	"sync"
)

func (c *Client) UploadDocFile(ctx context.Context, name, parentType, parentNode string, size int, reader io.Reader) (string, error) {
	resp, err := c.Drive.File.UploadAll(ctx, larkdrive.NewUploadAllFileReqBuilder().Body(
		larkdrive.NewUploadAllFileReqBodyBuilder().
			FileName(name).
			ParentType(parentType).
			ParentNode(parentNode).
			File(reader).
			Size(size).Build()).Build())
	if err != nil {
		return "", err
	}

	if !resp.Success() {
		return "", newLarkError(resp.Code, resp.Msg, "UploadFile")
	}

	return *resp.Data.FileToken, nil
}

func UploadDocFile(ctx context.Context, name, parentType, parentNode string, size int, reader io.Reader) (string, error) {
	return GlobalClient.UploadDocFile(ctx, name, parentType, parentNode, size, reader)
}

func (c *Client) UploadDocFileMultiPart(ctx context.Context, name, parentNode string, size int, reader io.ReaderAt) (string, error) {
	req := larkdrive.NewUploadPrepareFileReqBuilder().
		FileUploadInfo(larkdrive.NewFileUploadInfoBuilder().
			FileName(name).
			ParentType("explorer").
			ParentNode(parentNode).
			Size(size).
			Build()).
		Build()

	resp, err := c.Drive.File.UploadPrepare(ctx, req)
	if err != nil {
		return "", err
	}

	if !resp.Success() {
		return "", newLarkError(resp.Code, resp.Msg, "UploadFilePrepare")
	}

	uploadId := *resp.Data.UploadId
	blockSize := *resp.Data.BlockSize
	blockNum := *resp.Data.BlockNum

	var wg sync.WaitGroup
	var mu sync.Mutex
	limit := ratelimit.New(5)
	wg.Add(blockNum)

	errs := make([]error, 0)

	// Upload blocks
	for i := 0; i < blockNum; i++ {
		i := i
		go func(blockSize int) {
			defer wg.Done()

			// reader range: [i*blockSize, (i+1)*blockSize)
			start := int64(i * blockSize)

			if i == blockNum-1 {
				blockSize = size - i*blockSize
			}

			end := start + int64(blockSize)

			reader := io.NewSectionReader(reader, start, end-start)

			err := c.uploadDocFilePart(ctx, uploadId, i, blockSize, reader, limit)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				return
			}
		}(blockSize)
	}

	wg.Wait()

	if len(errs) > 0 {
		return "", errors.Join(errs...)
	}

	closeReq := larkdrive.NewUploadFinishFileReqBuilder().
		Body(larkdrive.NewUploadFinishFileReqBodyBuilder().
			UploadId(uploadId).
			BlockNum(blockNum).
			Build()).
		Build()

	// 发起请求
	closeResp, err := c.Drive.File.UploadFinish(ctx, closeReq)
	if err != nil {
		return "", err
	}

	if !closeResp.Success() {
		return "", newLarkError(closeResp.Code, closeResp.Msg, "UploadFileFinish")
	}

	return *closeResp.Data.FileToken, nil
}

func (c *Client) uploadDocFilePart(ctx context.Context, uploadId string, i, blockSize int, reader io.Reader, limiter ratelimit.Limiter) error {
	limiter.Take()

	req := larkdrive.NewUploadPartFileReqBuilder().
		Body(larkdrive.NewUploadPartFileReqBodyBuilder().
			UploadId(uploadId).
			Seq(i).
			Size(blockSize).
			File(reader).
			Build()).
		Build()

	// 发起请求
	resp, err := c.Drive.File.UploadPart(ctx, req)
	if err != nil {
		return err
	}

	if !resp.Success() {
		return newLarkError(resp.Code, resp.Msg, "UploadFilePart")
	}

	return nil
}

func UploadDocFileMultiPart(ctx context.Context, name, parentNode string, size int, reader io.ReaderAt) (string, error) {
	return GlobalClient.UploadDocFileMultiPart(ctx, name, parentNode, size, reader)
}

func (c *Client) ListDriveFolder(ctx context.Context, folderToken string) ([]*larkdrive.File, error) {
	limiter := ratelimit.New(20)
	pageToken := ""

	files := make([]*larkdrive.File, 0)

	for {
		limiter.Take()
		resp, err := c.Drive.File.List(ctx, larkdrive.NewListFileReqBuilder().PageToken(pageToken).FolderToken(folderToken).Build())

		if err != nil {
			return nil, err
		}

		if !resp.Success() {
			return nil, newLarkError(resp.Code, resp.Msg, "ListDriveFolder")
		}

		files = append(files, resp.Data.Files...)

		if !*resp.Data.HasMore {
			break
		}

		pageToken = *resp.Data.NextPageToken
	}

	return files, nil
}

func (c *Client) CreateDriveFolder(ctx context.Context, name, folderToken string) (string, error) {
	req := larkdrive.NewCreateFolderFileReqBuilder().
		Body(larkdrive.NewCreateFolderFileReqBodyBuilder().
			Name(name).
			FolderToken(folderToken).
			Build()).
		Build()

	// 发起请求
	resp, err := c.Drive.File.CreateFolder(ctx, req)

	if err != nil {
		return "", err
	}

	if !resp.Success() {
		return "", newLarkError(resp.Code, resp.Msg, "CreateDriveFolder")
	}

	return *resp.Data.Token, nil
}
