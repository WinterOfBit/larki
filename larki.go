package larki

import (
	"context"
	"os"
	"sync"

	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	larkapplication "github.com/larksuite/oapi-sdk-go/v3/service/application/v6"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

var GlobalClient *Client

func init() {
	if client, err := NewClientFromEnv(); err == nil {
		GlobalClient = client
	}
}

func NewClient(appId, appSecret, verifyToken, encryptKey string, options ...ClientOption) (*Client, error) {
	return NewClientWithConfig(&Config{
		AppID:       appId,
		AppSecret:   appSecret,
		VerifyToken: verifyToken,
		EncryptKey:  encryptKey,
	}, options...)
}

func NewClientFromEnv() (*Client, error) {
	return NewClientWithConfig(&Config{
		AppID:       os.Getenv("LARK_APP_ID"),
		AppSecret:   os.Getenv("LARK_APP_SECRET"),
		VerifyToken: os.Getenv("LARK_VERIFY_TOKEN"),
		EncryptKey:  os.Getenv("LARK_ENCRYPT_KEY"),
	})
}

func NewClientWithConfig(config *Config, options ...ClientOption) (*Client, error) {
	client := &Client{
		Config: config,
	}

	client.Client = lark.NewClient(config.AppID, config.AppSecret)

	bot, err := client.GetBotInfo()
	if err != nil {
		return nil, err
	}

	client.BotInfo = bot
	client.EventDispatcher = dispatcher.NewEventDispatcher(client.VerifyToken, client.EncryptKey)
	for _, option := range options {
		option(client)
	}

	return client, nil
}

var clientMu sync.Mutex

func SetGlobalClient(client *Client) {
	clientMu.Lock()
	defer clientMu.Unlock()
	GlobalClient = client
}

func WithMessageEventSubscribe(evtChan chan *MessageEvent) ClientOption {
	return func(client *Client) {
		client.EventDispatcher = client.EventDispatcher.OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			evtChan <- &MessageEvent{event.Event}
			return nil
		})
	}
}

func WithBotAddedEventSubscribe(evtChan chan *BotAddedEvent) ClientOption {
	return func(client *Client) {
		client.EventDispatcher = client.EventDispatcher.OnP2ChatMemberBotAddedV1(func(ctx context.Context, event *larkim.P2ChatMemberBotAddedV1) error {
			evtChan <- &BotAddedEvent{event.Event}
			return nil
		})
	}
}

func WithChatCreatedEventSubscribe(evtChan chan *ChatCreatedEvent) ClientOption {
	return func(client *Client) {
		client.EventDispatcher = client.EventDispatcher.OnP1P2PChatCreatedV1(func(ctx context.Context, event *larkim.P1P2PChatCreatedV1) error {
			evtChan <- &ChatCreatedEvent{event.Event}
			return nil
		})
	}
}

func WithMenuEventSubscribe(evtChan chan *MenuEvent) ClientOption {
	return func(client *Client) {
		client.EventDispatcher = client.EventDispatcher.OnP2BotMenuV6(func(ctx context.Context, event *larkapplication.P2BotMenuV6) error {
			evtChan <- &MenuEvent{event.Event}
			return nil
		})
	}
}

func WithCustomizedEventSubscribe(eventType string, evtChan chan *larkevent.EventReq) ClientOption {
	return func(client *Client) {
		client.EventDispatcher = client.EventDispatcher.OnCustomizedEvent(eventType, func(ctx context.Context, event *larkevent.EventReq) error {
			evtChan <- event
			return nil
		})
	}
}
