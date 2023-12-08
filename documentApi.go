package larki

import (
	"context"
	"fmt"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
	"io"

	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

func (c *Client) UpdateBaseRecord(ctx context.Context, baseId, tableId, recordId string, fields map[string]interface{}) error {
	req := larkbitable.NewUpdateAppTableRecordReqBuilder().
		AppToken(baseId).TableId(tableId).RecordId(recordId).
		AppTableRecord(larkbitable.NewAppTableRecordBuilder().Fields(fields).Build()).Build()

	resp, err := c.Bitable.AppTableRecord.Update(ctx, req)
	if err != nil {
		return err
	}

	if !resp.Success() {
		return fmt.Errorf("failed to update base record %s/%s/%s: [%d] %s", baseId, tableId, recordId, resp.Code, resp.Msg)
	}

	return nil
}

func (c *Client) GetRecords(ctx context.Context, baseId, tableId string, limit int) ([]*larkbitable.AppTableRecord, error) {
	var records []*larkbitable.AppTableRecord
	var pageToken string
	pageSize := 50
	if limit > 0 && limit < pageSize {
		pageSize = limit
		records = make([]*larkbitable.AppTableRecord, 0, limit)
	} else {
		records = make([]*larkbitable.AppTableRecord, 0, pageSize)
	}

	for {
		req := larkbitable.NewListAppTableRecordReqBuilder().
			AppToken(baseId).TableId(tableId).PageSize(pageSize).PageToken(pageToken).Build()

		resp, err := c.Bitable.AppTableRecord.List(ctx, req)
		if err != nil {
			return nil, err
		}

		if !resp.Success() {
			return nil, fmt.Errorf("failed to get records %s/%s: [%d] %s", baseId, tableId, resp.Code, resp.Msg)
		}

		records = append(records, resp.Data.Items...)

		if resp.Data.HasMore == nil || !*resp.Data.HasMore || resp.Data.PageToken == nil {
			break
		}

		size := len(records)

		if limit > 0 && size >= limit {
			records = records[:limit]
			break
		}
	}

	return records, nil
}

func (c *Client) GetRecord(ctx context.Context, baseId, tableId, recordId string) (*larkbitable.AppTableRecord, error) {
	req := larkbitable.NewGetAppTableRecordReqBuilder().
		AppToken(baseId).TableId(tableId).RecordId(recordId).Build()

	resp, err := c.Bitable.AppTableRecord.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, newLarkError(resp.Code, resp.Msg, "GetRecord")
	}

	return resp.Data.Record, nil
}

func (c *Client) GetDocResource(ctx context.Context, fileToken string) (io.Reader, string, error) {
	resp, err := c.Drive.File.Download(ctx, larkdrive.NewDownloadFileReqBuilder().FileToken(fileToken).Build())
	if err != nil {
		return nil, "", err
	}

	if !resp.Success() {
		return nil, "", newLarkError(resp.Code, resp.Msg, "GetDocResource")
	}

	return resp.File, resp.FileName, nil
}
