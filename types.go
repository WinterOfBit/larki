package larki

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type Client struct {
	*lark.Client
	*Config
	*BotInfo
	EventDispatcher *dispatcher.EventDispatcher
	MessageEvent    <-chan *MessageEvent
}

type Config struct {
	AppID       string
	AppSecret   string
	VerifyToken string
	EncryptKey  string
}

type BotInfo struct {
	ActivateStatus int    `json:"activate_status"`
	AppName        string `json:"app_name"`
	AvatarUrl      string `json:"avatar_url"`
	OpenID         string `json:"open_id"`
}

type MessageEvent struct {
	*larkim.P2MessageReceiveV1Data
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
	TemplateId        string                 `json:"template_id"`
	TemplateVariables map[string]interface{} `json:"template_variables"`
}
