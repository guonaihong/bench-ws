package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guonaihong/bench-ws/report"
)

func main() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源访问
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})
	// 定义路由处理程序
	r.GET("/chartData", func(c *gin.Context) {
		// 示例数据
		chartData := []report.ChartData{
			{
				Labels: []string{"一月", "二月", "三月", "四月", "五月", "六月", "七月"},
				Datasets: []report.Dataset{
					{
						Label:       "线条 1",
						Data:        []int{65, 59, 80, 81, 56, 55, 40},
						BorderColor: "rgb(75, 192, 192)",
						Tension:     0.1,
					},
					{
						Label:       "线条 2",
						Data:        []int{45, 69, 70, 61, 46, 65, 50},
						BorderColor: "rgb(255, 99, 132)",
						Tension:     0.1,
					},
					{
						Label:       "线条 3",
						Data:        []int{35, 49, 60, 71, 36, 45, 30},
						BorderColor: "rgb(54, 162, 235)",
						Tension:     0.1,
					},
				},
			},
		}

		// 返回JSON响应
		c.JSON(http.StatusOK, chartData)
		fmt.Printf("请求成功")
	})

	// 启动服务
	r.Run(":8082")
}
