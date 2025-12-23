package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/jiny3/mypower-monitor/library"
)

// 数据解析
func GetRoomHistory(c *gin.Context) {
	roomid := c.Param("roomid")
	// 近 30 天数据（create_at）
	metrics := library.Select(library.WithWhere("gid = ? AND key = ?", roomid, "power"), library.WithWhere("created_at >= datetime('now', '-31 days')"), library.WithOrder("created_at"))
	if len(metrics) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "无数据",
		})
		return
	}

	timeList := make([]string, len(metrics))
	for i, m := range metrics {
		timeList[i] = m.CreatedAt.Format("2006-01-02")
	}

	valueFloat := make([]float64, len(metrics))
	for i, m := range metrics {
		value, _ := strconv.ParseFloat(m.Value, 64)
		valueFloat[i] = value
	}

	// 将数据转化为差分数据
	n := len(metrics)
	valueDiff := []string{}
	timeDiff := []string{}
	for i := 1; i < n; i++ {
		// 求相隔天数(向上取整)
		day := metrics[i].CreatedAt.Day() - metrics[i-1].CreatedAt.Day()
		if day <= 0 {
			continue
		}
		if day > 1 {
			timeList[i] = fmt.Sprintf("%s (近%d天)", timeList[i], day)
		}
		valueDiff = append(valueDiff, fmt.Sprintf("%.2f", valueFloat[i-1]-valueFloat[i]))
		timeDiff = append(timeDiff, timeList[i])
	}

	c.JSON(http.StatusOK, gin.H{
		"current": valueFloat[len(valueFloat)-1],
		"time":    timeDiff,
		"value":   valueDiff,
	})
}
