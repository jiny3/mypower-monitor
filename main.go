package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jiny3/gopkg/configx"
	"github.com/jiny3/gopkg/logx"
	"github.com/sirupsen/logrus"

	"github.com/jiny3/mypower-monitor/checkdaily"
	"github.com/jiny3/mypower-monitor/server"
)

var checkdailyYaml struct {
	Token string            `json:"token"`
	Users []checkdaily.User `json:"users"`
}

func init() {
	logx.InitLogrus(logx.WithOpsJSON("logs/ops.log"))
	err := configx.Read("config/userlist.yaml", &checkdailyYaml)
	if err != nil {
		logrus.WithError(err).Fatal("read userlist failed")
	}
	logrus.Debug("read userlist success")
}

func main() {
	ctlCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(yaml struct {
		Token string            `json:"token"`
		Users []checkdaily.User `json:"users"`
	}) {
		for _, user := range yaml.Users {
			go user.Check(yaml.Token)
		}
		for {
			select {
			case <-ctlCtx.Done():
				return
			case <-time.After(24 * time.Hour):
				for _, user := range yaml.Users {
					go user.Check(yaml.Token)
				}
			}
		}
	}(checkdailyYaml)

	port := 7001 // master的端口
	r := gin.Default()

	server.Init(r)

	logrus.WithField("port", port).Info("server running ...")
	r.Run(fmt.Sprintf(":%d", port))
}
