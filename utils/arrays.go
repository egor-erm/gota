package utils

func AlignLengths(arr1, arr2 []float64) ([]float64, []float64) {
	aligned := AlignLengthsMulti(arr1, arr2)
	return aligned[0], aligned[1]
}

func AlignLengthsMulti(arrays ...[]float64) [][]float64 {
	if len(arrays) == 0 {
		return nil
	}

	maxLength := 0
	for _, arr := range arrays {
		if len(arr) > maxLength {
			maxLength = len(arr)
		}
	}

	aligned := make([][]float64, len(arrays))
	for i, arr := range arrays {
		if len(arr) == maxLength {
			aligned[i] = arr
			continue
		}

		startIndex := maxLength - len(arr)
		if startIndex < 0 {
			startIndex = 0
		}

		if startIndex < len(arr) {
			aligned[i] = arr[startIndex:]
		} else {
			aligned[i] = make([]float64, 0)
		}
	}

	return aligned
}
