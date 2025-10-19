
# Gota - Библиотека технического анализа на Go

Простая и эффективная библиотека для расчета технических индикаторов на финансовых данных.

## 📊 Поддерживаемые индикаторы

### Трендовые индикаторы
- **SMA (Simple Moving Average)** - Простая скользящая средняя
- **EMA (Exponential Moving Average)** - Экспоненциальная скользящая средняя
- **WMA (Weighted Moving Average)** - Взвешенная скользящая средняя
- **MACD (Moving Average Convergence/Divergence)** - Cхождение/Расхождение скользящих средних
 

### Индикаторы момента
- **RSI (Relative Strength Index)** - Индекс относительной силы


## 🚀 Быстрый старт

### Установка
```go
import "github.com/egor-erm/gota/pkg/api"
```

### Использование
```go
package main

import (
    "fmt"
    "time"
    "github.com/egor-erm/gota/pkg/api"
    "github.com/egor-erm/gota/pkg/gota"
)

func main() {
    // Создание тестовых данных
    candles := createCandles()
    
    // Инициализация анализатора
    analyzer := api.NewAnalyzer(candles)
    
    // Расчет индикаторов
    sma := analyzer.SMA(10)    // SMA с периодом 10
    ema := analyzer.EMA(10)    // EMA с периодом 10
    rsi := analyzer.RSI(14)    // RSI с периодом 14
    
    fmt.Println("SMA:", sma)
    fmt.Println("EMA:", ema)
    fmt.Println("RSI:", rsi)
}

func createCandles() gota.CandleSeries {
    baseTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	closePrices := []float64{50.00, 52.25, 51.50, 49.75, 48.90, 51.10, 52.40,      54.20, 55.80, 56.50}
    
    candles := make([]gota.Candle, len(closePrices))
    for i := 0; i < len(closePrices); i++ {
        candles[i] = gota.NewCandle(
            baseTime.AddDate(0, 0, i),
            closePrices[i]-0.5,  // open
            closePrices[i]+1.0,  // high
            closePrices[i]-1.0,  // low
            closePrices[i],      // close
            1000.0,             // volume
        )
    }
    
    return candles
}
```
- ### Больше примеров вы всегда можете посмотреть в папке /cmd


## 🔮 Планы по развитию
- Добавление большего количества индикаторов
- Визуализация результатов на графиках


## 📄 Лицензия
MIT License


## ✨ Можем обсудить общие планы по развитию библиотеки - https://t.me/erm_egor
