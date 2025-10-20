package gota

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
