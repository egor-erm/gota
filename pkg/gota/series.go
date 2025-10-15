package gota

// Series - интерфейс для работы масивами данных
type Series interface {
	Len() int
	At(index int) Candle
	Slice(start, end int) Series
}

// CandleSeries - структура для работы массивами свечей
type CandleSeries []Candle

func (cs CandleSeries) Len() int {
	return len(cs)
}

func (cs CandleSeries) At(index int) Candle {
	return cs[index]
}

func (cs CandleSeries) Slice(start, end int) Series {
	if start < 0 || end > len(cs) || start > end {
		return CandleSeries{}
	}

	return cs[start:end]
}
