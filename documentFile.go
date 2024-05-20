package larki

import (
	"context"
	"io"
	"strconv"
	"time"

	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
	"go.uber.org/ratelimit"
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

var (
	uploadDocFilePrepareLimit = ratelimit.New(2)
	uploadDocFileCloseLimit   = ratelimit.New(2)
)

func (c *Client) UploadDocFileMultiPart(ctx context.Context, name, parentNode string, size int, reader io.Reader) (string, error) {
	uploadDocFilePrepareLimit.Take()
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

	// Upload blocks
	for i := 0; i < blockNum; i++ {
		if i == blockNum-1 {
			blockSize = size - i*blockSize
		}

		reader := io.LimitReader(reader, int64(blockSize))

		err := c.uploadDocFilePart(ctx, uploadId, i, blockSize, reader)
		if err != nil {
			return "", err
		}
	}

	closeReq := larkdrive.NewUploadFinishFileReqBuilder().
		Body(larkdrive.NewUploadFinishFileReqBodyBuilder().
			UploadId(uploadId).
			BlockNum(blockNum).
			Build()).
		Build()

	// 发起请求
	uploadDocFileCloseLimit.Take()
	closeResp, err := c.Drive.File.UploadFinish(ctx, closeReq)
	if err != nil {
		return "", err
	}

	if !closeResp.Success() {
		return "", newLarkError(closeResp.Code, closeResp.Msg, "UploadFileFinish")
	}

	return *closeResp.Data.FileToken, nil
}

var uploadDocFilePartLimiter = ratelimit.New(5)

func (c *Client) uploadDocFilePart(ctx context.Context, uploadId string, i, blockSize int, reader io.Reader) error {
	uploadDocFilePartLimiter.Take()

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
		if resp.Code == 99991400 {
			// frequency limit
			ratelimitResetStr := resp.Header.Get("x-ogw-ratelimit-reset")
			ratelimitReset, err := strconv.ParseInt(ratelimitResetStr, 10, 64)
			if err != nil {
				return err
			}

			// wait for ratelimit reset
			time.Sleep(time.Duration(ratelimitReset) * time.Second)

			// retry
			return c.uploadDocFilePart(ctx, uploadId, i, blockSize, reader)
		}

		return newLarkError(resp.Code, resp.Msg, "UploadFilePart")
	}

	return nil
}

func UploadDocFileMultiPart(ctx context.Context, name, parentNode string, size int, reader io.Reader) (string, error) {
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

func ListDriveFolder(ctx context.Context, folderToken string) ([]*larkdrive.File, error) {
	return GlobalClient.ListDriveFolder(ctx, folderToken)
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

func CreateDriveFolder(ctx context.Context, name, folderToken string) (string, error) {
	return GlobalClient.CreateDriveFolder(ctx, name, folderToken)
}

var (
	uploadDocMediaMultiPartLimiter        = ratelimit.New(2)
	uploadDocMediaMultiPartPrepareLimiter = ratelimit.New(2)
)

func (c *Client) UploadDocMediaMultiPart(ctx context.Context, name, parentType, parentNode, extra string, size int, reader io.Reader) (string, error) {
	uploadDocMediaMultiPartPrepareLimiter.Take()

	req := larkdrive.NewUploadPrepareMediaReqBuilder().
		MediaUploadInfo(larkdrive.NewMediaUploadInfoBuilder().
			FileName(name).
			ParentType(parentType).
			Size(size).
			ParentNode(parentNode).
			Extra(extra).
			Build()).
		Build()

	resp, err := c.Drive.Media.UploadPrepare(ctx, req)
	if err != nil {
		return "", err
	}

	if !resp.Success() {
		return "", newLarkError(resp.Code, resp.Msg, "UploadMediaPrepare")
	}

	uploadId := *resp.Data.UploadId
	blockSize := *resp.Data.BlockSize
	blockNum := *resp.Data.BlockNum

	// Upload blocks
	for i := 0; i < blockNum; i++ {
		if i == blockNum-1 {
			blockSize = size - i*blockSize
		}

		reader := io.LimitReader(reader, int64(blockSize))

		err := c.uploadDocMediaPart(ctx, uploadId, i, blockSize, reader)
		if err != nil {
			return "", err
		}
	}

	closeReq := larkdrive.NewUploadFinishMediaReqBuilder().
		Body(larkdrive.NewUploadFinishMediaReqBodyBuilder().
			UploadId(uploadId).
			BlockNum(blockNum).
			Build()).
		Build()

	// 发起请求
	uploadDocMediaMultiPartLimiter.Take()
	closeResp, err := c.Drive.Media.UploadFinish(ctx, closeReq)
	if err != nil {
		return "", err
	}

	if !closeResp.Success() {
		return "", newLarkError(closeResp.Code, closeResp.Msg, "UploadMediaFinish")
	}

	return *closeResp.Data.FileToken, nil
}

func (c *Client) uploadDocMediaPart(ctx context.Context, uploadId string, i, blockSize int, reader io.Reader) error {
	uploadDocMediaMultiPartLimiter.Take()

	req := larkdrive.NewUploadPartMediaReqBuilder().
		Body(larkdrive.NewUploadPartMediaReqBodyBuilder().
			UploadId(uploadId).
			Seq(i).
			Size(blockSize).
			File(reader).
			Build()).
		Build()

	// 发起请求
	resp, err := c.Drive.Media.UploadPart(ctx, req)
	if err != nil {
		return err
	}

	if !resp.Success() {
		if resp.Code == 99991400 {
			// frequency limit
			ratelimitResetStr := resp.Header.Get("x-ogw-ratelimit-reset")
			ratelimitReset, err := strconv.ParseInt(ratelimitResetStr, 10, 64)
			if err != nil {
				return err
			}

			// wait for ratelimit reset
			time.Sleep(time.Duration(ratelimitReset) * time.Second)

			// retry
			return c.uploadDocMediaPart(ctx, uploadId, i, blockSize, reader)
		}

		return newLarkError(resp.Code, resp.Msg, "UploadMediaPart")
	}

	return nil
}

func UploadDocMediaMultiPart(ctx context.Context, name, parentType, parentNode, extra string, size int, reader io.Reader) (string, error) {
	return GlobalClient.UploadDocMediaMultiPart(ctx, name, parentType, parentNode, extra, size, reader)
}
