package api

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"

	"github.com/egor-erm/gota"
	"github.com/fogleman/gg"
)

// Visualizer - структура для визуализации данных и индикаторов
type Visualizer struct {
	series          gota.Series
	indicators      []IndicatorConfig
	candles         []Candle
	width           int
	height          int
	margin          int
	candleWidth     int
	topHeight       int // высота для свечного графика
	indicatorHeight int // высота для каждого индикатора
	colors          Colors
}

// IndicatorConfig - конфигурация индикатора для отображения
type IndicatorConfig struct {
	Name    string
	Type    IndicatorType
	Data    [][]float64 // данные индикатора (может быть несколько линий)
	Labels  []string    // названия линий
	Color   color.Color
	SubType string // подтип (например, для MACD)
	Overlay bool   // отображать ли индикатор поверх свечей
}

type IndicatorType string

const (
	IndicatorSMA      IndicatorType = "SMA"
	IndicatorEMA      IndicatorType = "EMA"
	IndicatorWMA      IndicatorType = "WMA"
	IndicatorMACD     IndicatorType = "MACD"
	IndicatorRSI      IndicatorType = "RSI"
	IndicatorStochRSI IndicatorType = "StochRSI"
	IndicatorATR      IndicatorType = "ATR"
	IndicatorBB       IndicatorType = "BollingerBands"
)

// Candle - структура для отрисовки свечи
type Candle struct {
	Open, High, Low, Close float64
	X                      int
	Width                  int
}

// Colors - цвета для визуализации
type Colors struct {
	Background     color.Color
	Grid           color.Color
	Text           color.Color
	Bullish        color.Color
	Bearish        color.Color
	Volume         color.Color
	IndicatorLines []color.Color
}

// NewVisualizer создает новый визуализатор
func NewVisualizer(series gota.Series, width, height int) *Visualizer {
	return &Visualizer{
		series:          series,
		width:           width,
		height:          height,
		margin:          50,
		candleWidth:     10,
		topHeight:       height / 2, // половина высоты под свечи
		indicatorHeight: 150,        // высота каждого индикатора
		colors: Colors{
			Background: color.RGBA{20, 20, 30, 255},
			Grid:       color.RGBA{60, 60, 80, 255},
			Text:       color.RGBA{200, 200, 220, 255},
			Bullish:    color.RGBA{0, 200, 100, 255},
			Bearish:    color.RGBA{255, 80, 80, 255},
			Volume:     color.RGBA{100, 100, 255, 200},
			IndicatorLines: []color.Color{
				color.RGBA{255, 100, 100, 255}, // красный
				color.RGBA{100, 255, 100, 255}, // зеленый
				color.RGBA{100, 100, 255, 255}, // синий
				color.RGBA{255, 255, 100, 255}, // желтый
				color.RGBA{255, 100, 255, 255}, // фиолетовый
				color.RGBA{100, 255, 255, 255}, // голубой
			},
		},
	}
}

// AddIndicator добавляет индикатор для отображения
func (v *Visualizer) AddIndicator(config IndicatorConfig) {
	v.indicators = append(v.indicators, config)
}

// AddSMA добавляет SMA индикатор
func (v *Visualizer) AddSMA(period int, color color.Color) {
	sma := NewAnalyzer(v.series).SMA(period)
	if len(sma) > 0 {
		v.AddIndicator(IndicatorConfig{
			Name:    fmt.Sprintf("SMA(%d)", period),
			Type:    IndicatorSMA,
			Data:    [][]float64{sma},
			Labels:  []string{"SMA"},
			Color:   color,
			Overlay: true,
		})
	}
}

// AddEMA добавляет EMA индикатор
func (v *Visualizer) AddEMA(period int, color color.Color) {
	ema := NewAnalyzer(v.series).EMA(period)
	if len(ema) > 0 {
		v.AddIndicator(IndicatorConfig{
			Name:    fmt.Sprintf("EMA(%d)", period),
			Type:    IndicatorEMA,
			Data:    [][]float64{ema},
			Labels:  []string{"EMA"},
			Color:   color,
			Overlay: true,
		})
	}
}

// AddRSI добавляет RSI индикатор
func (v *Visualizer) AddRSI(period int) {
	rsi := NewAnalyzer(v.series).RSI(period)
	if len(rsi) > 0 {
		v.AddIndicator(IndicatorConfig{
			Name:    fmt.Sprintf("RSI(%d)", period),
			Type:    IndicatorRSI,
			Data:    [][]float64{rsi},
			Labels:  []string{"RSI"},
			Color:   v.colors.IndicatorLines[0],
			Overlay: false,
		})
	}
}

