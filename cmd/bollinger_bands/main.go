package main

import (
	"fmt"
	"time"

	"github.com/egor-erm/gota"
	"github.com/egor-erm/gota/api"
)

// Пример с сайта
func main() {
	candles := createCandles()

	for i, candle := range candles {
		fmt.Printf("Day %d: %v $\n", i+1, candle.GetClosePrice())
	}

	analyser := api.NewAnalyzer(candles)
	upperBand, middleBand, lowerBand := analyser.BollingerBands(10, 2)

	fmt.Println("Bollinger Bands")
	fmt.Println("Upper Bands:", upperBand)
	fmt.Println("Middle Bands:", middleBand)
	fmt.Println("Lower Bands:", lowerBand)
}

func createCandles() gota.CandleSeries {
	baseTime := time.Now()
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
