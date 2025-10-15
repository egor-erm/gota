package trend

import (
	"github.com/egor-erm/gota/pkg/gota"
)

// SMA - Simple Moving Average
type SMA struct {
	period int
}

func NewSMA(period int) *SMA {
	return &SMA{period: period}
}

func (s SMA) Period() int {
	return s.period
}

func (s SMA) Calculate(series gota.Series) []float64 {
	if series.Len() < s.period {
		return nil
	}

	result := make([]float64, 0)

	// Начинаем с первой свечи, для которой можем рассчитать SMA (до неё должно быть period-1 свечей)
	for i := s.period - 1; i < series.Len(); i++ {
		sum := 0.0
		for j := 0; j < s.period; j++ {
			sum += series.At(i - j).GetClosePrice()
		}

		result = append(result, sum/float64(s.period))
	}

	return result
}
