package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

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
		// 读取output目录
		// output直接下一层目录名作为Title
		// 读取这个目录下面的所有json文件, Labels是1s, 2s, 3s
		// 把json文件的内容转换成ChartData

		dir, err := os.ReadDir("./output")
		if err != nil {
			fmt.Println(err)
			return
		}

		var chartData []report.ChartData
		for _, d := range dir {
			switch d.Name() {
			case ".", "..":
				continue
			default:
				if !d.IsDir() {
					continue
				}
				var data report.ChartData
				dir2Str := path.Clean("./output/" + d.Name())
				dir2, err := os.ReadDir(dir2Str)
				if err != nil {
					fmt.Printf("os.ReadDir:", err)
					return
				}

				fmt.Printf("first dir:%s\n", dir2Str)
				maxNum := 0
				for _, d2 := range dir2 {
					curName := path.Clean(dir2Str + "/" + d2.Name())
					all, err := os.ReadFile(curName)
					if err != nil {
						fmt.Printf("read file:%v, dir(%s), name(%s)\n", err, dir2Str, d2.Name())
						return
					}

					var dataset report.Dataset
					err = json.Unmarshal(all, &dataset)
					if err != nil {
						fmt.Printf("json.Unmarshal:%v, dir(%s), name(%s), filename(%s)\n", err, dir2Str, d2.Name(), curName)
						return
					}
					maxNum = max(maxNum, len(dataset.Data))
					data.Datasets = append(data.Datasets, dataset)

				}

				data.Title = d.Name()
				for i := 1; i < maxNum; i++ {
					data.Labels = append(data.Labels, fmt.Sprintf("%ds", i))
				}
				chartData = append(chartData, data)
			}
		}
		// 示例数据
		// chartData := []report.ChartData{
		// 	{
		// 		Labels: []string{"一月", "二月", "三月", "四月", "五月", "六月", "七月"},
		// 		Datasets: []report.Dataset{
		// 			{
		// 				Label:       "线条 1",
		// 				Data:        []int{65, 59, 80, 81, 56, 55, 40},
		// 				BorderColor: "rgb(75, 192, 192)",
		// 				Tension:     0.1,
		// 			},
		// 			{
		// 				Label:       "线条 2",
		// 				Data:        []int{45, 69, 70, 61, 46, 65, 50},
		// 				BorderColor: "rgb(255, 99, 132)",
		// 				Tension:     0.1,
		// 			},
		// 			{
		// 				Label:       "线条 3",
		// 				Data:        []int{35, 49, 60, 71, 36, 45, 30},
		// 				BorderColor: "rgb(54, 162, 235)",
		// 				Tension:     0.1,
		// 			},
		// 		},
		// 	},
		// 	{
		// 		Labels: []string{"1s", "2s", "3s", "4s", "5s", "6s", "7s"},
		// 		Datasets: []report.Dataset{
		// 			{
		// 				Label:       "线条 4",
		// 				Data:        []int{65, 59, 80, 81, 56, 55, 40},
		// 				BorderColor: "rgb(75, 192, 192)",
		// 				Tension:     0.1,
		// 			},
		// 			{
		// 				Label:       "线条 5",
		// 				Data:        []int{45, 69, 70, 61, 46, 65, 50},
		// 				BorderColor: "rgb(255, 99, 132)",
		// 				Tension:     0.1,
		// 			},
		// 			{
		// 				Label:       "线条 6",
		// 				Data:        []int{35, 49, 60, 71, 36, 45, 30},
		// 				BorderColor: "rgb(54, 162, 235)",
		// 				Tension:     0.1,
		// 			},
		// 		},
		// 	},
		// 	{
		// 		Labels: []string{"1s", "2s", "3s", "4s", "5s", "6s", "7s"},
		// 		Datasets: []report.Dataset{
		// 			{
		// 				Label:       "线条 7",
		// 				Data:        []int{65, 59, 80, 81, 56, 55, 40},
		// 				BorderColor: "rgb(75, 192, 192)",
		// 				Tension:     0.1,
		// 			},
		// 			{
		// 				Label:       "线条 8",
		// 				Data:        []int{45, 69, 70, 61, 46, 65, 50},
		// 				BorderColor: "rgb(255, 99, 132)",
		// 				Tension:     0.1,
		// 			},
		// 			{
		// 				Label:       "线条 9",
		// 				Data:        []int{35, 49, 60, 71, 36, 45, 30},
		// 				BorderColor: "rgb(54, 162, 235)",
		// 				Tension:     0.1,
		// 			},
		// 		},
		// 	},
		// }

		// 返回JSON响应
		c.JSON(http.StatusOK, chartData)
		fmt.Printf("请求成功")
	})

	// 启动服务
	r.Run(":8082")
}
