package larki

import (
	"context"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// GetMessage 获取指定消息
func (c *Client) GetMessage(ctx context.Context, messageId string) (*larkim.Message, error) {
	resp, err := c.Im.Message.Get(ctx, larkim.NewGetMessageReqBuilder().MessageId(messageId).Build())
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, newLarkError(resp.Code, resp.Msg, "GetMessage")
	}

	return resp.Data.Items[0], nil
}

// GetMessage 获取指定消息
func GetMessage(ctx context.Context, messageId string) (*larkim.Message, error) {
	return GlobalClient.GetMessage(ctx, messageId)
}

// ReplyMessage 回复消息
func (c *Client) ReplyMessage(ctx context.Context, message, messageId, messageType string) error {
	resp, err := c.Im.Message.Reply(ctx,
		larkim.NewReplyMessageReqBuilder().Body(
			larkim.NewReplyMessageReqBodyBuilder().
				MsgType(messageType).
				Content(message).
				Build()).
			MessageId(messageId).Build())
	if err != nil {
		return err
	}

	if !resp.Success() {
		return newLarkError(resp.Code, resp.Msg, "ReplyMessage")
	}

	return nil
}

// ReplyMessage 回复消息
func ReplyMessage(ctx context.Context, message, messageId, messageType string) error {
	return GlobalClient.ReplyMessage(ctx, message, messageId, messageType)
}

// ReplyText 使用文本回复消息
func (c *Client) ReplyText(ctx context.Context, messageId, title string, text ...string) error {
	content, err := buildPost(title, text)
	if err != nil {
		return err
	}

	return c.ReplyMessage(ctx, content, messageId, larkim.MsgTypePost)
}

// ReplyText 使用文本回复消息
func ReplyText(ctx context.Context, messageId, title string, text ...string) error {
	return GlobalClient.ReplyText(ctx, messageId, title, text...)
}

// ReplyImage 使用图片回复消息
func (c *Client) ReplyImage(ctx context.Context, messageId, imageKey string) error {
	return c.ReplyMessage(ctx, NewImageContent(imageKey), messageId, larkim.MsgTypeImage)
}

// ReplyImage 使用图片回复消息
func ReplyImage(ctx context.Context, messageId, imageKey string) error {
	return GlobalClient.ReplyImage(ctx, messageId, imageKey)
}

// ReplyCard 使用卡片回复消息
func (c *Client) ReplyCard(ctx context.Context, messageId, card string) error {
	return c.ReplyMessage(ctx, card, messageId, larkim.MsgTypeInteractive)
}

// ReplyCard 使用卡片回复消息
func ReplyCard(ctx context.Context, messageId, card string) error {
	return GlobalClient.ReplyCard(ctx, messageId, card)
}

// ReplyCardTemplate 使用模板卡片回复消息
func (c *Client) ReplyCardTemplate(ctx context.Context, messageId, templateId string, vars map[string]interface{}) error {
	str, err := buildTemplateCard(templateId, vars)
	if err != nil {
		return err
	}

	return c.ReplyCard(ctx, messageId, str)
}

// ReplyCardTemplate 使用模板卡片回复消息
func ReplyCardTemplate(ctx context.Context, messageId, templateId string, vars map[string]interface{}) error {
	return GlobalClient.ReplyCardTemplate(ctx, messageId, templateId, vars)
}

func (c *Client) sendMessage(ctx context.Context, receiverIdType, message, receiveId, messageType string) (string, error) {
	resp, err := c.Im.Message.Create(ctx,
		larkim.NewCreateMessageReqBuilder().Body(
			larkim.NewCreateMessageReqBodyBuilder().
				MsgType(messageType).
				ReceiveId(receiveId).
				Content(message).
				Build()).
			ReceiveIdType(receiverIdType).Build())
	if err != nil {
		return "", err
	}

	if !resp.Success() {
		return "", newLarkError(resp.Code, resp.Msg, "SendMessage")
	}

	return *resp.Data.MessageId, nil
}

// SendMessage 发送消息
func SendMessage(ctx context.Context, receiverIdType, message, receiveId, messageType string) (string, error) {
	return GlobalClient.sendMessage(ctx, receiverIdType, message, receiveId, messageType)
}

// SendMessageToGroup 发送消息到群组
func (c *Client) SendMessageToGroup(ctx context.Context, groupId, message, messageType string) (string, error) {
	return c.sendMessage(ctx, larkim.ReceiveIdTypeChatId, message, groupId, messageType)
}

// SendMessageToGroup 发送消息到群组
func SendMessageToGroup(ctx context.Context, groupId, message, messageType string) (string, error) {
	return GlobalClient.SendMessageToGroup(ctx, groupId, message, messageType)
}

