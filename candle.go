package gota

import (
	"time"
)

// Candle - структура для хранения данных по свече
type Candle struct {
	StartTime  time.Time `json:"start_time"`
	OpenPrice  float64   `json:"open"`
	HighPrice  float64   `json:"high"`
	LowPrice   float64   `json:"low"`
	ClosePrice float64   `json:"close"`
	Volume     float64   `json:"volume"`
}

func NewCandle(startTime time.Time, open, high, low, close, volume float64) Candle {
	return Candle{
		StartTime:  startTime,
		OpenPrice:  open,
		HighPrice:  high,
		LowPrice:   low,
		ClosePrice: close,
		Volume:     volume,
	}
}

func (c Candle) GetStartTime() time.Time {
	return c.StartTime
}

func (c Candle) GetOpenPrice() float64 {
	return c.OpenPrice
}

func (c Candle) GetHighPrice() float64 {
	return c.HighPrice
}

func (c Candle) GetLowPrice() float64 {
	return c.LowPrice
}

func (c Candle) GetClosePrice() float64 {
	return c.ClosePrice
}

func (c Candle) GetVolume() float64 {
	return c.Volume
}
