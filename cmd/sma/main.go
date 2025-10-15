package main

import (
	"fmt"
	"time"

	"github.com/egor-erm/gota/pkg/api"
	"github.com/egor-erm/gota/pkg/gota"
)

// Пример с сайта https://learn.bybit.com/ru/indicators/what-is-simple-moving-average-sma
func main() {
	candles := createCandles()

	for i, candle := range candles {
		fmt.Printf("Day %d: %v $\n", i+1, candle.GetClosePrice())
	}

	analyser := api.NewAnalyzer(candles)
	sma := analyser.SMA(10)

	fmt.Println("SMA", sma)
}

func createCandles() gota.CandleSeries {
	baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	closePrices := []float64{50.00, 52.25, 51.50, 49.75, 48.90, 51.10, 52.40, 54.20, 55.80, 56.50}

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
