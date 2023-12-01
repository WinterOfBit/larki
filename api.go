package larki

import (
	"context"

	"github.com/bytedance/sonic"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

const botInfoUrl = "https://open.feishu.cn/open-apis/bot/v3/info"

func (c *Client) getBotInfo() (*BotInfo, error) {
	resp, err := c.Get(context.Background(), botInfoUrl, nil, larkcore.AccessTokenTypeTenant)
	if err != nil {
		return nil, err
	}

	var data botInfoResp

	if err = sonic.Unmarshal(resp.RawBody, &data); err != nil {
		return nil, err
	}

	if data.Code != 0 {
		return nil, newLarkError(data.Code, data.Msg, "GetBotInfo")
	}

	return &data.Bot, nil
}

func (m *MessageEvent) TrimTextContent() (string, bool) {
	content, ok := ParseTextContent(*m.Message.Content)
	if !ok {
		return "", false
	}

	filtered, ignore := FilterTextContent(content, m.Message.Mentions)
	return filtered, ignore
}
