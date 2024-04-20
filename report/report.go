package report

type ChartData struct {
	Labels   []string  `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}

type Dataset struct {
	Label       string  `json:"label"`
	Data        []int   `json:"data"`
	BorderColor string  `json:"borderColor"`
	Tension     float64 `json:"tension"`
}