// SendTextToGroup 使用文本发送消息到群组
func (c *Client) SendTextToGroup(ctx context.Context, groupId, title string, text ...string) (string, error) {
	content, err := buildPost(title, text)
	if err != nil {
		return "", err
	}

	return c.SendMessageToGroup(ctx, groupId, content, larkim.MsgTypePost)
}

// SendTextToGroup 使用文本发送消息到群组
func SendTextToGroup(ctx context.Context, groupId, title string, text ...string) (string, error) {
	return GlobalClient.SendTextToGroup(ctx, groupId, title, text...)
}

// SendImageToGroup 使用图片发送消息到群组
func (c *Client) SendImageToGroup(ctx context.Context, groupId, imageKey string) (string, error) {
	return c.sendMessage(ctx, larkim.ReceiveIdTypeChatId, NewImageContent(imageKey), groupId, larkim.MsgTypeImage)
}

// SendImageToGroup 使用图片发送消息到群组
func SendImageToGroup(ctx context.Context, groupId, imageKey string) (string, error) {
	return GlobalClient.SendImageToGroup(ctx, groupId, imageKey)
}

// SendCardToGroup 使用卡片发送消息到群组
func (c *Client) SendCardToGroup(ctx context.Context, groupId, card string) (string, error) {
	return c.sendMessage(ctx, larkim.ReceiveIdTypeChatId, card, groupId, larkim.MsgTypeInteractive)
}

// SendCardToGroup 使用卡片发送消息到群组
func SendCardToGroup(ctx context.Context, groupId, card string) (string, error) {
	return GlobalClient.SendCardToGroup(ctx, groupId, card)
}

// SendCardTemplateToGroup 使用模板卡片发送消息到群组
func (c *Client) SendCardTemplateToGroup(ctx context.Context, groupId, templateId string, vars map[string]interface{}) (string, error) {
	str, err := buildTemplateCard(templateId, vars)
	if err != nil {
		return "", err
	}

	return c.SendCardToGroup(ctx, groupId, str)
}

// SendCardTemplateToGroup 使用模板卡片发送消息到群组
func SendCardTemplateToGroup(ctx context.Context, groupId, templateId string, vars map[string]interface{}) (string, error) {
	return GlobalClient.SendCardTemplateToGroup(ctx, groupId, templateId, vars)
}

// SendMessageToUser 发送消息到用户
func (c *Client) SendMessageToUser(ctx context.Context, openId, message, messageType string) (string, error) {
	return c.sendMessage(ctx, larkim.ReceiveIdTypeOpenId, message, openId, messageType)
}

// SendMessageToUser 发送消息到用户
func SendMessageToUser(ctx context.Context, openId, message, messageType string) (string, error) {
	return GlobalClient.SendMessageToUser(ctx, openId, message, messageType)
}

// SendTextToUser 使用文本发送消息到用户
func (c *Client) SendTextToUser(ctx context.Context, openId, title string, text ...string) (string, error) {
	content, err := buildPost(title, text)
	if err != nil {
		return "", err
	}

	return c.SendMessageToUser(ctx, openId, content, larkim.MsgTypePost)
}

// SendTextToUser 使用文本发送消息到用户
func SendTextToUser(ctx context.Context, openId, title string, text ...string) (string, error) {
	return GlobalClient.SendTextToUser(ctx, openId, title, text...)
}

// SendImageToUser 使用图片发送消息到用户
func (c *Client) SendImageToUser(ctx context.Context, openId, imageKey string) (string, error) {
	return c.sendMessage(ctx, larkim.ReceiveIdTypeOpenId, NewImageContent(imageKey), openId, larkim.MsgTypeImage)
}

// SendImageToUser 使用图片发送消息到用户
func SendImageToUser(ctx context.Context, openId, imageKey string) (string, error) {
	return GlobalClient.SendImageToUser(ctx, openId, imageKey)
}

// SendCardToUser 使用卡片发送消息到用户
func (c *Client) SendCardToUser(ctx context.Context, openId, card string) (string, error) {
	return c.sendMessage(ctx, larkim.ReceiveIdTypeOpenId, card, openId, larkim.MsgTypeInteractive)
}

// SendCardToUser 使用卡片发送消息到用户
func SendCardToUser(ctx context.Context, openId, card string) (string, error) {
	return GlobalClient.SendCardToUser(ctx, openId, card)
}

// SendCardTemplateToUser 使用模板卡片发送消息到用户
func (c *Client) SendCardTemplateToUser(ctx context.Context, openId, templateId string, vars map[string]interface{}) (string, error) {
	str, err := buildTemplateCard(templateId, vars)
	if err != nil {
		return "", err
	}

	return c.SendCardToUser(ctx, openId, str)
}

// SendCardTemplateToUser 使用模板卡片发送消息到用户
func SendCardTemplateToUser(ctx context.Context, openId, templateId string, vars map[string]interface{}) (string, error) {
	return GlobalClient.SendCardTemplateToUser(ctx, openId, templateId, vars)
}
