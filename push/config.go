package push

type config struct {
	Db        *string `json:"db"`
	Interval  *int    `json:"interval"`
	BatchSize *int    `json:"batchSize"`
}

var Config config
