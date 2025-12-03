package momentum

import (
	"github.com/egor-erm/gota"
)

// StochRSI - Stochastic RSI
// https://www.investopedia.com/terms/s/stochrsi.asp
type StochRSI struct {
	rsiPeriod   int
	stochPeriod int
	smoothK     int
	smoothD     int
}

type StochRSIResult struct {
	K []float64 // Быстрая %K линия
	D []float64 // Медленная %D линия
}

func NewStochRSI(rsiPeriod, stochPeriod, smoothK, smoothD int) *StochRSI {
	return &StochRSI{
		rsiPeriod:   rsiPeriod,
		stochPeriod: stochPeriod,
		smoothK:     smoothK,
		smoothD:     smoothD,
	}
}

func (s StochRSI) Calculate(series gota.Series) *StochRSIResult {
	// Сначала вычисляем RSI
	rsi := NewRSI(s.rsiPeriod)
	rsiValues := rsi.Calculate(series)

	if len(rsiValues) < s.stochPeriod {
		return nil
	}

	// Вычисляем Stochastic для RSI значений
	stochasticK := make([]float64, 0)

	// Для каждого окна вычисляем %K
	for i := s.stochPeriod - 1; i < len(rsiValues); i++ {
		// Находим min и max RSI в окне
		minRSI := rsiValues[i]
		maxRSI := rsiValues[i]

		for j := 0; j < s.stochPeriod; j++ {
			value := rsiValues[i-j]
			if value < minRSI {
				minRSI = value
			}
			if value > maxRSI {
				maxRSI = value
			}
		}

		// Вычисляем %K для StochRSI
		if maxRSI-minRSI == 0 {
			stochasticK = append(stochasticK, 100.0)
		} else {
			k := 100 * (rsiValues[i] - minRSI) / (maxRSI - minRSI)
			stochasticK = append(stochasticK, k)
		}
	}

	// Сглаживаем %K линию (если нужно)
	kLine := stochasticK
	if s.smoothK > 1 {
		kLine = smoothValues(stochasticK, s.smoothK)
	}

	// Вычисляем %D линию (сглаженную версию %K)
	dLine := kLine
	if s.smoothD > 1 {
		dLine = smoothValues(kLine, s.smoothD)
	}

	// Выравниваем длины
	if len(kLine) > len(dLine) {
		kLine = kLine[len(kLine)-len(dLine):]
	}

	return &StochRSIResult{
		K: kLine,
		D: dLine,
	}
}

// smoothValues сглаживает значения используя SMA
func smoothValues(values []float64, period int) []float64 {
	if len(values) < period {
		return values
	}

	result := make([]float64, len(values)-period+1)

	for i := period - 1; i < len(values); i++ {
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += values[i-j]
		}
		result[i-period+1] = sum / float64(period)
	}

	return result
}
