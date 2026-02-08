package trend

import (
	"math"

	"github.com/egor-erm/gota"
)

// ADX - Average Directional Index
type ADX struct {
	period int
}

type ADXResult struct {
	ADXValues []float64
	PlusDI    []float64
	MinusDI   []float64
}

func NewADX(period int) *ADX {
	return &ADX{period: period}
}

func (a ADX) Period() int {
	return a.period
}

// Calculate вычисляет ADX, +DI, -DI
func (a *ADX) Calculate(series gota.Series) *ADXResult {
	n := series.Len()
	if n < a.period*2 {
		return nil
	}

	// 1. Вычисляем TR, +DM, -DM для каждого бара (начиная со второго = индекс 1)
	tr := make([]float64, n)
	plusDM := make([]float64, n)
	minusDM := make([]float64, n)

	for i := 1; i < n; i++ {
		curr := series.At(i)
		prev := series.At(i - 1)

		high := curr.GetHighPrice()
		low := curr.GetLowPrice()
		prevClose := prev.GetClosePrice()
		prevHigh := prev.GetHighPrice()
		prevLow := prev.GetLowPrice()

		// True Range
		tr[i] = math.Max(high-low, math.Max(
			math.Abs(high-prevClose),
			math.Abs(low-prevClose),
		))

		// Directional Movement
		upMove := high - prevHigh
		downMove := prevLow - low

		if upMove > downMove && upMove > 0 {
			plusDM[i] = upMove
		}

		if downMove > upMove && downMove > 0 {
			minusDM[i] = downMove
		}
	}

	// 2. Сглаживание по Уайлдеру (Wilder's smoothing = RMА с α = 1/period)
	// Первое значение — простая сумма первых period ненулевых изменений (с индекса 1 по period)
	smoothedTR := make([]float64, n)
	smoothedPlusDM := make([]float64, n)
	smoothedMinusDM := make([]float64, n)

	var sumTR, sumPlus, sumMinus float64
	for i := 1; i <= a.period; i++ {
		sumTR += tr[i]
		sumPlus += plusDM[i]
		sumMinus += minusDM[i]
	}
	smoothedTR[a.period] = sumTR
	smoothedPlusDM[a.period] = sumPlus
	smoothedMinusDM[a.period] = sumMinus

	// Дальше — сглаживание: prev - prev/period + current
	for i := a.period + 1; i < n; i++ {
		smoothedTR[i] = smoothedTR[i-1] - smoothedTR[i-1]/float64(a.period) + tr[i]
		smoothedPlusDM[i] = smoothedPlusDM[i-1] - smoothedPlusDM[i-1]/float64(a.period) + plusDM[i]
		smoothedMinusDM[i] = smoothedMinusDM[i-1] - smoothedMinusDM[i-1]/float64(a.period) + minusDM[i]
	}

	// 3. +DI и -DI (начиная с индекса period)
	plusDI := make([]float64, n)
	minusDI := make([]float64, n)
	dx := make([]float64, n)

	for i := a.period; i < n; i++ {
		if smoothedTR[i] > 0 {
			plusDI[i] = 100 * smoothedPlusDM[i] / smoothedTR[i]
			minusDI[i] = 100 * smoothedMinusDM[i] / smoothedTR[i]

			diff := math.Abs(plusDI[i] - minusDI[i])
			sumDI := plusDI[i] + minusDI[i]
			if sumDI > 0 {
				dx[i] = 100 * diff / sumDI
			}
		}
	}

	// 4. ADX — сглаживание DX (первый ADX — простое среднее первых period значений DX)
	startADX := a.period * 2

	adx := make([]float64, 0, n-startADX)
	validPlusDI := make([]float64, 0, n-startADX)
	validMinusDI := make([]float64, 0, n-startADX)

	var sumDX float64
	for i := a.period + 1; i <= a.period*2; i++ {
		sumDX += dx[i]
	}

	firstADX := sumDX / float64(a.period)

	adx = append(adx, firstADX)
	validPlusDI = append(validPlusDI, plusDI[startADX])
	validMinusDI = append(validMinusDI, minusDI[startADX])

	// Последующие значения ADX — Wilder's smoothing
	for i := startADX + 1; i < n; i++ {
		nextADX := (adx[len(adx)-1]*float64(a.period-1) + dx[i]) / float64(a.period)
		adx = append(adx, nextADX)
		validPlusDI = append(validPlusDI, plusDI[i])
		validMinusDI = append(validMinusDI, minusDI[i])
	}

	return &ADXResult{
		ADXValues: adx,
		PlusDI:    validPlusDI,
		MinusDI:   validMinusDI,
	}
}
