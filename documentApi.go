package larki

import (
	"context"
	"fmt"
	"io"
	"time"

	larkwiki "github.com/larksuite/oapi-sdk-go/v3/service/wiki/v2"

	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"

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

func UpdateBaseRecord(ctx context.Context, baseId, tableId, recordId string, fields map[string]interface{}) error {
	return GlobalClient.UpdateBaseRecord(ctx, baseId, tableId, recordId, fields)
}

func (c *Client) GetRecords(ctx context.Context, baseId, tableId, viewId string, limit int) ([]*larkbitable.AppTableRecord, error) {
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
			AppToken(baseId).TableId(tableId).ViewId(viewId).PageSize(pageSize).PageToken(pageToken).Build()

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

func GetRecords(ctx context.Context, baseId, tableId, viewId string, limit int) ([]*larkbitable.AppTableRecord, error) {
	return GlobalClient.GetRecords(ctx, baseId, tableId, viewId, limit)
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

func GetRecord(ctx context.Context, baseId, tableId, recordId string) (*larkbitable.AppTableRecord, error) {
	return GlobalClient.GetRecord(ctx, baseId, tableId, recordId)
}

func (c *Client) GetDocMedia(ctx context.Context, fileToken string) (io.Reader, string, error) {
	resp, err := c.Drive.Media.Download(ctx, larkdrive.NewDownloadMediaReqBuilder().FileToken(fileToken).Build())
	if err != nil {
		return nil, "", err
	}

	if !resp.Success() {
		return nil, "", newLarkError(resp.Code, resp.Msg, "GetDocResource")
	}

	return resp.File, resp.FileName, nil
}

func GetDocMedia(ctx context.Context, fileToken string) (io.Reader, string, error) {
	return GlobalClient.GetDocMedia(ctx, fileToken)
}

func (c *Client) GetDocFile(ctx context.Context, fileToken string) (io.Reader, string, error) {
	resp, err := c.Drive.File.Download(ctx, larkdrive.NewDownloadFileReqBuilder().FileToken(fileToken).Build())
	if err != nil {
		return nil, "", err
	}

	if !resp.Success() {
		return nil, "", newLarkError(resp.Code, resp.Msg, "GetDocResource")
	}

	return resp.File, resp.FileName, nil
}

func GetDocFile(ctx context.Context, fileToken string) (io.Reader, string, error) {
	return GlobalClient.GetDocFile(ctx, fileToken)
}

func (c *Client) ListBaseTables(ctx context.Context, baseId string) ([]*larkbitable.AppTable, error) {
	resp, err := c.Bitable.AppTable.List(ctx, larkbitable.NewListAppTableReqBuilder().AppToken(baseId).Build())
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, newLarkError(resp.Code, resp.Msg, "ListBaseTables")
	}

	return resp.Data.Items, nil
}

func ListBaseTables(ctx context.Context, baseId string) ([]*larkbitable.AppTable, error) {
	return GlobalClient.ListBaseTables(ctx, baseId)
}

func (c *Client) UploadDocMedia(ctx context.Context, fileName, parentType, parentNode, extras string, size int, reader io.Reader) (string, error) {
	resp, err := c.Drive.Media.UploadAll(ctx, larkdrive.NewUploadAllMediaReqBuilder().Body(
		larkdrive.NewUploadAllMediaReqBodyBuilder().
			FileName(fileName).
			ParentType(parentType).
			ParentNode(parentNode).
			File(reader).
			Extra(extras).
			Size(size).Build()).Build())
	if err != nil {
		return "", err
	}

	if !resp.Success() {
		return "", newLarkError(resp.Code, resp.Msg, "UploadMedia")
	}

	return *resp.Data.FileToken, nil
}

func UploadDocMedia(ctx context.Context, fileName, parentType, parentNode, extras string, size int, reader io.Reader) (string, error) {
	return GlobalClient.UploadDocMedia(ctx, fileName, parentType, parentNode, extras, size, reader)
}

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

func (c *Client) ImportDoc(ctx context.Context, fileExt, fileToken, targetType, fileName string, mountType int, mountKey string) (string, error) {
	resp, err := c.Drive.ImportTask.Create(ctx, larkdrive.NewCreateImportTaskReqBuilder().
		ImportTask(larkdrive.NewImportTaskBuilder().
			FileExtension(fileExt).
			FileToken(fileToken).
			Type(targetType).
			Point(
				larkdrive.NewImportTaskMountPointBuilder().
					MountType(mountType).MountKey(mountKey).Build()).
			FileName(fileName).Build()).Build())
	if err != nil {
		return "", err
	}

	if !resp.Success() {
		return "", newLarkError(resp.Code, resp.Msg, "ImportDoc")
	}

	return *resp.Data.Ticket, nil
}

func ImportDoc(ctx context.Context, fileExt, fileToken, targetType, fileName string, mountType int, mountKey string) (string, error) {
	return GlobalClient.ImportDoc(ctx, fileExt, fileToken, targetType, fileName, mountType, mountKey)
}

func (c *Client) GetImportDocStatus(ctx context.Context, ticket string) (*larkdrive.ImportTask, error) {
	resp, err := c.Drive.ImportTask.Get(ctx, larkdrive.NewGetImportTaskReqBuilder().Ticket(ticket).Build())
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, newLarkError(resp.Code, resp.Msg, "GetImportDocStatus")
	}

	return resp.Data.Result, nil
}

func GetImportDocStatus(ctx context.Context, ticket string) (*larkdrive.ImportTask, error) {
	return GlobalClient.GetImportDocStatus(ctx, ticket)
}

func (c *Client) MoveDocToWiki(ctx context.Context, spaceId, objType, objToken, parentWikiToken string) (*larkwiki.MoveDocsToWikiSpaceNodeRespData, error) {
	resp, err := c.Wiki.SpaceNode.MoveDocsToWiki(ctx, larkwiki.NewMoveDocsToWikiSpaceNodeReqBuilder().
		Body(larkwiki.NewMoveDocsToWikiSpaceNodeReqBodyBuilder().
			ParentWikiToken(parentWikiToken).
			ObjToken(objToken).
			ObjType(objType).Build()).SpaceId(spaceId).Build())
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, newLarkError(resp.Code, resp.Msg, "MoveDocToWiki")
	}

	return resp.Data, nil
}

