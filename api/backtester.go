package api

import (
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/egor-erm/gota"
)

// TradeType - тип сделки
type TradeType string

const (
	TradeTypeLong  TradeType = "LONG"
	TradeTypeShort TradeType = "SHORT"
)

// Trade - структура сделки
type Trade struct {
	EntryTime     time.Time
	ExitTime      time.Time
	EntryPrice    float64
	ExitPrice     float64
	Type          TradeType
	Profit        float64
	ProfitPercent float64
}

// TradeSignal - сигнал для входа/выхода
type TradeSignal struct {
	Time     time.Time
	Price    float64
	IsEntry  bool
	Type     TradeType
	Strength float64 // сила сигнала (0-1)
}

// Strategy - интерфейс стратегии
type Strategy interface {
	// Analyze анализирует свечи и возвращает сигналы
	Analyze(candles gota.CandleSeries) []TradeSignal
	// Name возвращает название стратегии
	Name() string
}

// BacktestResult - результаты бектеста
type BacktestResult struct {
	TotalTrades   int
	WinningTrades int
	LosingTrades  int
	TotalProfit   float64
	TotalReturn   float64
	MaxDrawdown   float64
	WinRate       float64
	AvgWin        float64
	AvgLoss       float64
	ProfitFactor  float64
	SharpeRatio   float64
	Trades        []Trade
	EquityCurve   []float64
	MaxEquity     float64
	MinEquity     float64
}

// Backtester - структура для бектестинга
type Backtester struct {
	initialCapital float64
	commission     float64
	positionSize   float64 // в процентах от капитала (0-1)
	slippage       float64 // проскальзывание в процентах
}

// NewBacktester создает новый бектестер
func NewBacktester(initialCapital float64) *Backtester {
	return &Backtester{
		initialCapital: initialCapital,
		commission:     0.001, // 0.1% по умолчанию
		positionSize:   1.0,   // 100% капитала
		slippage:       0.001, // 0.1% по умолчанию
	}
}

// SetCommission устанавливает комиссию
func (b *Backtester) SetCommission(commission float64) {
	b.commission = commission
}

// SetPositionSize устанавливает размер позиции
func (b *Backtester) SetPositionSize(size float64) error {
	if size <= 0 || size > 1 {
		return errors.New("размер позиции должен быть между 0 и 1")
	}
	b.positionSize = size
	return nil
}

// SetSlippage устанавливает проскальзывание
func (b *Backtester) SetSlippage(slippage float64) {
	b.slippage = slippage
}

// Backtest выполняет бектест стратегии
func (b *Backtester) Backtest(strategy Strategy, candles gota.CandleSeries) *BacktestResult {
	if candles.Len() == 0 {
		return nil
	}

	// Получаем сигналы от стратегии
	signals := strategy.Analyze(candles)

	// Инициализируем результат
	result := &BacktestResult{
		Trades:      make([]Trade, 0),
		EquityCurve: make([]float64, candles.Len()),
	}

	// Переменные для отслеживания состояния
	var currentTrade *Trade
	equity := b.initialCapital
	peakEquity := equity
	maxDrawdown := 0.0

	// Обрабатываем каждую свечу
	for i := 0; i < candles.Len(); i++ {
		candle := candles.At(i)
		price := candle.GetClosePrice()

		// Проверяем сигналы для этой свечи
		for _, signal := range signals {
			if !signal.Time.Equal(candle.GetStartTime()) {
				continue
			}

			if signal.IsEntry {
				// Сигнал на вход
				if currentTrade == nil {
					// Учитываем проскальзывание
					entryPrice := price * (1 + b.slippage)

					currentTrade = &Trade{
						EntryTime:  signal.Time,
						EntryPrice: entryPrice,
						Type:       signal.Type,
					}
				}
			} else {
				// Сигнал на выход
				if currentTrade != nil && currentTrade.Type == signal.Type {
					// Учитываем проскальзывание
					exitPrice := price * (1 - b.slippage)

					currentTrade.ExitTime = signal.Time
					currentTrade.ExitPrice = exitPrice

					// Рассчитываем прибыль
					var profit float64
					if currentTrade.Type == TradeTypeLong {
						profit = (exitPrice - currentTrade.EntryPrice) / currentTrade.EntryPrice
					} else {
						profit = (currentTrade.EntryPrice - exitPrice) / currentTrade.EntryPrice
					}

					// Учитываем комиссию (вход + выход)
					profit -= b.commission * 2

					// Рассчитываем абсолютную прибыль
					positionValue := equity * b.positionSize
					profitAbs := profit * positionValue

					currentTrade.Profit = profitAbs
					currentTrade.ProfitPercent = profit * 100

					// Обновляем капитал
					equity += profitAbs

					// Сохраняем сделку
					result.Trades = append(result.Trades, *currentTrade)
					currentTrade = nil
				}
			}
		}

		// Обновляем кривую капитала
		result.EquityCurve[i] = equity

		// Рассчитываем максимальную просадку
		if equity > peakEquity {
			peakEquity = equity
		}
		drawdown := (peakEquity - equity) / peakEquity
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
	}

	// Форсируем выход из открытой позиции в конце
	if currentTrade != nil {
		lastCandle := candles.At(candles.Len() - 1)
		currentTrade.ExitTime = lastCandle.GetStartTime()
		currentTrade.ExitPrice = lastCandle.GetClosePrice()

		var profit float64
		if currentTrade.Type == TradeTypeLong {
			profit = (currentTrade.ExitPrice - currentTrade.EntryPrice) / currentTrade.EntryPrice
		} else {
			profit = (currentTrade.EntryPrice - currentTrade.ExitPrice) / currentTrade.EntryPrice
		}

		profit -= b.commission * 2
		positionValue := equity * b.positionSize
		profitAbs := profit * positionValue

		currentTrade.Profit = profitAbs
		currentTrade.ProfitPercent = profit * 100

		equity += profitAbs
		result.Trades = append(result.Trades, *currentTrade)
	}

	// Рассчитываем статистику
	b.calculateStatistics(result, b.initialCapital)

	return result
}

