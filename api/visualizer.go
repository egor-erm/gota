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

var (
	RED        = color.RGBA{255, 80, 80, 255}
	GREEN      = color.RGBA{0, 200, 100, 255}
	BLUE       = color.RGBA{100, 100, 255, 255}
	YELLOW     = color.RGBA{255, 200, 50, 255}
	PURPLE     = color.RGBA{200, 100, 255, 255}
	CYAN       = color.RGBA{100, 200, 255, 255}
	ORANGE     = color.RGBA{255, 150, 50, 255}
	PINK       = color.RGBA{255, 100, 200, 255}
	WHITE      = color.RGBA{255, 255, 255, 255}
	GRAY       = color.RGBA{150, 150, 150, 255}
	DARK_GRAY  = color.RGBA{60, 60, 80, 255}
	LIGHT_GRAY = color.RGBA{200, 200, 220, 255}
	BLACK      = color.RGBA{0, 0, 0, 255}

	// Полупрозрачные версии
	RED_TRANSPARENT    = color.RGBA{255, 80, 80, 150}
	GREEN_TRANSPARENT  = color.RGBA{0, 200, 100, 150}
	BLUE_TRANSPARENT   = color.RGBA{100, 100, 255, 150}
	YELLOW_TRANSPARENT = color.RGBA{255, 200, 50, 150}
	PURPLE_TRANSPARENT = color.RGBA{200, 100, 255, 150}
	GRAY_TRANSPARENT   = color.RGBA{150, 150, 150, 100}
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
	Name      string
	Type      IndicatorType
	Data      [][]float64   // данные индикатора (может быть несколько линий)
	Colors    []color.Color // цвета для каждой линии
	Labels    []string      // названия линий
	LineWidth float64       // толщина линии
	Overlay   bool          // отображать ли индикатор поверх свечей
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
func (v *Visualizer) AddSMA(period int, c color.Color) {
	sma := NewAnalyzer(v.series).SMA(period)
	if len(sma) > 0 {
		// Обрезаем данные, чтобы они соответствовали свечам
		alignedSma := v.alignIndicatorData(sma)
		v.AddIndicator(IndicatorConfig{
			Name:      fmt.Sprintf("SMA(%d)", period),
			Type:      IndicatorSMA,
			Data:      [][]float64{alignedSma},
			Colors:    []color.Color{c},
			Labels:    []string{"SMA"},
			LineWidth: 2.0,
			Overlay:   true,
		})
	}
}

// AddEMA добавляет EMA индикатор
func (v *Visualizer) AddEMA(period int, c color.Color) {
	ema := NewAnalyzer(v.series).EMA(period)
	if len(ema) > 0 {
		// Обрезаем данные, чтобы они соответствовали свечам
		alignedEma := v.alignIndicatorData(ema)
		v.AddIndicator(IndicatorConfig{
			Name:      fmt.Sprintf("EMA(%d)", period),
			Type:      IndicatorEMA,
			Data:      [][]float64{alignedEma},
			Colors:    []color.Color{c},
			Labels:    []string{"EMA"},
			LineWidth: 2.0,
			Overlay:   true,
		})
	}
}

// AddRSI добавляет RSI индикатор
func (v *Visualizer) AddRSI(period int, c color.Color) {
	rsi := NewAnalyzer(v.series).RSI(period)
	if len(rsi) > 0 {
		// Обрезаем данные, чтобы они соответствовали свечам
		alignedRsi := v.alignIndicatorData(rsi)
		v.AddIndicator(IndicatorConfig{
			Name:      fmt.Sprintf("RSI(%d)", period),
			Type:      IndicatorRSI,
			Data:      [][]float64{alignedRsi},
			Colors:    []color.Color{c},
			Labels:    []string{"RSI"},
			LineWidth: 1.5,
			Overlay:   false,
		})
	}
}

// AddMACD добавляет MACD индикатор
func (v *Visualizer) AddMACD(fast, slow, signal int, macdColor, signalColor, histColor color.Color) {
	macdLine, signalLine, histogram := NewAnalyzer(v.series).MACD(fast, slow, signal)
	if macdLine != nil {
		// Обрезаем данные, чтобы они соответствовали свечам
		alignedMacdLine := v.alignIndicatorData(macdLine)
		alignedSignalLine := v.alignIndicatorData(signalLine)
		alignedHistogram := v.alignIndicatorData(histogram)

		v.AddIndicator(IndicatorConfig{
			Name:      fmt.Sprintf("MACD(%d,%d,%d)", fast, slow, signal),
			Type:      IndicatorMACD,
			Data:      [][]float64{alignedMacdLine, alignedSignalLine, alignedHistogram},
			Colors:    []color.Color{macdColor, signalColor, histColor},
			Labels:    []string{"MACD", "Signal", "Histogram"},
			LineWidth: 1.5,
			Overlay:   false,
		})
	}
}

// AddBollingerBands добавляет Bollinger Bands
func (v *Visualizer) AddBollingerBands(period int, stdDev float64, upperColor, middleColor, lowerColor color.Color) {
	upper, middle, lower := NewAnalyzer(v.series).BollingerBands(period, stdDev)
	if upper != nil {
		// Обрезаем данные, чтобы они соответствовали свечам
		alignedUpper := v.alignIndicatorData(upper)
		alignedMiddle := v.alignIndicatorData(middle)
		alignedLower := v.alignIndicatorData(lower)

		v.AddIndicator(IndicatorConfig{
			Name:      fmt.Sprintf("BB(%d,%.1f)", period, stdDev),
			Type:      IndicatorBB,
			Data:      [][]float64{alignedUpper, alignedMiddle, alignedLower},
			Colors:    []color.Color{upperColor, middleColor, lowerColor},
			Labels:    []string{"Upper", "Middle", "Lower"},
			LineWidth: 1.5,
			Overlay:   true,
		})
	}
}

// AddStochRSI добавляет StochRSI индикатор
func (v *Visualizer) AddStochRSI(rsiPeriod, stochPeriod, smoothK, smoothD int, kColor, dColor color.Color) {
	kLine, dLine := NewAnalyzer(v.series).StochRSI(rsiPeriod, stochPeriod, smoothK, smoothD)
	if len(kLine) > 0 {
		// Обрезаем данные, чтобы они соответствовали свечам
		alignedKLine := v.alignIndicatorData(kLine)
		alignedDLine := v.alignIndicatorData(dLine)

		v.AddIndicator(IndicatorConfig{
			Name:      fmt.Sprintf("StochRSI(%d,%d,%d,%d)", rsiPeriod, stochPeriod, smoothK, smoothD),
			Type:      IndicatorStochRSI,
			Data:      [][]float64{alignedKLine, alignedDLine},
			Colors:    []color.Color{kColor, dColor},
			Labels:    []string{"%K", "%D"},
			LineWidth: 1.5,
			Overlay:   false,
		})
	}
}

// AddATR добавляет ATR индикатор
func (v *Visualizer) AddATR(period int, c color.Color) {
	atr := NewAnalyzer(v.series).ATR(period)
	if len(atr) > 0 {
		// Обрезаем данные, чтобы они соответствовали свечам
		alignedAtr := v.alignIndicatorData(atr)
		v.AddIndicator(IndicatorConfig{
			Name:      fmt.Sprintf("ATR(%d)", period),
			Type:      IndicatorATR,
			Data:      [][]float64{alignedAtr},
			Colors:    []color.Color{c},
			Labels:    []string{"ATR"},
			LineWidth: 1.5,
			Overlay:   false,
		})
	}
}

// alignIndicatorData обрезает данные индикатора для соответствия свечам
func (v *Visualizer) alignIndicatorData(indicatorData []float64) []float64 {
	if len(indicatorData) == 0 {
		return nil
	}

	// Вычисляем смещение между данными индикатора и свечами
	// Индикаторы обычно имеют задержку из-за периода расчета
	offset := v.series.Len() - len(indicatorData)

	// Если индикатор короче свечей, добавляем нули в начало
	if offset > 0 {
		aligned := make([]float64, v.series.Len())
		for i := 0; i < offset; i++ {
			aligned[i] = math.NaN() // Используем NaN для пропуска точек
		}
		copy(aligned[offset:], indicatorData)
		return aligned
	}

	// Если индикатор длиннее свечей, обрезаем начало
	if offset < 0 {
		return indicatorData[-offset:]
	}

	// Если длины равны
	return indicatorData
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
		availableHeight := v.height - v.topHeight - v.margin*2
		indicatorHeight := availableHeight / max(indicatorCount, 1)

		// Рисуем отдельные индикаторы под свечами
		currentY := v.topHeight + v.margin
		indIndex := 0

		for _, ind := range v.indicators {
			if !ind.Overlay {
				v.drawIndicator(dc, ind, currentY, currentY+indicatorHeight)
				currentY += indicatorHeight + v.margin/2
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
	candleSpacing := (v.width - 2*v.margin) / max(len(v.candles), 1)
	if candleSpacing < 3 {
		candleSpacing = 3
	}
	v.candleWidth = max(candleSpacing-2, 1)

	for i := 0; i < v.series.Len(); i++ {
		candle := v.series.At(i)
		v.candles[i] = Candle{
			Open:  candle.GetOpenPrice(),
			High:  candle.GetHighPrice(),
			Low:   candle.GetLowPrice(),
			Close: candle.GetClosePrice(),
			X:     v.margin + i*candleSpacing + candleSpacing/2,
			Width: v.candleWidth,
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
	graphHeight := float64(v.topHeight - v.margin*2)

	// Рисуем каждую свечу
	for _, candle := range v.candles {
		// Преобразуем цены в координаты Y
		highY := v.priceToY(candle.High, minPrice, priceRange, graphHeight)
		lowY := v.priceToY(candle.Low, minPrice, priceRange, graphHeight)
		openY := v.priceToY(candle.Open, minPrice, priceRange, graphHeight)
		closeY := v.priceToY(candle.Close, minPrice, priceRange, graphHeight)

		// Выбираем цвет свечи
		if candle.Close >= candle.Open {
			dc.SetColor(v.colors.Bullish)
		} else {
			dc.SetColor(v.colors.Bearish)
		}

		// Рисуем тень (high-low)
		dc.SetLineWidth(1)
		dc.DrawLine(
			float64(candle.X),
			float64(v.margin)+highY,
			float64(candle.X),
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
			float64(candle.X-candle.Width/2),
			float64(v.margin)+bodyTop,
			float64(candle.Width),
			bodyHeight,
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
	graphHeight := float64(v.topHeight - v.margin*2)

	for _, ind := range v.indicators {
		if ind.Overlay && len(ind.Data) > 0 {
			for lineIdx, lineData := range ind.Data {
				if len(lineData) == 0 || len(lineData) != len(v.candles) {
					continue
				}

				// Выбираем цвет для линии
				var lineColor color.Color
				if lineIdx < len(ind.Colors) {
					lineColor = ind.Colors[lineIdx]
				} else {
					lineColor = v.colors.IndicatorLines[lineIdx%len(v.colors.IndicatorLines)]
				}

				dc.SetColor(lineColor)
				dc.SetLineWidth(ind.LineWidth)

				// Рисуем линию, пропуская NaN значения
				startX, startY := -1.0, -1.0
				for i, val := range lineData {
					if math.IsNaN(val) {
						startX, startY = -1.0, -1.0
						continue
					}

					x := float64(v.candles[i].X)
					y := float64(v.margin) + v.priceToY(val, minPrice, priceRange, graphHeight)

					if startX >= 0 && startY >= 0 {
						dc.DrawLine(startX, startY, x, y)
						dc.Stroke()
					}

					startX, startY = x, y
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
	if err := dc.LoadFontFace("", 12); err == nil {
		dc.DrawStringAnchored(ind.Name, float64(v.width/2), float64(topY+15), 0.5, 0.5)
	}

	// Рисуем данные индикатора
	indicatorWidth := float64(v.width - 2*v.margin)
	indicatorHeight := float64(bottomY - topY - 30) // оставляем место для названия

	for lineIdx, lineData := range ind.Data {
		if len(lineData) == 0 || len(lineData) != len(v.candles) {
			continue
		}

		// Находим min и max значения для масштабирования
		minVal, maxVal := v.getIndicatorRange(lineData)
		valRange := maxVal - minVal

		if valRange == 0 {
			valRange = 1
		}

		// Выбираем цвет для линии
		var lineColor color.Color
		if lineIdx < len(ind.Colors) {
			lineColor = ind.Colors[lineIdx]
		} else {
			lineColor = v.colors.IndicatorLines[lineIdx%len(v.colors.IndicatorLines)]
		}

		dc.SetColor(lineColor)
		dc.SetLineWidth(ind.LineWidth)

		// Рисуем линию, пропуская NaN значения
		startX, startY := -1.0, -1.0
		for i, val := range lineData {
			if math.IsNaN(val) {
				startX, startY = -1.0, -1.0
				continue
			}

			x := float64(v.margin) + (float64(i)/float64(len(lineData)-1))*indicatorWidth
			y := float64(topY) + 30 + indicatorHeight - ((val-minVal)/valRange)*indicatorHeight

			if startX >= 0 && startY >= 0 {
				dc.DrawLine(startX, startY, x, y)
				dc.Stroke()
			}

			startX, startY = x, y
		}

		// Для RSI и StochRSI рисуем уровни
		if ind.Type == IndicatorRSI {
			dc.SetColor(color.RGBA{150, 150, 150, 100})
			dc.SetLineWidth(0.5)

			// Уровень 30
			y30 := float64(topY) + 30 + indicatorHeight - ((30-minVal)/valRange)*indicatorHeight
			dc.DrawLine(float64(v.margin), y30, float64(v.width-v.margin), y30)

			// Уровень 70
			y70 := float64(topY) + 30 + indicatorHeight - ((70-minVal)/valRange)*indicatorHeight
			dc.DrawLine(float64(v.margin), y70, float64(v.width-v.margin), y70)

			dc.Stroke()
		} else if ind.Type == IndicatorStochRSI {
			dc.SetColor(color.RGBA{150, 150, 150, 100})
			dc.SetLineWidth(0.5)

			// Уровень 20
			y20 := float64(topY) + 30 + indicatorHeight - ((20-minVal)/valRange)*indicatorHeight
			dc.DrawLine(float64(v.margin), y20, float64(v.width-v.margin), y20)

			// Уровень 80
			y80 := float64(topY) + 30 + indicatorHeight - ((80-minVal)/valRange)*indicatorHeight
			dc.DrawLine(float64(v.margin), y80, float64(v.width-v.margin), y80)

			dc.Stroke()
		}

		// Для MACD рисуем нулевую линию
		if ind.Type == IndicatorMACD && lineIdx == 0 {
			dc.SetColor(color.RGBA{150, 150, 150, 100})
			dc.SetLineWidth(0.5)

			zeroY := float64(topY) + 30 + indicatorHeight - ((0-minVal)/valRange)*indicatorHeight
			dc.DrawLine(float64(v.margin), zeroY, float64(v.width-v.margin), zeroY)
			dc.Stroke()
		}
	}
}

// drawTitleAndLegend добавляет заголовок и легенду
func (v *Visualizer) drawTitleAndLegend(dc *gg.Context) {
	dc.SetColor(v.colors.Text)

	// Загружаем шрифт (если доступен)
	if err := dc.LoadFontFace("", 14); err == nil {
		// Заголовок
		dc.DrawStringAnchored("Technical Analysis Chart",
			float64(v.width/2), float64(v.margin/2), 0.5, 0.5)

		// Легенда
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

			var indicatorColor color.Color
			if len(ind.Colors) > 0 {
				indicatorColor = ind.Colors[0]
			} else {
				indicatorColor = v.colors.IndicatorLines[i%len(v.colors.IndicatorLines)]
			}

			dc.SetColor(indicatorColor)

			// Обрезаем длинное имя для легенды
			name := ind.Name
			if len(name) > 15 {
				name = name[:12] + "..."
			}

			dc.DrawString(name,
				float64(legendX+(i%2)*70),
				float64(legendY+20*(i/2)))
		}
	}
}

// Вспомогательные методы

func (v *Visualizer) getPriceRange() (min, max float64) {
	if len(v.candles) == 0 {
		return 0, 1
	}

	min = math.Inf(1)
	max = math.Inf(-1)

	for _, candle := range v.candles {
		min = math.Min(min, candle.Low)
		max = math.Max(max, candle.High)
	}

	// Учитываем оверлейные индикаторы при определении диапазона
	for _, ind := range v.indicators {
		if ind.Overlay && len(ind.Data) > 0 {
			for _, lineData := range ind.Data {
				for _, val := range lineData {
					if !math.IsNaN(val) {
						min = math.Min(min, val)
						max = math.Max(max, val)
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

	min = math.Inf(1)
	max = math.Inf(-1)

	for _, val := range data {
		if !math.IsNaN(val) {
			min = math.Min(min, val)
			max = math.Max(max, val)
		}
	}

	// Добавляем небольшой зазор
	gap := (max - min) * 0.1
	if gap == 0 {
		gap = 1
	}
	return min - gap, max + gap
}

func (v *Visualizer) priceToY(price, minPrice, priceRange, height float64) float64 {
	normalized := (price - minPrice) / priceRange
	return height * (1 - normalized)
}