func MoveDocToWiki(ctx context.Context, spaceId, objType, objToken, parentWikiToken string) (*larkwiki.MoveDocsToWikiSpaceNodeRespData, error) {
	return GlobalClient.MoveDocToWiki(ctx, spaceId, objType, objToken, parentWikiToken)
}

func (c *Client) GetMoveDocToWikiStatus(ctx context.Context, taskId string) ([]*larkwiki.MoveResult, error) {
	resp, err := c.Wiki.Task.Get(ctx, larkwiki.NewGetTaskReqBuilder().
		TaskId(taskId).
		TaskType(larkwiki.TaskTypeMove).
		Build())
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, newLarkError(resp.Code, resp.Msg, "GetMoveDocToWikiStatus")
	}

	return resp.Data.Task.MoveResult, nil
}

func GetMoveDocToWikiStatus(ctx context.Context, taskId string) ([]*larkwiki.MoveResult, error) {
	return GlobalClient.GetMoveDocToWikiStatus(ctx, taskId)
}

// UploadToWiki 上传文件到知识库
func (c *Client) UploadToWiki(ctx context.Context,
	name, ext, docType, spaceId, parentNode string,
	size int, reader io.Reader,
) ([]*larkwiki.MoveResult, error) {
	extras := fmt.Sprintf(`{"file_extension":"%s", "obj_type": "%s"}`, ext, docType)
	fileToken, err := c.UploadDocMedia(ctx,
		fmt.Sprintf("%s-%d.%s", name, time.Now().UnixMilli(), ext),
		"ccm_import_open",
		"", extras, size, reader)
	if err != nil {
		return nil, err
	}

	ticket, err := c.ImportDoc(ctx, ext, fileToken, docType, name, 1, "")
	if err != nil {
		return nil, err
	}

	var status *larkdrive.ImportTask
	for {
		status, err = c.GetImportDocStatus(ctx, ticket)
		if err != nil {
			return nil, err
		}

		if *status.JobStatus == 0 {
			break
		} else if *status.JobStatus < 3 {
			time.Sleep(3 * time.Second)
		} else {
			return nil, newLarkError(*status.JobStatus, *status.JobErrorMsg, "UploadToWiki")
		}
	}

	resp, err := c.MoveDocToWiki(ctx, spaceId, docType, *status.Token, parentNode)
	if err != nil {
		return nil, err
	}

	var moveStatus []*larkwiki.MoveResult
	for {
		moveStatus, err = c.GetMoveDocToWikiStatus(ctx, *resp.TaskId)
		if err != nil {
			return nil, err
		}

		if len(moveStatus) > 0 {
			if *moveStatus[0].Status == 0 {
				break
			} else if *moveStatus[0].Status == 1 {
				time.Sleep(3 * time.Second)
			} else {
				return nil, newLarkError(*moveStatus[0].Status, *moveStatus[0].StatusMsg, "UploadToWiki")
			}
		}

		time.Sleep(3 * time.Second)
	}

	return moveStatus, nil
}

func UploadToWiki(ctx context.Context,
	name, ext, docType, spaceId, parentNode string,
	size int, reader io.Reader,
) ([]*larkwiki.MoveResult, error) {
	return GlobalClient.UploadToWiki(ctx, name, ext, docType, spaceId, parentNode, size, reader)
}