// calculateStatistics рассчитывает статистику бектеста
func (b *Backtester) calculateStatistics(result *BacktestResult, initialCapital float64) {
	result.TotalTrades = len(result.Trades)
	result.TotalProfit = result.EquityCurve[len(result.EquityCurve)-1] - initialCapital
	result.TotalReturn = (result.TotalProfit / initialCapital) * 100

	var totalWins, totalLosses float64
	var wins, losses int

	for _, trade := range result.Trades {
		if trade.Profit > 0 {
			wins++
			totalWins += trade.Profit
		} else {
			losses++
			totalLosses += trade.Profit
		}
	}

	result.WinningTrades = wins
	result.LosingTrades = losses

	if result.TotalTrades > 0 {
		result.WinRate = float64(wins) / float64(result.TotalTrades) * 100
	}

	if wins > 0 {
		result.AvgWin = totalWins / float64(wins)
	}

	if losses > 0 {
		result.AvgLoss = totalLosses / float64(losses)
	}

	if totalLosses != 0 {
		result.ProfitFactor = math.Abs(totalWins / totalLosses)
	}

	// Находим максимальное и минимальное значение капитала
	result.MaxEquity = initialCapital
	result.MinEquity = initialCapital
	for _, equity := range result.EquityCurve {
		if equity > result.MaxEquity {
			result.MaxEquity = equity
		}
		if equity < result.MinEquity {
			result.MinEquity = equity
		}
	}

	result.MaxDrawdown = (result.MaxEquity - result.MinEquity) / result.MaxEquity * 100
}

