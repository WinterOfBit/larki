package larki

import (
	"context"
	"sync"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

var GlobalClient *Client

func NewClient(appId, appSecret, verifyToken, encryptKey string) (*Client, error) {
	return NewClientWithConfig(&Config{
		AppID:       appId,
		AppSecret:   appSecret,
		VerifyToken: verifyToken,
		EncryptKey:  encryptKey,
	})
}

func NewClientWithConfig(config *Config) (*Client, error) {
	messageEventChan := make(chan *MessageEvent, 8)
	botAddedEventChan := make(chan *BotAddedEvent, 8)
	chatCreatedEventChan := make(chan *ChatCreatedEvent, 8)
	client := &Client{
		Config:       config,
		MessageEvent: messageEventChan,
	}

	client.Client = lark.NewClient(config.AppID, config.AppSecret)

	bot, err := client.GetBotInfo()
	if err != nil {
		return nil, err
	}

	client.BotInfo = bot
	client.EventDispatcher = dispatcher.NewEventDispatcher(client.VerifyToken, client.EncryptKey).
		OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			messageEventChan <- &MessageEvent{event.Event}
			return nil
		}).
		OnP2ChatMemberBotAddedV1(func(ctx context.Context, event *larkim.P2ChatMemberBotAddedV1) error {
			botAddedEventChan <- &BotAddedEvent{event.Event}
			return nil
		}).
		OnP1P2PChatCreatedV1(func(ctx context.Context, event *larkim.P1P2PChatCreatedV1) error {
			chatCreatedEventChan <- &ChatCreatedEvent{event.Event}
			return nil
		})

	return client, nil
}

var clientMu sync.Mutex

func SetGlobalClient(client *Client) {
	clientMu.Lock()
	defer clientMu.Unlock()
	GlobalClient = client
}
