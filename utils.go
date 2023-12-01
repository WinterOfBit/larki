package larki

import (
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

const larkErrFormat = "lark %s failed, code: %d, msg: %s"

func newLarkError(code int, msg, field string) error {
	return fmt.Errorf(larkErrFormat, field, code, msg)
}

func ParseTextContent(text string) (string, bool) {
	var content textContent
	if err := sonic.UnmarshalString(text, &content); err != nil {
		return "", false
	}
	return content.Text, true
}

func ParseImageKey(context string) (string, bool) {
	var content imageContent
	if err := sonic.UnmarshalString(context, &content); err != nil {
		return "", false
	}
	return content.ImageKey, true
}

func NewImageContent(imageKey string) string {
	content, _ := sonic.Marshal(imageContent{ImageKey: imageKey})
	return string(content)
}

// FilterTextContent filters out mentions and returns the text content and should ignore it
func FilterTextContent(text string, mentions []*larkim.MentionEvent) (string, bool) {
	text = strings.TrimSpace(text)
	if len(mentions) == 0 {
		return text, true
	}

	if strings.Contains(text, "@_all") {
		return text, true
	}

	for _, mention := range mentions {
		if mention.Key != nil {
			text = strings.ReplaceAll(text, *mention.Key, "")
		}
	}

	return strings.TrimSpace(text), false
}

func buildTemplateCard(templateId string, vars map[string]interface{}) (string, error) {
	template := templateCardContentData{
		TemplateId:        templateId,
		TemplateVariables: vars,
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
