package larki

import (
	"bytes"
	"context"
	"io"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func (c *Client) GetImage(ctx context.Context, messageId, imageKey string) (io.Reader, error) {
	resp, err := c.Im.MessageResource.Get(ctx, larkim.NewGetMessageResourceReqBuilder().MessageId(messageId).FileKey(imageKey).Build())
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, newLarkError(resp.Code, resp.Msg, "GetImage")
	}

	reader := bytes.NewReader(resp.RawBody)
	return reader, nil
}

func (c *Client) UploadImage(ctx context.Context, reader io.Reader) (string, error) {
	resp, err := c.Im.Image.Create(ctx, larkim.NewCreateImageReqBuilder().Body(
		larkim.NewCreateImageReqBodyBuilder().ImageType(larkim.ImageTypeMessage).Image(reader).Build()).Build())
	if err != nil {
		return "", err
	}

	if !resp.Success() {
		return "", newLarkError(resp.Code, resp.Msg, "UploadImage")
	}

	return *resp.Data.ImageKey, nil
}
