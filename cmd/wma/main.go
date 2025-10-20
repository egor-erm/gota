package main

import (
	"fmt"
	"time"

	"github.com/egor-erm/gota"
	"github.com/egor-erm/gota/api"
)

// Пример с сайта https://www.binance.com/ru/academy/glossary/weighted-moving-average-wma
func main() {
	candles := createCandles()

	for i, candle := range candles {
		fmt.Printf("Day %d: %v $\n", i+1, candle.GetClosePrice())
	}

	analyser := api.NewAnalyzer(candles)
	wma := analyser.WMA(5)

	fmt.Println("WMA", wma)
}

func createCandles() gota.CandleSeries {
	baseTime := time.Now()
	closePrices := []float64{10, 11, 12, 13, 14}

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
