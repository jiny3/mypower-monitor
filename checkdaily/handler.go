package checkdaily

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/jiny3/gopkg/filex"
	"github.com/sirupsen/logrus"
)

func (user *User) Check(token string) {
	tryCounter := 3

	for i := 0; i < tryCounter; i++ {
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()
		ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
		defer cancel()

		url := "http://ehall.njc.ucas.ac.cn/qljfwapp/sys/lwPsXykApp/index.do?#/dledcx"

		var dataValue string
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible("#username"),
			chromedp.WaitVisible("#password"),
			chromedp.SendKeys("#username", user.Account),
			chromedp.SendKeys("#password", user.Password),
			chromedp.Click("#login_submit", chromedp.NodeVisible),
			chromedp.AttributeValue(`//*[@name="REMAINEQ"]`, "data-value", &dataValue, nil),
		)
		if err != nil {
			logrus.WithField("user", user.Homeid).WithError(err).Error("查询电量失败")
		} else {
			logrus.WithField("user", user.Homeid).WithField("current-power", dataValue).Info("查询电量成功")
			user.send(token, fmt.Sprintf("当前电量: %s", dataValue), fmt.Sprintf("http://157.0.19.2:10063/mypower/%s", user.Homeid))
			appendFile(fmt.Sprintf("data/%s/value.txt", user.Homeid), fmt.Sprintf("%s\n", dataValue))
			appendFile(fmt.Sprintf("data/%s/time.txt", user.Homeid), fmt.Sprintf("%s\n", time.Now().Format("2006-01-02")))
			return
		}
	}

	time.Sleep(5 * time.Minute)

	user.send(token, "查询电量失败", fmt.Sprintf("http://157.0.19.2:10063/mypower/%s", user.Homeid))
}

// 向指定文件追加写入内容
func appendFile(filename string, content string) error {
	f, err := filex.FileCreate(filename)
	if err != nil {
		logrus.WithError(err).Error("打开文件失败")
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		logrus.WithError(err).Error("写入文件失败")
		return err
	}
	return nil
}

func (user *User) send(token, title, msg string) {
	if token == "" {
		logrus.Debug("未设置pushplus token")
		return
	}
	var data []byte
	if user.To == "" {
		data = []byte(fmt.Sprintf("{\"token\": \"%s\", \"title\": \"%s\", \"content\": \"%s\"}", token, title, msg))
	} else {
		data = []byte(fmt.Sprintf("{\"token\": \"%s\", \"title\": \"%s\", \"content\": \"%s\", \"to\": \"%s\"}", token, title, msg, user.To))
	}
	response, err := http.Post("http://www.pushplus.plus/send", "application/json", bytes.NewBuffer(data))
	if err != nil {
		logrus.WithField("user", user.Homeid).WithError(err).Errorf("发送邮件失败")
		return
	}
	defer response.Body.Close()
}
