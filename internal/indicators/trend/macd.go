package trend

import (
	"fmt"
	"time"

	"github.com/egor-erm/gota/pkg/gota"
	"github.com/egor-erm/gota/pkg/utils"
)

// MACD - Moving Average Convergence/Divergence
// https://gerchik.com/journal/foreks/macd-indicator/
type MACD struct {
	fastPeriod   int
	slowPeriod   int
	signalPeriod int
}

type MACDResult struct {
	MACDLine   []float64 // Основная линия MACD
	SignalLine []float64 // Сигнальная линия
	Histogram  []float64 // Гистограмма (MACD - Signal)
}

func NewMACD(fastPeriod, slowPeriod, signalPeriod int) *MACD {
	return &MACD{
		fastPeriod:   fastPeriod,
		slowPeriod:   slowPeriod,
		signalPeriod: signalPeriod,
	}
}

func (m MACD) Calculate(series gota.Series) *MACDResult {
	if series.Len() < m.slowPeriod {
		return nil
	}

	// Вычисляем EMA для быстрой и медленной линии
	fastEMA := NewEMA(m.fastPeriod).Calculate(series)
	slowEMA := NewEMA(m.slowPeriod).Calculate(series)
	fmt.Println(fastEMA, slowEMA)

	// Выравниваем длины (EMA начинаются с разных индексов)
	fastEMA, slowEMA = utils.AlignLengths(fastEMA, slowEMA)

	macdLine := make([]float64, 0)
	for i := 0; i < len(fastEMA); i++ {
		macdLine = append(macdLine, fastEMA[i]-slowEMA[i])
	}

	// Вычисляем сигнальную линию (EMA от MACD)
	macdSeries := createFloatSeries(macdLine)
	signalLine := NewEMA(m.signalPeriod).Calculate(macdSeries)

	// Выравниваем длины (EMA начинаются с разных индексов)
	macdLine, signalLine = utils.AlignLengths(macdLine, signalLine)

	// Вычисляем гистограмму
	histogram := make([]float64, len(signalLine))
	for i := 0; i < len(signalLine); i++ {
		histogram[i] = macdLine[i] - signalLine[i]
	}

	return &MACDResult{macdLine, signalLine, histogram}
}

// Вспомогательная функция для создания Series из []float64
func createFloatSeries(data []float64) gota.Series {
	baseTime := time.Now()
	candles := make([]gota.Candle, len(data))

	for i, value := range data {
		candles[i] = gota.NewCandle(
			baseTime.AddDate(0, 0, i),
			value, value, value, value, 1.0,
		)
	}

	return gota.CandleSeries(candles)
}
