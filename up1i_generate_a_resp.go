package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/lucasb-eyer/go-colorful"
)

type DataPoint struct {
	Label string  `csv:"label"`
	Value float64 `csv:"value"`
}

type RespVisualization struct {
	Title  string      `json:"title"`
	Data   []DataPoint `json:"data"`
	Colors []string    `json:"colors"`
}

func main() {
	http.HandleFunc("/viz", handleViz)
	http.ListenAndServe(":8080", nil)
}

func handleViz(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("data.csv")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var dataPoints []DataPoint
	if err := gocsv.UnmarshalFile(file, &dataPoints); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	colors := make([]string, len(dataPoints))
	for i := range dataPoints {
		color := colorful.Hsv(float64(i)/float64(len(dataPoints)), 1, 1)
		colors[i] = color.Hex()
	}

	viz := RespVisualization{
		Title:  "Responsive Data Visualization",
		Data:   dataPoints,
		Colors: colors,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(viz)
}