package volatility

import (
	"math"

	"github.com/egor-erm/gota"
)

// ATR - Average True Range
// https://ru.tradingview.com/chart/BTCUSD/54j9Q1zR-polnoe-rukovodstvo-po-indikatoru-atr/
type ATR struct {
	period int
}

func NewATR(period int) *ATR {
	return &ATR{period: period}
}

func (a ATR) Calculate(series gota.Series) []float64 {
	if series.Len() < a.period {
		return nil
	}

	trueRanges := make([]float64, series.Len())

	// Вычисляем True Range для каждой свечи
	for i := 1; i < series.Len(); i++ {
		current := series.At(i)
		previous := series.At(i - 1)

		tr1 := current.GetHighPrice() - current.GetLowPrice()
		tr2 := math.Abs(current.GetHighPrice() - previous.GetClosePrice())
		tr3 := math.Abs(current.GetLowPrice() - previous.GetClosePrice())

		trueRanges[i] = math.Max(tr1, math.Max(tr2, tr3))
	}

	// Вычисляем ATR как SMA от True Range
	result := make([]float64, 0)

	// Первое значение ATR - среднее за период
	firstATR := 0.0
	for i := 0; i < a.period; i++ {
		firstATR += trueRanges[i]
	}

	firstATR /= float64(a.period)
	result = append(result, firstATR)

	// Остальные значения используя формулу Вилдера
	for i := a.period + 1; i < len(trueRanges); i++ {
		atr := (result[len(result)-1]*float64(a.period-1) + trueRanges[i]) / float64(a.period)
		result = append(result, atr)
	}

	return result
}
