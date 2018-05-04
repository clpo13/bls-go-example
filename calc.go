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
	"flag"
	"fmt"
	"os"

	"github.com/clpo13/bls-go"
)

func main() {
	startPtr := flag.String("start", "", "first `year` to get data for")
	endPtr := flag.String("end", "", "last `year` to get data for")
	seriesPtr := flag.String("series", "", "series to get data for")

	avgPtr := flag.Bool("avg", false, "request annual average of monthly values")
	calcPtr := flag.Bool("calc", false, "request data calculations")
	catPtr := flag.Bool("cat", false, "request series catalog data")

	keyPtr := flag.String("key", "", "API key to use")

	flag.Parse()

	if *startPtr == "" {
		fmt.Println("Need a start year")
		os.Exit(1)
	}

	if *endPtr == "" {
		fmt.Println("Need an end year")
		os.Exit(1)
	}

	if *seriesPtr == "" {
		fmt.Println("Need a series")
		os.Exit(1)
	}

	fmt.Printf("Querying series %s for years %s through %s...\n", *seriesPtr, *startPtr, *endPtr)

	// Create a JSON payload.
	seriesArray := []string{*seriesPtr} // Convert seriesID to an array since that's what the API expects.
	payload := blsgo.Payload{
	  Start:   *startPtr,
	  End:     *endPtr,
	  Series:  seriesArray,
		Catalog: *catPtr,
		Calc:    *calcPtr,
	  Avg:     *avgPtr,
	  Key:     *keyPtr,
	}

	tr := blsgo.GetData(payload)

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