// PrintResults выводит результаты бектеста
func (b *Backtester) PrintResults(result *BacktestResult) {
	fmt.Println("=== РЕЗУЛЬТАТЫ БЕКТЕСТА ===")
	fmt.Printf("Всего сделок: %d\n", result.TotalTrades)
	fmt.Printf("Прибыльных: %d (%.1f%%)\n", result.WinningTrades, result.WinRate)
	fmt.Printf("Убыточных: %d\n", result.LosingTrades)
	fmt.Printf("Общая прибыль: $%.2f\n", result.TotalProfit)
	fmt.Printf("Общая доходность: %.2f%%\n", result.TotalReturn)
	fmt.Printf("Макс. просадка: %.2f%%\n", result.MaxDrawdown)
	fmt.Printf("Средняя прибыль: $%.2f\n", result.AvgWin)
	fmt.Printf("Средний убыток: $%.2f\n", result.AvgLoss)
	fmt.Printf("Фактор прибыли: %.2f\n", result.ProfitFactor)
	fmt.Println("==========================")

	// Выводим детали по сделкам
	if len(result.Trades) > 0 {
		fmt.Println("\n=== ДЕТАЛИ СДЕЛОК ===")

		for i, trade := range result.Trades {
			// Определяем статус сделки
			status := "ПРИБЫЛЬ"
			if trade.Profit < 0 {
				status = "УБЫТОК"
			}

			// Форматируем время для читаемости
			entryTime := trade.EntryTime.Format("2006-01-02 15:04:05")
			exitTime := trade.ExitTime.Format("2006-01-02 15:04:05")

			fmt.Printf("\nСделка #%d [%s]\n", i+1, status)
			fmt.Printf("  Тип: %s\n", trade.Type)
			fmt.Printf("  Вход:  %s по $%.2f\n", entryTime, trade.EntryPrice)
			fmt.Printf("  Выход: %s по $%.2f\n", exitTime, trade.ExitPrice)
			fmt.Printf("  Длительность: %v\n", trade.ExitTime.Sub(trade.EntryTime))
			fmt.Printf("  Изменение цены: $%.2f (%.2f%%)\n",
				trade.ExitPrice-trade.EntryPrice,
				((trade.ExitPrice-trade.EntryPrice)/trade.EntryPrice)*100)
			fmt.Printf("  Прибыль/Убыток: $%.2f (%.2f%%)\n", trade.Profit, trade.ProfitPercent)
		}

		// Сводка по сделкам
		fmt.Println("\n=== СВОДКА ПО СДЕЛКАМ ===")

		var totalLongProfit, totalShortProfit float64
		var longTrades, shortTrades int
		var winningLongs, winningShorts int

		for _, trade := range result.Trades {
			if trade.Type == TradeTypeLong {
				longTrades++
				totalLongProfit += trade.Profit
				if trade.Profit > 0 {
					winningLongs++
				}
			} else {
				shortTrades++
				totalShortProfit += trade.Profit
				if trade.Profit > 0 {
					winningShorts++
				}
			}
		}

		fmt.Printf("LONG сделок: %d\n", longTrades)
		if longTrades > 0 {
			fmt.Printf("  Прибыльных LONG: %d (%.1f%%)\n",
				winningLongs, float64(winningLongs)/float64(longTrades)*100)
			fmt.Printf("  Общая прибыль по LONG: $%.2f\n", totalLongProfit)
			fmt.Printf("  Средняя прибыль по LONG: $%.2f\n", totalLongProfit/float64(longTrades))
		}

		fmt.Printf("SHORT сделок: %d\n", shortTrades)
		if shortTrades > 0 {
			fmt.Printf("  Прибыльных SHORT: %d (%.1f%%)\n",
				winningShorts, float64(winningShorts)/float64(shortTrades)*100)
			fmt.Printf("  Общая прибыль по SHORT: $%.2f\n", totalShortProfit)
			fmt.Printf("  Средняя прибыль по SHORT: $%.2f\n", totalShortProfit/float64(shortTrades))
		}
	}
}

// CSVExporter экспортирует результаты в CSV
type CSVExporter struct{}

// ExportTrades экспортирует сделки в CSV
func (e *CSVExporter) ExportTrades(trades []Trade, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Заголовки
	headers := []string{
		"EntryTime", "ExitTime", "Type", "EntryPrice", "ExitPrice",
		"Profit", "ProfitPercent",
	}
	writer.Write(headers)

	// Данные
	for _, trade := range trades {
		record := []string{
			trade.EntryTime.Format(time.RFC3339),
			trade.ExitTime.Format(time.RFC3339),
			string(trade.Type),
			strconv.FormatFloat(trade.EntryPrice, 'f', 2, 64),
			strconv.FormatFloat(trade.ExitPrice, 'f', 2, 64),
			strconv.FormatFloat(trade.Profit, 'f', 2, 64),
			strconv.FormatFloat(trade.ProfitPercent, 'f', 2, 64),
		}
		writer.Write(record)
	}

	return nil
}

// ExportEquityCurve экспортирует кривую капитала в CSV
func (e *CSVExporter) ExportEquityCurve(equityCurve []float64, candles gota.CandleSeries, filename string) error {
	if len(equityCurve) != candles.Len() {
		return errors.New("длина кривой капитала не совпадает с количеством свечей")
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Заголовки
	writer.Write([]string{"Time", "Equity"})

	// Данные
	for i := 0; i < candles.Len(); i++ {
		record := []string{
			candles.At(i).GetStartTime().Format(time.RFC3339),
			strconv.FormatFloat(equityCurve[i], 'f', 2, 64),
		}
		writer.Write(record)
	}

	return nil
}
