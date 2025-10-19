package main

import (
	"fmt"
	"time"

	"github.com/egor-erm/gota/pkg/api"
	"github.com/egor-erm/gota/pkg/gota"
)

// Реализация RSI совпадает с реализацией с сайта https://habr.com/ru/companies/otus/articles/823500/
func main() {
	candles := createCandles()

	for i, candle := range candles {
		fmt.Printf("Day %d: %v $\n", i+1, candle.GetClosePrice())
	}

	analyser := api.NewAnalyzer(candles)
	rsi := analyser.RSI(14)

	fmt.Println("RSI", rsi)
}

func createCandles() gota.CandleSeries {
	baseTime := time.Now()
	closePrices := []float64{49.85, 50.00, 52.25, 51.50, 49.75, 48.90, 51.10, 52.40, 54.20, 55.80, 56.50, 50.25, 49.35, 49.90, 52.10}

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
