package volatility

import (
	"math"

	"github.com/egor-erm/gota"
	"github.com/egor-erm/gota/indicators/trend"
)

// Bollinger Bands
type BollingerBands struct {
	period int
	stdDev float64
}

type BollingerBandsResult struct {
	UpperBand  []float64
	MiddleBand []float64
	LowerBand  []float64
}

func NewBollingerBands(period int, stdDev float64) *BollingerBands {
	return &BollingerBands{
		period: period,
		stdDev: stdDev,
	}
}

func (bb BollingerBands) Calculate(series gota.Series) *BollingerBandsResult {
	if series.Len() < bb.period {
		return nil
	}

	upperBand := make([]float64, 0)
	middleBand := make([]float64, 0)
	lowerBand := make([]float64, 0)

	sma := trend.NewSMA(bb.period)
	smaValues := sma.Calculate(series)

	for i := bb.period - 1; i < series.Len(); i++ {
		sumSquares := 0.0
		smaValue := smaValues[i-bb.period+1]

		for j := 0; j < bb.period; j++ {
			diff := series.At(i-j).GetClosePrice() - smaValue
			sumSquares += diff * diff
		}

		stdDev := math.Sqrt(sumSquares / float64(bb.period))

		middleBand = append(middleBand, smaValue)
		upperBand = append(upperBand, smaValue+bb.stdDev*stdDev)
		lowerBand = append(lowerBand, smaValue-bb.stdDev*stdDev)
	}

	return &BollingerBandsResult{upperBand, middleBand, lowerBand}
}
