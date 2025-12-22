package main

import (
	"fmt"
	"math"
	"time"

	"github.com/egor-erm/gota"
	"github.com/egor-erm/gota/api"
)

func main() {
	// Создаем тестовые данные
	candles := make([]gota.Candle, 100)
	baseTime := time.Now()

	candles[0] = gota.NewCandle(
		baseTime.AddDate(0, 0, 0),
		3, 4, 1, 2, 1000.0,
	)

	for i := 1; i < 100; i++ {
		open := candles[i-1].GetClosePrice()
		close := open + (float64(i%10)-5)*0.3
		high := math.Max(open, close) + 0.5
		low := math.Min(open, close) - 0.5

		candles[i] = gota.NewCandle(
			baseTime.AddDate(0, 0, i),
			open, high, low, close, 1000.0,
		)
	}

	series := gota.CandleSeries(candles)

	// Создаем визуализатор
	visualizer := api.NewVisualizer(series, 1200, 800)

	// Добавляем индикаторы
	visualizer.AddStochRSI(14, 14, 3, 3)

	// Рендерим и сохраняем
	err := visualizer.RenderToFile("chart.png")
	if err != nil {
		panic(err)
	}

	fmt.Println("Chart saved to chart.png")
}
