package larki

import (
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

const larkErrFormat = "lark %s failed, code: %d, msg: %s"

// newLarkError 构造 飞书服务端报错
func newLarkError(code int, msg, field string) error {
	return fmt.Errorf(larkErrFormat, field, code, msg)
}

// ParseTextContent 解析文本消息内容
func ParseTextContent(text string) (string, bool) {
	var content textContent
	if err := sonic.UnmarshalString(text, &content); err != nil {
		return "", false
	}
	return content.Text, true
}

// ParseImageKey 解析图片消息imageKey
func ParseImageKey(context string) (string, bool) {
	var content imageContent
	if err := sonic.UnmarshalString(context, &content); err != nil {
		return "", false
	}
	return content.ImageKey, true
}

// NewImageContent 构造图片消息内容
func NewImageContent(imageKey string) string {
	content, _ := sonic.Marshal(imageContent{ImageKey: imageKey})
	return string(content)
}

// FilterTextContent 返回过滤掉 @ 信息后的文本内容和是否需要忽略，若包含@全体成员，则忽略，否则返回去除@信息后的文本内容
// @return text, atbot, atall
func (c *Client) FilterTextContent(text string, mentions []*larkim.MentionEvent) (string, bool, bool) {
	text = strings.TrimSpace(text)
	if len(mentions) == 0 {
		return text, false, false
	}

	if strings.Contains(text, "@_all") {
		return text, false, true
	}

	atBot := false
	for _, mention := range mentions {
		if mention.Key != nil {
			text = strings.ReplaceAll(text, *mention.Key, "")
			if *mention.Id.OpenId == c.BotInfo.OpenID {
				atBot = true
			}
		}
	}

	return strings.TrimSpace(text), atBot, false
}

// buildTemplateCard 构造模板卡片消息
func buildTemplateCard(templateId string, vars map[string]interface{}) (string, error) {
	template := templateCardContentData{
		TemplateId: templateId,
		Vars:       vars,
	}

	card := templateCardContent{
		Type: "template",
		Data: template,
	}

	str, err := sonic.MarshalString(card)
	if err != nil {
		return "", err
	}

	return str, nil
}

// buildPost 构造富文本消息
func buildPost(title string, content []string) (string, error) {
	post := larkim.NewMessagePostContent()
	post.ContentTitle(title)
	for _, t := range content {
		post.AppendContent([]larkim.MessagePostElement{
			&larkim.MessagePostText{
				Text: t,
			},
		})
	}

	var err error
	p, err := larkim.NewMessagePost().ZhCn(post).Build()
	if err != nil {
		return "", err
	}

	return p, nil
}
