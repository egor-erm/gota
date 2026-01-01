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
	candles := createTestCandles()

	// Создаем стратегию
	strategy := NewSimpleStrategy()

	// Создаем бектестер
	backtester := api.NewBacktester(5000)

	// Выполняем бектест
	result := backtester.Backtest(strategy, candles)

	// Выводим результаты
	backtester.PrintResults(result)

	// Экспортируем результаты
	exporter := &api.CSVExporter{}

	// Экспорт сделок
	err := exporter.ExportTrades(result.Trades, "trades.csv")
	if err != nil {
		fmt.Printf("Ошибка экспорта сделок: %v\n", err)
	}

	// Экспорт кривой капитала
	err = exporter.ExportEquityCurve(result.EquityCurve, candles, "equity_curve.csv")
	if err != nil {
		fmt.Printf("Ошибка экспорта кривой капитала: %v\n", err)
	}
}

func createTestCandles() gota.CandleSeries {
	candles := make([]gota.Candle, 100)
	baseTime := time.Now()

	candles[0] = gota.NewCandle(
		baseTime.AddDate(0, 0, 0),
		3000, 3050, 2900, 2950, 1000.0,
	)

	for i := 1; i < 100; i++ {
		open := candles[i-1].GetClosePrice()
		close := open + (float64(i%10)-5)*10
		high := math.Max(open, close) + 50
		low := math.Min(open, close) - 50

		candles[i] = gota.NewCandle(
			baseTime.AddDate(0, 0, i),
			open, high, low, close, 1000.0,
		)
	}

	return candles
}

// SimpleStrategy - простая стратегия
type SimpleStrategy struct {
}

func NewSimpleStrategy() *SimpleStrategy {
	return &SimpleStrategy{}
}

func (s *SimpleStrategy) Name() string {
	return "Simple Strategy"
}

func (s *SimpleStrategy) Analyze(candles gota.CandleSeries) []api.TradeSignal {
	signals := make([]api.TradeSignal, 0)

	for i := 1; i < len(candles); i++ {
		candle := candles[i]
		price := candle.GetClosePrice()

		// Сигнал на покупку
		if i == 1 {
			fmt.Println("Вход:", price)
			signals = append(signals, api.TradeSignal{
				Time:     candle.GetStartTime(),
				Price:    price,
				IsEntry:  true,
				Type:     api.TradeTypeShort,
				Strength: 1.0,
			})
		}

		// Сигнал на продажу
		if i == 50 {
			fmt.Println("Выход:", price)
			signals = append(signals, api.TradeSignal{
				Time:     candle.GetStartTime(),
				Price:    price,
				IsEntry:  false,
				Type:     api.TradeTypeShort,
				Strength: 1.0,
			})
		}
	}

	return signals
}
