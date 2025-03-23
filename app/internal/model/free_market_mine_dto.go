package model

type FreeMarketMineResponse struct {
	SelledList []SelledEggDTO `json:"selled_list"`
	NoSellList []NoSellEggDTO `json:"no_sell_list"`
	Summary    SummaryDTO     `json:"summary"`
	Others     []NoSellEggDTO `json:"others"`
}

// SummaryDTO for summary
type SummaryDTO struct {
	TotalSelled int `json:"total_selled"`
	TotalNoSell int `json:"total_no_sell"`
}

// SelledEggDTO for sold eggs
type SelledEggDTO struct {
	EggPrice float64 `json:"egg_price"`
	EggNum   int     `json:"egg_num"`
	Date     string  `json:"date"`
	Total    float64 `json:"total"`
}

// NoSellEggDTO for unsold eggs
type NoSellEggDTO struct {
	Id       int     `json:"id"`
	EggPrice float64 `json:"egg_price"`
	EggNum   int     `json:"egg_num"`
	Date     string  `json:"date"`
}
