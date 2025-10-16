package main

import (
	"fmt"
	"time"

	"github.com/egor-erm/gota/pkg/api"
	"github.com/egor-erm/gota/pkg/gota"
)

// Пример с сайта https://www.binance.com/ru/academy/glossary/exponential-moving-average-ema
func main() {
	candles := createCandles()

	for i, candle := range candles {
		fmt.Printf("Day %d: %v $\n", i+1, candle.GetClosePrice())
	}

	analyser := api.NewAnalyzer(candles)
	sma := analyser.EMA(10)

	fmt.Println("EMA", sma)
}

func createCandles() gota.CandleSeries {
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	closePrices := []float64{50, 57, 58, 53, 55, 49, 56, 54, 63, 64, 60}

	candles := make([]gota.Candle, len(closePrices))
	for i := 0; i < len(closePrices); i++ {
		candles[i] = gota.NewCandle(
			baseTime.AddDate(0, 0, i),
			closePrices[i]-0.5,
			closePrices[i]+1.0,
			closePrices[i]-1.0,
			closePrices[i],
			1000.0,
		)
	}

	return candles
}
