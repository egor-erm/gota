package utils

func AlignLengths(arr1, arr2 []float64) ([]float64, []float64) {
	if len(arr1) == len(arr2) {
		return arr1, arr2
	}

	if len(arr1) > len(arr2) {
		startIndex := len(arr1) - len(arr2)

		return arr1[startIndex:], arr2
	} else {
		startIndex := len(arr2) - len(arr1)

		return arr1, arr2[startIndex:]
	}
}
