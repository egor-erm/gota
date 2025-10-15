package gota

import (
	"time"
)

// Candle - структура для хранения данных по свече
type Candle struct {
	startTime  time.Time
	openPrice  float64
	highPrice  float64
	lowPrice   float64
	closePrice float64
	volume     float64
}

func NewCandle(startTime time.Time, open, high, low, close, volume float64) Candle {
	return Candle{
		startTime:  startTime,
		openPrice:  open,
		highPrice:  high,
		lowPrice:   low,
		closePrice: close,
		volume:     volume,
	}
}

func (c Candle) GetStartTime() time.Time {
	return c.startTime
}

func (c Candle) GetOpenPrice() float64 {
	return c.openPrice
}

func (c Candle) GetHighPrice() float64 {
	return c.highPrice
}

func (c Candle) GetLowPrice() float64 {
	return c.lowPrice
}

func (c Candle) GetClosePrice() float64 {
	return c.closePrice
}

func (c Candle) GetVolume() float64 {
	return c.volume
}
