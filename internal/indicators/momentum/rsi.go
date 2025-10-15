package momentum

import (
	"math"

	"github.com/egor-erm/gota/pkg/gota"
)

// RSI - Relative Strength Index
type RSI struct {
	period int
}

func NewRSI(period int) *RSI {
	return &RSI{period: period}
}

func (r RSI) Period() int {
	return r.period
}

func (r RSI) Calculate(series gota.Series) []float64 {
	if series.Len() < r.period+1 {
		return nil
	}

	// Инициализация gains и losses
	gains := make([]float64, series.Len())
	losses := make([]float64, series.Len())

	// Вычисляем изменения и разделяем на gains/losses
	for i := 1; i < series.Len(); i++ {
		change := series.At(i).GetClosePrice() - series.At(i-1).GetClosePrice()
		if change > 0 {
			gains[i] = change
			losses[i] = 0
		} else {
			gains[i] = 0
			losses[i] = math.Abs(change)
		}
	}

	// Средние значения gain/loss
	avgGain := 0.0
	avgLoss := 0.0

	for i := 1; i <= r.period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}

	avgGain /= float64(r.period)
	avgLoss /= float64(r.period)

	result := make([]float64, 0)

	// Расчёт для первой свечи
	result = append(result, calculateRSIValue(avgGain, avgLoss))

	// Расчёт для остальных свечей
	for i := r.period + 1; i < series.Len(); i++ {
		// Обновляем средние значения
		avgGain = (avgGain*float64(r.period-1) + gains[i]) / float64(r.period)
		avgLoss = (avgLoss*float64(r.period-1) + losses[i]) / float64(r.period)

		result = append(result, calculateRSIValue(avgGain, avgLoss))
	}

	return result
}

// Sычисляет значение RSI из средних значений gain/loss
func calculateRSIValue(avgGain, avgLoss float64) float64 {
	if avgLoss == 0 {
		return 100.0
	}

	rs := avgGain / avgLoss
	return 100.0 - (100.0 / (1.0 + rs))
}
