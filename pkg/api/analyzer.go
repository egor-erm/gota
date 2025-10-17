package api

import (
	"github.com/egor-erm/gota/internal/indicators/momentum"
	"github.com/egor-erm/gota/internal/indicators/trend"
	"github.com/egor-erm/gota/pkg/gota"
)

// Analyzer - структура для анализа данных
type Analyzer struct {
	series gota.Series
}

func NewAnalyzer(series gota.Series) *Analyzer {
	return &Analyzer{
		series: series,
	}
}

func (a *Analyzer) SMA(period int) []float64 {
	sma := trend.NewSMA(period)

	return sma.Calculate(a.series)
}

func (a *Analyzer) EMA(period int) []float64 {
	ema := trend.NewEMA(period)

	return ema.Calculate(a.series)
}

func (a *Analyzer) WMA(period int) []float64 {
	wma := trend.NewWMA(period)

	return wma.Calculate(a.series)
}

func (a *Analyzer) RSI(period int) []float64 {
	rsi := momentum.NewRSI(period)

	return rsi.Calculate(a.series)
}
