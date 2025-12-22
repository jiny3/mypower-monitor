package library

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

type User struct {
	PushPlusClient `mapstructure:",squash"`
	Account        string `mapstructure:"account"`
	Password       string `mapstructure:"password"`
	RoomID         string `mapstructure:"room_id"`
}

func (u *User) Check() {
	tryCounter := 3

	for i := range tryCounter {
		dataValue, err := u.browserOps()
		if err != nil {
			if i == tryCounter-1 {
				break
			}
			logrus.WithField("retried", i).WithField("user", u.RoomID).WithError(err).Warn("查询电量失败, 重试中...")
			continue
		}
		logrus.WithField("user", u.RoomID).WithField("current-power", dataValue).Info("查询电量成功")
		err = u.Send(fmt.Sprintf("当前电量: %s", dataValue), fmt.Sprintf("http://157.0.19.2:10063/mypower/%s", u.RoomID))
		if err != nil {
			logrus.WithField("user", u.RoomID).WithError(err).Warn("发送电量通知失败")
		}

		Insert(Metric{
			GID:   u.RoomID,
			Key:   "power",
			Value: dataValue,
		})

		return
	}

	err := u.Send("查询电量失败", fmt.Sprintf("http://157.0.19.2:10063/mypower/%s", u.RoomID))
	if err != nil {
		logrus.WithField("user", u.RoomID).WithError(err).Warn("发送电量通知失败")
	}
}

// browser 自动化相关操作
func (u *User) browserOps() (string, error) {
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
		chromedp.SendKeys("#username", u.Account),
		chromedp.SendKeys("#password", u.Password),
		chromedp.Click("#login_submit", chromedp.NodeVisible),
		chromedp.AttributeValue(`//*[@name="REMAINEQ"]`, "data-value", &dataValue, nil),
	)
	if err != nil {
		return "", err
	}
	return dataValue, nil
}
