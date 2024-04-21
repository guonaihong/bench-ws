package report

type ChartData struct {
	Labels   []string  `json:"labels"`   // x轴的标签
	Datasets []Dataset `json:"datasets"` // y数据集
	Title    string    `json:"title"`    //折线图的标题
}

type Dataset struct {
	Label       string  `json:"label"`       // 数据集的名称
	Data        []int   `json:"data"`        // 数据集的数据
	BorderColor string  `json:"borderColor"` // 数据集的边框颜色
	Tension     float64 `json:"tension"`     // 数据集的曲线弯曲程度
}
