package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func parseContents(contents [][]string, column int) []float64 {
	data := []float64{}
	for _, row := range contents {
		entry, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			continue
		}
		data = append(data, entry)
	}
	return data
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE: gohisto [data_file]\n")
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("error opening file")
	}
	reader := csv.NewReader(file)

	contents, err := reader.ReadAll()
	if err != nil {
		log.Fatal("error reading file")
	}

	data := parseContents(contents, 1)

	hist := NewHistogram("Title", 10, .5, 10.5)
	hist.Fill(data...)
	hist.Draw()
}

type Histogram struct {
	data      []float64
	nBins     int
	lowBin    float64
	highBin   float64
	Bins      []int
	underflow int
	overflow  int
	Title     string
}

func NewHistogram(title string, nBins int, lowBin, highBin float64) *Histogram {
	bins := []int{}
	for i := 0; i < nBins; i++ {
		bins = append(bins, 0)
	}
	return &Histogram{
		data:      []float64{},
		nBins:     nBins,
		lowBin:    lowBin,
		highBin:   highBin,
		Bins:      bins,
		underflow: 0,
		overflow:  0,
		Title:     title,
	}
}

func (h *Histogram) Fill(values ...float64) {
	for _, val := range values {
		h.data = append(h.data, val)
		binIndex := int(float64(h.nBins) * (val - h.lowBin) / (h.highBin - h.lowBin))
		if binIndex < 0 {
			h.underflow++
		} else if binIndex >= h.nBins {
			h.overflow++
		} else {
			h.Bins[binIndex]++
		}
	}
}

func (h *Histogram) Percentiles() []float64 {
	sort.Float64s(h.data)
	entries := len(h.data)

	percentiles := []float64{}
	for i := 0; i < 10; i++ {
		percentiles = append(percentiles, h.data[entries*(i+1)/10-1])
	}
	return percentiles
}

func (h *Histogram) BinCenter(binIndex int) float64 {
	return (float64(binIndex)+.5)*(h.highBin-h.lowBin)/float64(h.nBins) + h.lowBin
}

func (h *Histogram) Bars() []int {
	return h.Bins
}

func (h *Histogram) Draw() {
	h.PrintTitle()
	h.PrintBars()
	h.PrintStats()
}

func (h *Histogram) PrintTitle() {
	fmt.Printf("%s\n", "\x1B[;1m"+h.Title+"\x1B[0m")
}

func (h *Histogram) PrintStats() {
	var entries int
	for _, bin := range h.Bins {
		entries += bin
	}
	fmt.Printf("entries: %d\n", entries)
	fmt.Printf("underflow: %d | overflow: %d\n", h.underflow, h.overflow)
	fmt.Printf("Percentiles: ")
	for i, p := range h.Percentiles() {
		fmt.Printf(` %dth  %.2f |`, 10*(i+1), p)
	}
}

func (h *Histogram) PrintBars() {
	bars := h.Bars()
	for binIndex, bar := range bars {
		for i := 0; i < 1; i++ {
			fmt.Printf("%s\n", line(40, bar, h.BinCenter(binIndex)))
		}
	}
	for i := 0; i < 40; i++ {
		fmt.Printf("-")
	}
	fmt.Printf("\n")
}

func line(width, fill int, label float64) string {
	line := fmt.Sprintf("%.1e|", label)
	for i := 0; i < width; i++ {
		if i <= fill {
			line += "\x1B[31;1m" + "\u2588" + "\x1B[0m"
		} else {
			line += " "
		}
	}
	return line
}
