package library

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type PushPlusClient struct {
	Token string `mapstructure:"token"`
	To    string `mapstructure:"to"`
}

func (p *PushPlusClient) Send(title, msg string) {
	if p.Token == "" {
		logrus.Error("pushplus token 为空")
		return
	}
	var data []byte
	if p.To == "" {
		data = fmt.Appendf(nil, "{\"token\": \"%s\", \"title\": \"%s\", \"content\": \"%s\"}", p.Token, title, msg)
	} else {
		data = fmt.Appendf(nil, "{\"token\": \"%s\", \"title\": \"%s\", \"content\": \"%s\", \"to\": \"%s\"}", p.Token, title, msg, p.To)
	}
	response, err := http.Post("http://www.pushplus.plus/send", "application/json", bytes.NewBuffer(data))
	if err != nil {
		logrus.WithError(err).Error("发送 PushPlus 消息失败")
		return
	}
	defer response.Body.Close()
}
