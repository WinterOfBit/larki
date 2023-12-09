# Larki 
[![github action](https://github.com/wintbiit/larki/actions/workflows/auto-release.yml/badge.svg?branch=master)](https://github.com/wintbiit/larki/actions/workflows/auto-release.yml)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/wintbiit/larki)](https://github.com/wintbiit/larki/releases)


## 对 [lark-oapi-go](https://github.com/larksuite/oapi-sdk-go) 的再封装
## Usage
### Init Client
```go
package larki
func NewClient(appId, appSecret, verifyToken, encryptKey string) (*Client, error) {
	return NewClientWithConfig(&Config{
		AppID:       appId,
		AppSecret:   appSecret,
		VerifyToken: verifyToken,
		EncryptKey:  encryptKey,
	})
}
```

Use it:
```go
package main

import (
    "fmt"
    "github.com/wintbiit/larki"
)

func main() {
    client, err := larki.NewClient("appId", "appSecret", "verifyToken", "encryptKey")
    if err != nil {
        panic(err)
    }
	
	// set global client
	larki.SetGlobalClient(client)
	
	// do something
}
```

### Send Message
```go
package main

import (
    "fmt"
    "github.com/wintbiit/larki"
)

var client *larki.Client

func main() {
	client.ReplyText(ctx, "om_v1234151", "hello world title", "hello world content")
}
```

### Subscribe Event
```go
package main

import (
    "fmt"
    "github.com/wintbiit/larki"
)

var client *larki.Client

func main() {
	// larkim.P2MessageReceiveV1Data
	for event := range client.MessageEvent {
		fmt.Printf("event: %+v\n", event)
	}
	
	// larkim.P2ChatMemberBotAddedV1Data
	for event := range client.BotAddedEvent {
        fmt.Printf("event: %+v\n", event)
    }
	
	// larkim.P1P2PChatCreatedV1Data
	for event := range client.ChatCreatedEvent {
        fmt.Printf("event: %+v\n", event)
    }
}
```

### Get Bot Meta
bot info is fetched when client is initialized
```go
var client *larki.Client
botInfo := client.BotInfo
```