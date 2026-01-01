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

func (cs CandleSeries) Last(n int) CandleSeries {
	if n <= 0 {
		return CandleSeries{}
	}

	if n >= len(cs) {
		return cs
	}

	return cs[len(cs)-n:]
}

func (cs CandleSeries) LastMin(nums ...int) CandleSeries {
	if len(nums) == 0 {
		return CandleSeries{}
	}

	minN := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] < minN {
			minN = nums[i]
		}
	}

	return cs.Last(minN)
}
