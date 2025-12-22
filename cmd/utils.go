package cmd

import (
	"github.com/jiny3/gopkg/configx"
	"github.com/jiny3/gopkg/filex"
	"github.com/sirupsen/logrus"

	"github.com/jiny3/mypower-monitor/library"
)

var usersConf struct {
	Users []library.User `mapstructure:"users"`
}

func initUsersConf() {
	filePath := "users.toml"

	filex.FileCreate(filePath)
	err := configx.Read(filePath, &usersConf)
	if err != nil {
		logrus.WithError(err).Fatal("userlist config read failed")
		return
	}
	logrus.WithField("users", usersConf.Users).Debug("read userlist success")
}
