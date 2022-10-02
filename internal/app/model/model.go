package model

type Data struct {
	Ethereum Ethereum `json:"ethereum"`
}

type Ethereum struct {
	Transactions []Transactions `json:"transactions"`
}

type Transactions struct {
	Time           string  `json:"time"`
	GasPrice       float64 `json:"gasPrice"`
	GasValue       float64 `json:"gasValue"`
	Average        float64 `json:"average"`
	MaxGasPrice    float64 `json:"maxGasPrice"`
	MedianGasPrice float64 `json:"medianGasPrice"`
}

type ProcessedGasInfo struct {
	MonthsGasValue       map[string]float64 `json:"monthGasValue,omitempty"`
	DayAverageGasPrice   float64            `json:"dayAverageGasPrice,omitempty"`
	HoursAverageGasPrice map[int]float64    `json:"hourAverageGasPrice,omitempty"`
	AllTimePaid          float64            `json:"allTimePaid,omitempty"`
}
