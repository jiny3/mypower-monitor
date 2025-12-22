package cmd

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jiny3/gopkg/logx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/jiny3/mypower-monitor/service"
)

func serverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "启动电量查询服务",
		Long: "启动电量查询服务, 提供电量历史数据查询界面, " +
			"支持通过 /data/:roomid 接口查询指定宿舍的电量历史数据, " +
			"端口默认为 8080, 也可以通过命令行参数指定",
		Args: cobra.MaximumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			logx.InitLogrus(logx.WithOpsJSON("ops.log"))
			initUsersConf()
			logrus.Info("starting server...")
		},
		Run: runServer,
	}
	return cmd
}

func runServer(cmd *cobra.Command, args []string) {
	port := "8080"
	if len(args) > 0 {
		port = args[0]
	}
	r := gin.Default()

	service.Init(r)
	r.GET("/data/:roomid", service.GetRoomHistory)

	logrus.WithField("port", port).Info("server running ...")
	r.Run(fmt.Sprintf(":%s", port))
}
