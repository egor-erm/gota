package trend

import (
	"github.com/egor-erm/gota/pkg/gota"
)

// WMA - Weighted Moving Average
// https://www.binance.com/ru/academy/glossary/weighted-moving-average-wma
type WMA struct {
	period int
}

func NewWMA(period int) *WMA {
	return &WMA{period: period}
}

func (w WMA) Period() int {
	return w.period
}

func (w WMA) Calculate(series gota.Series) []float64 {
	if series.Len() < w.period {
		return nil
	}

	result := make([]float64, 0)

	for i := w.period - 1; i < series.Len(); i++ {
		sum := 0.0
		weightSum := 0.0

		for j := 0; j < w.period; j++ {
			weight := float64(w.period - j)
			sum += series.At(i-j).GetClosePrice() * weight
			weightSum += weight
		}

		result = append(result, sum/weightSum)
	}

	return result
}
