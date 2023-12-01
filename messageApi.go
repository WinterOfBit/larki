package larki

import (
	"context"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

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

func GetMessage(ctx context.Context, messageId string) (*larkim.Message, error) {
	return globalClient.GetMessage(ctx, messageId)
}

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

func ReplyMessage(ctx context.Context, message, messageId, messageType string) error {
	return globalClient.ReplyMessage(ctx, message, messageId, messageType)
}

func (c *Client) ReplyText(ctx context.Context, messageId, title string, text ...string) error {
	content, err := buildPost(title, text)
	if err != nil {
		return err
	}

	return c.ReplyMessage(ctx, content, messageId, larkim.MsgTypePost)
}

func ReplyText(ctx context.Context, messageId, title string, text ...string) error {
	return globalClient.ReplyText(ctx, messageId, title, text...)
}

func (c *Client) ReplyImage(ctx context.Context, messageId, imageKey string) error {
	return c.ReplyMessage(ctx, NewImageContent(imageKey), messageId, larkim.MsgTypeImage)
}

func ReplyImage(ctx context.Context, messageId, imageKey string) error {
	return globalClient.ReplyImage(ctx, messageId, imageKey)
}

func (c *Client) ReplyCard(ctx context.Context, messageId, card string) error {
	return c.ReplyMessage(ctx, card, messageId, larkim.MsgTypeInteractive)
}

func ReplyCard(ctx context.Context, messageId, card string) error {
	return globalClient.ReplyCard(ctx, messageId, card)
}

func (c *Client) ReplyCardTemplate(ctx context.Context, messageId, templateId string, vars map[string]interface{}) error {
	str, err := buildTemplateCard(templateId, vars)
	if err != nil {
		return err
	}

	return c.ReplyCard(ctx, messageId, str)
}

func ReplyCardTemplate(ctx context.Context, messageId, templateId string, vars map[string]interface{}) error {
	return globalClient.ReplyCardTemplate(ctx, messageId, templateId, vars)
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

func SendMessage(ctx context.Context, receiverIdType, message, receiveId, messageType string) (string, error) {
	return globalClient.sendMessage(ctx, receiverIdType, message, receiveId, messageType)
}

func (c *Client) SendMessageToGroup(ctx context.Context, groupId, message, messageType string) (string, error) {
	return c.sendMessage(ctx, larkim.ReceiveIdTypeChatId, message, groupId, messageType)
}

func SendMessageToGroup(ctx context.Context, groupId, message, messageType string) (string, error) {
	return globalClient.SendMessageToGroup(ctx, groupId, message, messageType)
}

func (c *Client) SendTextToGroup(ctx context.Context, groupId, title string, text ...string) (string, error) {
	content, err := buildPost(title, text)
	if err != nil {
		return "", err
	}

	return c.SendMessageToGroup(ctx, groupId, content, larkim.MsgTypePost)
}

func SendTextToGroup(ctx context.Context, groupId, title string, text ...string) (string, error) {
	return globalClient.SendTextToGroup(ctx, groupId, title, text...)
}

func (c *Client) SendImageToGroup(ctx context.Context, groupId, imageKey string) (string, error) {
	return c.sendMessage(ctx, larkim.ReceiveIdTypeChatId, NewImageContent(imageKey), groupId, larkim.MsgTypeImage)
}

func SendImageToGroup(ctx context.Context, groupId, imageKey string) (string, error) {
	return globalClient.SendImageToGroup(ctx, groupId, imageKey)
}

func (c *Client) SendCardToGroup(ctx context.Context, groupId, card string) (string, error) {
	return c.sendMessage(ctx, larkim.ReceiveIdTypeChatId, card, groupId, larkim.MsgTypeInteractive)
}

func SendCardToGroup(ctx context.Context, groupId, card string) (string, error) {
	return globalClient.SendCardToGroup(ctx, groupId, card)
}

func (c *Client) SendCardTemplateToGroup(ctx context.Context, groupId, templateId string, vars map[string]interface{}) (string, error) {
	str, err := buildTemplateCard(templateId, vars)
	if err != nil {
		return "", err
	}

	return c.SendCardToGroup(ctx, groupId, str)
}

func SendCardTemplateToGroup(ctx context.Context, groupId, templateId string, vars map[string]interface{}) (string, error) {
	return globalClient.SendCardTemplateToGroup(ctx, groupId, templateId, vars)
}

func (c *Client) SendMessageToUser(ctx context.Context, openId, message, messageType string) (string, error) {
	return c.sendMessage(ctx, larkim.ReceiveIdTypeOpenId, message, openId, messageType)
}

func SendMessageToUser(ctx context.Context, openId, message, messageType string) (string, error) {
	return globalClient.SendMessageToUser(ctx, openId, message, messageType)
}

func (c *Client) SendTextToUser(ctx context.Context, openId, title string, text ...string) (string, error) {
	content, err := buildPost(title, text)
	if err != nil {
		return "", err
	}

	return c.SendMessageToUser(ctx, openId, content, larkim.MsgTypePost)
}

func SendTextToUser(ctx context.Context, openId, title string, text ...string) (string, error) {
	return globalClient.SendTextToUser(ctx, openId, title, text...)
}

func (c *Client) SendImageToUser(ctx context.Context, openId, imageKey string) (string, error) {
	return c.sendMessage(ctx, larkim.ReceiveIdTypeOpenId, NewImageContent(imageKey), openId, larkim.MsgTypeImage)
}

func SendImageToUser(ctx context.Context, openId, imageKey string) (string, error) {
	return globalClient.SendImageToUser(ctx, openId, imageKey)
}

func (c *Client) SendCardToUser(ctx context.Context, openId, card string) (string, error) {
	return c.sendMessage(ctx, larkim.ReceiveIdTypeOpenId, card, openId, larkim.MsgTypeInteractive)
}

func SendCardToUser(ctx context.Context, openId, card string) (string, error) {
	return globalClient.SendCardToUser(ctx, openId, card)
}

func (c *Client) SendCardTemplateToUser(ctx context.Context, openId, templateId string, vars map[string]interface{}) (string, error) {
	str, err := buildTemplateCard(templateId, vars)
	if err != nil {
		return "", err
	}

	return c.SendCardToUser(ctx, openId, str)
}

func SendCardTemplateToUser(ctx context.Context, openId, templateId string, vars map[string]interface{}) (string, error) {
	return globalClient.SendCardTemplateToUser(ctx, openId, templateId, vars)
}
