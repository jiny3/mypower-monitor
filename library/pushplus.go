package library

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

type PushPlusClient struct {
	Token string `mapstructure:"token"`
	To    string `mapstructure:"to"`
}

func (p *PushPlusClient) Send(title, msg string) error {
	if p.Token == "" {
		return errors.New("pushplus token 为空")
	}
	var data []byte
	if p.To == "" {
		data = fmt.Appendf(nil, "{\"token\": \"%s\", \"title\": \"%s\", \"content\": \"%s\"}", p.Token, title, msg)
	} else {
		data = fmt.Appendf(nil, "{\"token\": \"%s\", \"title\": \"%s\", \"content\": \"%s\", \"to\": \"%s\"}", p.Token, title, msg, p.To)
	}
	response, err := http.Post("http://www.pushplus.plus/send", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}
