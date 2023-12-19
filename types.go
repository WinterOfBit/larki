package larki

import (
	"context"
	"io"

	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	larkapplication "github.com/larksuite/oapi-sdk-go/v3/service/application/v6"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type Client struct {
	*lark.Client
	*Config
	*BotInfo
	EventDispatcher *dispatcher.EventDispatcher
	MessageClient
	ImageClient
}

type ClientOption func(*Client)

type Config struct {
	AppID       string
	AppSecret   string
	VerifyToken string
	EncryptKey  string
}

type MessageClient interface {
	GetMessage(ctx context.Context, messageId string) (*larkim.Message, error)
	ReplyMessage(ctx context.Context, message, messageId, messageType string) error
	ReplyText(ctx context.Context, messageId, title string, text ...string) error
	ReplyImage(ctx context.Context, messageId, imageKey string) error
	ReplyCard(ctx context.Context, messageId, card string) error
	ReplyCardTemplate(ctx context.Context, messageId, templateId string, vars map[string]interface{}) error
	SendMessage(ctx context.Context, receiverIdType, message, receiveId, messageType string) (string, error)
	SendMessageToGroup(ctx context.Context, groupId, message, messageType string) (string, error)
	SendTextToGroup(ctx context.Context, groupId, title string, text ...string) (string, error)
	SendImageToGroup(ctx context.Context, groupId, imageKey string) (string, error)
	SendCardToGroup(ctx context.Context, groupId, card string) (string, error)
	SendCardTemplateToGroup(ctx context.Context, groupId, templateId string, vars map[string]interface{}) (string, error)
	SendMessageToUser(ctx context.Context, openId, message, messageType string) (string, error)
	SendTextToUser(ctx context.Context, openId, title string, text ...string) (string, error)
	SendImageToUser(ctx context.Context, openId, imageKey string) (string, error)
	SendCardToUser(ctx context.Context, openId, card string) (string, error)
	SendCardTemplateToUser(ctx context.Context, openId, templateId string, vars map[string]interface{}) (string, error)
}

type ImageClient interface {
	GetImage(ctx context.Context, messageId, imageKey string) (io.Reader, error)
	UploadImage(ctx context.Context, reader io.Reader) (string, error)
}

type DocumentClient interface{}

type BotInfo struct {
	ActivateStatus int    `json:"activate_status"`
	AppName        string `json:"app_name"`
	AvatarUrl      string `json:"avatar_url"`
	OpenID         string `json:"open_id"`
}

type MessageEvent struct {
	*larkim.P2MessageReceiveV1Data
}

type BotAddedEvent struct {
	*larkim.P2ChatMemberBotAddedV1Data
}

type ChatCreatedEvent struct {
	*larkim.P1P2PChatCreatedV1Data
}

type MenuEvent struct {
	*larkapplication.P2BotMenuV6Data
}

type CustomizedEvent struct {
	*larkevent.EventReq
}

type botInfoResp struct {
	Code int     `json:"code"`
	Msg  string  `json:"msg"`
	Bot  BotInfo `json:"bot"`
}

type textContent struct {
	Text string `json:"text"`
}

type imageContent struct {
	ImageKey string `json:"image_key"`
}

type templateCardContent struct {
	Type string                  `json:"type"`
	Data templateCardContentData `json:"data"`
}

type templateCardContentData struct {
	TemplateId string                 `json:"template_id"`
	Vars       map[string]interface{} `json:"template_variable"`
}

type MenuEventBody struct {
	Operator  Operator `json:"operator"`
	EventKey  string   `json:"event_key"`
	Timestamp int64    `json:"timestamp"`
}

type Operator struct {
	OperatorName string        `json:"operator_name"`
	OperatorId   larkim.UserId `json:"operator_id"`
}

type JsApiTicketResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    JsApiTicket `json:"data"`
}

type JsApiTicket struct {
	Ticket   string `json:"ticket"`
	ExpireIn int    `json:"expire_in"`
}

type JsApiAuthResponse struct {
	Appid     interface{} `json:"appid"`
	Signature interface{} `json:"signature"`
	Noncestr  interface{} `json:"noncestr"`
	Timestamp interface{} `json:"timestamp"`
}

type MiniProgramTokenResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"msg"`
	Data    MiniProgramToken `json:"data"`
}

type MiniProgramToken struct {
	OpenId       string `json:"open_id"`
	EmployeeId   string `json:"employee_id"`
	SessionKey   string `json:"session_key"`
	TenantKey    string `json:"tenant_key"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
