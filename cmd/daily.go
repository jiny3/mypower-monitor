package cmd

import (
	"sync"

	"github.com/jiny3/gopkg/logx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/jiny3/mypower-monitor/library"
)

func checkCmd() *cobra.Command {
	localUser := library.User{}
	cmd := &cobra.Command{
		Use:   "check",
		Short: "执行电量查询脚本",
		Long:  "执行电量查询脚本, 支持从配置文件读取用户列表, 也支持通过命令行参数指定单个用户(声明后忽略配置文件)",
		PreRun: func(cmd *cobra.Command, args []string) {
			logx.InitLogrus(logx.WithLevel(logrus.DebugLevel))
			initUsersConf()

			if localUser.Account == "" && localUser.Password == "" && localUser.RoomID == "" {
				return
			}
			if localUser.Account == "" || localUser.Password == "" || localUser.RoomID == "" {
				cmd.Help()
				logrus.WithFields(logrus.Fields{
					"account":  localUser.Account,
					"password": localUser.Password,
					"room_id":  localUser.RoomID,
				}).Fatal("参数不完整")
			}
			usersConf.Users = []library.User{localUser}
		},
		Run: check,
	}
	cmd.Flags().StringVarP(&localUser.Account, "account", "a", "", "登录账号")
	cmd.Flags().StringVarP(&localUser.Password, "password", "p", "", "登录密码")
	cmd.Flags().StringVarP(&localUser.RoomID, "room_id", "r", "", "宿舍房间号")
	cmd.Flags().StringVarP(&localUser.PushPlusClient.Token, "pushplus_token", "t", "", "PushPlus 推送 Token")
	cmd.Flags().StringVarP(&localUser.PushPlusClient.To, "pushplus_to", "c", "", "PushPlus 推送目标")
	return cmd
}

func check(cmd *cobra.Command, args []string) {
	var wg sync.WaitGroup
	for _, user := range usersConf.Users {
		wg.Add(1)
		go func(u library.User) {
			defer wg.Done()
			u.Check()
		}(user)
	}
	wg.Wait()

	logrus.Info("all users checked")
}
