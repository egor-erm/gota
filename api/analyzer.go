package api

import (
	"github.com/egor-erm/gota"
	"github.com/egor-erm/gota/indicators/momentum"
	"github.com/egor-erm/gota/indicators/trend"
	"github.com/egor-erm/gota/indicators/volatility"
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

func (a *Analyzer) MACD(fast, slow, signal int) ([]float64, []float64, []float64) {
	macd := trend.NewMACD(fast, slow, signal)

	result := macd.Calculate(a.series)
	if result == nil {
		return nil, nil, nil
	}

	return result.MACDLine, result.SignalLine, result.Histogram
}

func (a *Analyzer) RSI(period int) []float64 {
	rsi := momentum.NewRSI(period)

	return rsi.Calculate(a.series)
}

func (a *Analyzer) ATR(period int) []float64 {
	atr := volatility.NewATR(period)

	return atr.Calculate(a.series)
}

func (a *Analyzer) BollingerBands(period int, stdDev float64) ([]float64, []float64, []float64) {
	bb := volatility.NewBollingerBands(period, stdDev)

	result := bb.Calculate(a.series)
	if result == nil {
		return nil, nil, nil
	}

	return result.UpperBand, result.MiddleBand, result.LowerBand
}
