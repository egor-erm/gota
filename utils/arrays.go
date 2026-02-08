package utils

func AlignLengths(arr1, arr2 []float64) ([]float64, []float64) {
	aligned := AlignLengthsMulti(arr1, arr2)
	return aligned[0], aligned[1]
}

func AlignLengthsMulti(arrays ...[]float64) [][]float64 {
	if len(arrays) == 0 {
		return nil
	}

	// Находим МИНИМАЛЬНУЮ длину
	minLength := len(arrays[0])
	for _, arr := range arrays[1:] {
		if len(arr) < minLength {
			minLength = len(arr)
		}
	}

	aligned := make([][]float64, len(arrays))
	for i, arr := range arrays {
		if len(arr) == minLength {
			aligned[i] = arr
		} else {
			// Берем последние minLength элементов
			startIndex := len(arr) - minLength
			aligned[i] = arr[startIndex:]
		}
	}

	return aligned
}
