package trend

import (
	"github.com/egor-erm/gota"
)

// EMA - Exponential Moving Average
// https://www.binance.com/ru/academy/glossary/exponential-moving-average-ema
type EMA struct {
	period int
}

func NewEMA(period int) *EMA {
	return &EMA{period: period}
}

func (e EMA) Period() int {
	return e.period
}

func (e EMA) Calculate(series gota.Series) []float64 {
	if series.Len() < e.period {
		return nil
	}

	result := make([]float64, 0)

	// Вычисляем множитель для EMA
	multiplier := 2.0 / (float64(e.period) + 1.0)

	// Первое значение EMA - это SMA за тот же период
	firstSMA := 0.0
	for i := 0; i < e.period; i++ {
		firstSMA += series.At(i).GetClosePrice()
	}
	firstSMA /= float64(e.period)

	// Добавляем первое значение EMA
	result = append(result, firstSMA)

	// Вычисляем остальные значения EMA
	for i := e.period; i < series.Len(); i++ {
		prevEMA := result[len(result)-1]

		// Формула EMA: (Close - PrevEMA) * multiplier + PrevEMA
		currentEMA := (series.At(i).GetClosePrice()-prevEMA)*multiplier + prevEMA
		result = append(result, currentEMA)
	}

	return result
}
