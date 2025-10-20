package gota

// Series - интерфейс для работы масивами данных
type Series interface {
	Len() int
	At(index int) Candle
	Slice(start, end int) Series
}
