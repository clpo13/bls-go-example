// Copyright 2018 Cody Logan
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	 http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/clpo13/bls-go"
)

// Footnote represents arbitrary footnotes sent with the received data.
type Footnote struct {
	Code *string `json:",omitempty"`
	Text *string `json:",omitempty"`
}

// Change represents calculations of change over 1, 3, 6, and 12 months.
type Change struct {
	OneMonth    string `json:"1"`
	ThreeMonth  string `json:"3"`
	SixMonth    string `json:"6"`
	TwelveMonth string `json:"12"`
}

// Calculation represents calculations sent with the received data.
type Calculation struct {
	NetChange Change `json:"net_changes"`
	PctChange Change `json:"pct_changes"`
}

// Period represents an individual period (usually a month) of data.
type Period struct {
	Year         string      `json:"year"`
	Num          string      `json:"period"`
	Name         string      `json:"periodName"`
	Value        string      `json:"value"`
	Footnotes    []Footnote  `json:"footnotes"`
	Calculations *Calculation `json:"calculations,omitempty"`
}

// Catalog is...
type Catalog struct {
	Title	 	 string `json:"series_title"`
	ID			 string `json:"series_id"`
	Season	 string `json:"seasonality"`
	Name		 string `json:"survey_name"`
	Abbr		 string `json:"survey_abbreviation"`
	DataType string `json:"measure_data_type"`
	Area     string `json:"area"`
	AreaType string `json:"area_type"`
}

// SeriesData represents data from a single series.
type SeriesData struct {
	SeriesID string
	Catalog  *Catalog  `json:"catalog,omitempty"`
	Data     []Period
}

// Series represents a JSON object containing arrays of series.
type Series struct {
	Series []SeriesData
}

// ResultData represents the top-level structure of the received data.
type ResultData struct {
	Status       string
	ResponseTime int
	Message      []string
	Results      Series
}

func main() {
	f, err := os.Open("testdata2.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	var tr ResultData
	err = dec.Decode(&tr)
	if err != nil {
		panic(err)
	}

	catalog := tr.Results.Series[0].Catalog
	if catalog != nil {
		fmt.Println("Series title:", catalog.Title)
		fmt.Println("Series ID:", catalog.ID)
		fmt.Println("Seasonality:", catalog.Season)
		fmt.Println("Survey name:", catalog.Name)
		fmt.Println("Survey abbreviation:", catalog.Abbr)
		fmt.Println("Measure data type:", catalog.DataType)
		fmt.Println("Area:", catalog.Area)
		fmt.Println("Area type:", catalog.AreaType)
		fmt.Println()
	}

	data := tr.Results.Series[0].Data[0]
	fmt.Println("Year:", data.Year)
	fmt.Println("Period:", data.Num)
	fmt.Println("Period name:", data.Name)
	fmt.Println("Value:", data.Value)

	fn := data.Footnotes
	if fn[0].Code != nil {
		fmt.Println("Footnotes:")
		for _, v := range fn {
			fmt.Println("\t", *v.Code, "-", *v.Text)
		}
	}

	calcs := data.Calculations
	if calcs != nil {
		net := calcs.NetChange
		pct := calcs.PctChange
		fmt.Println("Net change:")
		fmt.Println("\t1 month:", net.OneMonth)
		fmt.Println("\t3 months:", net.ThreeMonth)
		fmt.Println("\t6 months:", net.SixMonth)
		fmt.Println("\t12 months:", net.TwelveMonth)
		fmt.Println("Percent change:")
		fmt.Println("\t1 month:", pct.OneMonth)
		fmt.Println("\t3 months:", pct.ThreeMonth)
		fmt.Println("\t6 months:", pct.SixMonth)
		fmt.Println("\t12 months:", pct.TwelveMonth)
	}
}