// AddMACD добавляет MACD индикатор
func (v *Visualizer) AddMACD(fast, slow, signal int) {
	macdLine, signalLine, histogram := NewAnalyzer(v.series).MACD(fast, slow, signal)
	if macdLine != nil {
		v.AddIndicator(IndicatorConfig{
			Name:    fmt.Sprintf("MACD(%d,%d,%d)", fast, slow, signal),
			Type:    IndicatorMACD,
			Data:    [][]float64{macdLine, signalLine, histogram},
			Labels:  []string{"MACD", "Signal", "Histogram"},
			Color:   v.colors.IndicatorLines[0],
			SubType: "line",
			Overlay: false,
		})
	}
}

// AddBollingerBands добавляет Bollinger Bands
func (v *Visualizer) AddBollingerBands(period int, stdDev float64) {
	upper, middle, lower := NewAnalyzer(v.series).BollingerBands(period, stdDev)
	if upper != nil {
		v.AddIndicator(IndicatorConfig{
			Name:    fmt.Sprintf("BB(%d,%.1f)", period, stdDev),
			Type:    IndicatorBB,
			Data:    [][]float64{upper, middle, lower},
			Labels:  []string{"Upper", "Middle", "Lower"},
			Color:   v.colors.IndicatorLines[2],
			Overlay: true,
		})
	}
}

// AddRSI добавляет StochRSI индикатор
func (v *Visualizer) AddStochRSI(rsiPeriod, stochPeriod, smoothK, smoothD int) {
	rsi1, rsi2 := NewAnalyzer(v.series).StochRSI(rsiPeriod, stochPeriod, smoothK, smoothD)
	if len(rsi1) > 0 {
		v.AddIndicator(IndicatorConfig{
			Name:    fmt.Sprintf("StochRSI(%d %d %d %d)", rsiPeriod, stochPeriod, smoothK, smoothD),
			Type:    IndicatorStochRSI,
			Data:    [][]float64{rsi1, rsi2},
			Labels:  []string{"StochRSI"},
			Color:   v.colors.IndicatorLines[2],
			Overlay: false,
		})
	}
}

// Render визуализирует график и возвращает изображение
func (v *Visualizer) Render() (image.Image, error) {
	// Создаем контекст для рисования
	dc := gg.NewContext(v.width, v.height)

	// Очищаем фон
	dc.SetColor(v.colors.Background)
	dc.Clear()

	// Подготавливаем свечи для отрисовки
	v.prepareCandles()

	// Рисуем сетку и оси
	v.drawGrid(dc)

	// Рисуем свечной график
	v.drawCandles(dc)

	// Рисуем оверлейные индикаторы (на том же графике, что и свечи)
	v.drawOverlayIndicators(dc)

	// Рассчитываем высоту для каждого индикатора
	indicatorCount := 0
	for _, ind := range v.indicators {
		if !ind.Overlay {
			indicatorCount++
		}
	}

	if indicatorCount > 0 {
		indicatorHeight := (v.height - v.topHeight - v.margin*2) / indicatorCount

		// Рисуем отдельные индикаторы под свечами
		currentY := v.topHeight + v.margin
		indIndex := 0

		for _, ind := range v.indicators {
			if !ind.Overlay {
				v.drawIndicator(dc, ind, currentY, currentY+indicatorHeight)
				currentY += indicatorHeight + v.margin
				indIndex++
			}
		}
	}

	// Добавляем заголовок и легенду
	v.drawTitleAndLegend(dc)

	return dc.Image(), nil
}

