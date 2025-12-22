package service

import (
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

	// 将 value 数据转化为 差分数据
	valueDiff := make([]float64, len(valueFloat))
	for i := 1; i < len(valueFloat); i++ {
		valueDiff[i] = valueFloat[i-1] - valueFloat[i]
	}

	c.JSON(http.StatusOK, gin.H{
		"current": valueFloat[len(valueFloat)-1],
		"time":    timeList[1:],
		"value":   valueDiff[1:],
	})
}
