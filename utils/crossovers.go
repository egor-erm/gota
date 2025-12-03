package utils

import (
	"math"
)

// CrossoverType тип пересечения
type CrossoverType int

const (
	// NoCross - нет пересечения
	NoCross CrossoverType = iota
	// BullishCross - бычий кросс (нижняя линия пересекает верхнюю снизу вверх)
	BullishCross
	// BearishCross - медвежий кросс (верхняя линия пересекает нижнюю сверху вниз)
	BearishCross
)

// CrossoverEvent событие пересечения
type CrossoverEvent struct {
	Type       CrossoverType
	Index      int     // Индекс свечи где произошло пересечение
	Value      float64 // Значение на момент пересечения
	Line1Value float64 // Значение первой линии
	Line2Value float64 // Значение второй линии
}

// FindCrossovers находит все пересечения между двумя линиями
func FindCrossovers(line1, line2 []float64) []CrossoverEvent {
	if len(line1) != len(line2) || len(line1) < 2 {
		return nil
	}

	var events []CrossoverEvent
	for i := 1; i < len(line1); i++ {
		prevDiff := line1[i-1] - line2[i-1]
		currDiff := line1[i] - line2[i]

		// Проверяем пересечение
		if (prevDiff < 0 && currDiff > 0) ||
			(math.Abs(prevDiff) < 1e-10 && currDiff > 0) {
			// Бычий кросс: line1 пересекает line2 снизу вверх
			events = append(events, CrossoverEvent{
				Type:       BullishCross,
				Index:      i,
				Value:      (line1[i] + line2[i]) / 2,
				Line1Value: line1[i],
				Line2Value: line2[i],
			})
		} else if (prevDiff > 0 && currDiff < 0) ||
			(math.Abs(prevDiff) < 1e-10 && currDiff < 0) {
			// Медвежий кросс: line1 пересекает line2 сверху вниз
			events = append(events, CrossoverEvent{
				Type:       BearishCross,
				Index:      i,
				Value:      (line1[i] + line2[i]) / 2,
				Line1Value: line1[i],
				Line2Value: line2[i],
			})
		}
	}

	return events
}

// FindGoldenCross находит "золотой крест" (быстрая MA пересекает медленную MA снизу вверх)
func FindGoldenCross(fastMA, slowMA []float64) []CrossoverEvent {
	return findMACross(fastMA, slowMA, BullishCross)
}

// FindDeathCross находит "крест смерти" (быстрая MA пересекает медленную MA сверху вниз)
func FindDeathCross(fastMA, slowMA []float64) []CrossoverEvent {
	return findMACross(fastMA, slowMA, BearishCross)
}

// findMACross находит пересечения скользящих средних
func findMACross(fastMA, slowMA []float64, crossType CrossoverType) []CrossoverEvent {
	crossovers := FindCrossovers(fastMA, slowMA)
	var result []CrossoverEvent

	for _, crossover := range crossovers {
		if crossover.Type == crossType {
			result = append(result, crossover)
		}
	}

	return result
}