// RenderToFile сохраняет график в файл
func (v *Visualizer) RenderToFile(filename string) error {
	img, err := v.Render()
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

// RenderToWriter записывает график в writer
func (v *Visualizer) RenderToWriter(w io.Writer) error {
	img, err := v.Render()
	if err != nil {
		return err
	}

	return png.Encode(w, img)
}

// prepareCandles подготавливает данные свечей для отрисовки
func (v *Visualizer) prepareCandles() {
	v.candles = make([]Candle, v.series.Len())

	// Рассчитываем позиции свечей
	candleSpacing := (v.width - 2*v.margin) / len(v.candles)
	if candleSpacing < 3 {
		candleSpacing = 3
	}

	for i := 0; i < v.series.Len(); i++ {
		candle := v.series.At(i)
		v.candles[i] = Candle{
			Open:  candle.GetOpenPrice(),
			High:  candle.GetHighPrice(),
			Low:   candle.GetLowPrice(),
			Close: candle.GetClosePrice(),
			X:     v.margin + i*candleSpacing,
			Width: candleSpacing - 1,
		}
	}
}

// drawGrid рисует сетку и оси
func (v *Visualizer) drawGrid(dc *gg.Context) {
	dc.SetColor(v.colors.Grid)
	dc.SetLineWidth(1)

	// Вертикальные линии
	for x := v.margin; x <= v.width-v.margin; x += 50 {
		dc.DrawLine(float64(x), float64(v.margin), float64(x), float64(v.height-v.margin))
		dc.Stroke()
	}

	// Горизонтальные линии
	for y := v.margin; y <= v.height-v.margin; y += 50 {
		dc.DrawLine(float64(v.margin), float64(y), float64(v.width-v.margin), float64(y))
		dc.Stroke()
	}

	// Оси
	dc.SetColor(v.colors.Text)
	dc.SetLineWidth(2)

	// X-ось
	dc.DrawLine(float64(v.margin), float64(v.height-v.margin),
		float64(v.width-v.margin), float64(v.height-v.margin))

	// Y-ось
	dc.DrawLine(float64(v.margin), float64(v.margin),
		float64(v.margin), float64(v.height-v.margin))

	dc.Stroke()
}

// drawCandles рисует свечной график
func (v *Visualizer) drawCandles(dc *gg.Context) {
	if len(v.candles) == 0 {
		return
	}

	// Находим min и max цены для масштабирования
	minPrice, maxPrice := v.getPriceRange()
	priceRange := maxPrice - minPrice

	// Рисуем каждую свечу
	for _, candle := range v.candles {
		// Преобразуем цены в координаты Y
		highY := v.priceToY(candle.High, minPrice, priceRange, v.topHeight-v.margin*2)
		lowY := v.priceToY(candle.Low, minPrice, priceRange, v.topHeight-v.margin*2)
		openY := v.priceToY(candle.Open, minPrice, priceRange, v.topHeight-v.margin*2)
		closeY := v.priceToY(candle.Close, minPrice, priceRange, v.topHeight-v.margin*2)

		// Выбираем цвет свечи
		if candle.Close >= candle.Open {
			dc.SetColor(v.colors.Bullish)
		} else {
			dc.SetColor(v.colors.Bearish)
		}

		// Рисуем тень (high-low)
		dc.SetLineWidth(1)
		dc.DrawLine(
			float64(candle.X+candle.Width/2),
			float64(v.margin)+highY,
			float64(candle.X+candle.Width/2),
			float64(v.margin)+lowY,
		)
		dc.Stroke()

		// Рисуем тело свечи
		bodyTop := math.Min(openY, closeY)
		bodyBottom := math.Max(openY, closeY)
		bodyHeight := bodyBottom - bodyTop

		if bodyHeight < 1 {
			bodyHeight = 1
			bodyTop -= 0.5
		}

		dc.DrawRectangle(
			float64(candle.X),
			float64(v.margin)+bodyTop,
			float64(candle.Width),
			float64(bodyHeight),
		)

		if candle.Close >= candle.Open {
			dc.Fill()
		} else {
			dc.Stroke()
		}
	}
}

// drawOverlayIndicators рисует индикаторы поверх свечей
func (v *Visualizer) drawOverlayIndicators(dc *gg.Context) {
	minPrice, maxPrice := v.getPriceRange()
	priceRange := maxPrice - minPrice
	graphHeight := v.topHeight - v.margin*2

	for _, ind := range v.indicators {
		if ind.Overlay && len(ind.Data) > 0 {
			for lineIdx, lineData := range ind.Data {
				if len(lineData) == 0 {
					continue
				}

				// Выбираем цвет для линии
				colorIdx := lineIdx % len(v.colors.IndicatorLines)
				dc.SetColor(v.colors.IndicatorLines[colorIdx])
				dc.SetLineWidth(2)

				// Находим смещение для выравнивания с свечами
				offset := len(v.candles) - len(lineData)

				// Рисуем линию
				for i := 1; i < len(lineData); i++ {
					if i+offset < 0 || i+offset >= len(v.candles) {
						continue
					}

					x1 := float64(v.candles[i-1+offset].X + v.candles[i-1+offset].Width/2)
					y1 := float64(v.margin) + v.priceToY(lineData[i-1], minPrice, priceRange, graphHeight)

					x2 := float64(v.candles[i+offset].X + v.candles[i+offset].Width/2)
					y2 := float64(v.margin) + v.priceToY(lineData[i], minPrice, priceRange, graphHeight)

					dc.DrawLine(x1, y1, x2, y2)
					dc.Stroke()
				}
			}
		}
	}
}

// drawIndicator рисует отдельный индикатор в заданной области
func (v *Visualizer) drawIndicator(dc *gg.Context, ind IndicatorConfig, topY, bottomY int) {
	if len(ind.Data) == 0 {
		return
	}

	// Рисуем рамку для индикатора
	dc.SetColor(v.colors.Grid)
	dc.SetLineWidth(1)
	dc.DrawRectangle(float64(v.margin), float64(topY),
		float64(v.width-2*v.margin), float64(bottomY-topY))
	dc.Stroke()

	// Добавляем название индикатора
	dc.SetColor(v.colors.Text)
	dc.DrawStringAnchored(ind.Name, float64(v.width/2), float64(topY+15), 0.5, 0.5)

	// Рисуем данные индикатора
	indicatorWidth := v.width - 2*v.margin
	indicatorHeight := bottomY - topY - 30 // оставляем место для названия

	for lineIdx, lineData := range ind.Data {
		if len(lineData) == 0 {
			continue
		}

		// Находим min и max значения для масштабирования
		minVal, maxVal := v.getIndicatorRange(lineData)
		valRange := maxVal - minVal

		if valRange == 0 {
			valRange = 1
		}

		// Выбираем цвет для линии
		colorIdx := lineIdx % len(v.colors.IndicatorLines)
		dc.SetColor(v.colors.IndicatorLines[colorIdx])
		dc.SetLineWidth(1.5)

		// Рисуем линию
		for i := 1; i < len(lineData); i++ {
			x1 := float64(v.margin + (i-1)*indicatorWidth/len(lineData))
			y1 := float64(topY + 30 + indicatorHeight - int((lineData[i-1]-minVal)/valRange*float64(indicatorHeight)))

			x2 := float64(v.margin + i*indicatorWidth/len(lineData))
			y2 := float64(topY + 30 + indicatorHeight - int((lineData[i]-minVal)/valRange*float64(indicatorHeight)))

			dc.DrawLine(x1, y1, x2, y2)
			dc.Stroke()
		}

		// Для RSI рисуем уровни 30 и 70
		if ind.Type == IndicatorRSI {
			dc.SetColor(color.RGBA{150, 150, 150, 100})
			dc.SetLineWidth(0.5)

			// Уровень 30
			y30 := float64(topY + 30 + indicatorHeight - int((30-minVal)/valRange*float64(indicatorHeight)))
			dc.DrawLine(float64(v.margin), y30, float64(v.width-v.margin), y30)

			// Уровень 70
			y70 := float64(topY + 30 + indicatorHeight - int((70-minVal)/valRange*float64(indicatorHeight)))
			dc.DrawLine(float64(v.margin), y70, float64(v.width-v.margin), y70)

			dc.Stroke()
		}
	}
}

// drawTitleAndLegend добавляет заголовок и легенду
func (v *Visualizer) drawTitleAndLegend(dc *gg.Context) {
	dc.SetColor(v.colors.Text)

	// Заголовок
	dc.DrawStringAnchored("Technical Analysis Chart",
		float64(v.width/2), float64(v.margin/2), 0.5, 0.5)

	// Легенда для свечей
	legendX := v.width - v.margin - 150
	legendY := v.margin - 20

	// Бычьи свечи
	dc.SetColor(v.colors.Bullish)
	dc.DrawString("Bullish", float64(legendX), float64(legendY))

	// Медвежьи свечи
	dc.SetColor(v.colors.Bearish)
	dc.DrawString("Bearish", float64(legendX+70), float64(legendY))

	// Легенда для индикаторов
	legendY += 20
	for i, ind := range v.indicators {
		if i >= 6 { // ограничиваем количество индикаторов в легенде
			break
		}

		colorIdx := i % len(v.colors.IndicatorLines)
		dc.SetColor(v.colors.IndicatorLines[colorIdx])
		dc.DrawString(ind.Name, float64(legendX+(i%2)*70), float64(legendY+20*(i/2)))
	}
}

// Вспомогательные методы

func (v *Visualizer) getPriceRange() (min, max float64) {
	if len(v.candles) == 0 {
		return 0, 1
	}

	min = v.candles[0].Low
	max = v.candles[0].High

	for _, candle := range v.candles {
		if candle.Low < min {
			min = candle.Low
		}
		if candle.High > max {
			max = candle.High
		}
	}

	// Учитываем оверлейные индикаторы при определении диапазона
	for _, ind := range v.indicators {
		if ind.Overlay && len(ind.Data) > 0 {
			for _, lineData := range ind.Data {
				for _, val := range lineData {
					if val < min {
						min = val
					}
					if val > max {
						max = val
					}
				}
			}
		}
	}

	// Добавляем небольшой зазор
	gap := (max - min) * 0.05
	return min - gap, max + gap
}

func (v *Visualizer) getIndicatorRange(data []float64) (min, max float64) {
	if len(data) == 0 {
		return 0, 1
	}

	min = data[0]
	max = data[0]

	for _, val := range data {
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}
	}

	// Добавляем небольшой зазор
	gap := (max - min) * 0.1
	return min - gap, max + gap
}

func (v *Visualizer) priceToY(price, minPrice, priceRange float64, height int) float64 {
	normalized := (price - minPrice) / priceRange
	return float64(float64(height) * (1 - normalized))
}
