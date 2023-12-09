package larki

import "github.com/bytedance/sonic"

func (e *CustomizedEvent) GetAsMenuEvent() (*MenuEventBody, error) {
	body := e.Body
	var menuEvent MenuEventBody
	if err := sonic.Unmarshal(body, &menuEvent); err != nil {
		return nil, err
	}

	return &menuEvent, nil
}
