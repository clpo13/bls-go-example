// bls-go-example - a command-line program to demo the bls-go library
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
	"strings"

	"github.com/clpo13/bls-go"
)

func main() {
	// Required flags
	startPtr := flag.String("start", "", "first `year` to get data for (required)")
	endPtr := flag.String("end", "", "last `year` to get data for (required)")
	seriesPtr := flag.String("series", "", "series to get data for (required)")

	// Optional flags
	avgPtr := flag.Bool("avg", false, "request annual average of monthly values")
	calcPtr := flag.Bool("calc", false, "request data calculations")
	catPtr := flag.Bool("cat", false, "request series catalog data")

	// The API call will still work with no API key, but the results are limited.
	keyPtr := flag.String("key", "", "API key to use")

	flag.Parse()

	// Start year, end year, and series ID are required, so print an error and
	// quit if they aren't found.
	if *startPtr == "" || *endPtr == "" || *seriesPtr == "" {
		fmt.Println("Missing a required flag!")
		fmt.Println("Try 'bls-go-example --help' for more information.")
		os.Exit(1)
	}

	fmt.Printf("Querying series %s for years %s through %s...\n\n", *seriesPtr, *startPtr, *endPtr)

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

	// Send the payload to the API.
	tr := blsgo.GetData(payload)

	// Eventually, this error handling will be taken care of by the library
	// itself, which will send an error object if the server doesn't send
	// what we expect.
	if tr.Status != "REQUEST_SUCCEEDED" {
		fmt.Println("Server error:", tr.Status)
		for _, v := range tr.Message {
			fmt.Println(v)
		}
		os.Exit(1)
	}

	// Print out any messages.
	if len(tr.Message) > 0 {
		for _, v := range tr.Message {
			fmt.Println(v)
		}
		// If the first message is about an invalid series, quit with an error.
		if strings.HasPrefix(tr.Message[0], "Invalid Series") {
			os.Exit(1)
		}
	}

	// Catalog data is optional, so only print it out if we got it.
	// TODO: not all fields are present in all series, so make them optional
	// as well.
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

	// TODO: print out more than just the first period.
	data := tr.Results.Series[0].Data[0]
	fmt.Println("Year:", data.Year)
	fmt.Println("Period:", data.Num)
	fmt.Println("Period name:", data.Name)
	fmt.Println("Value:", data.Value)

	// An empty array of footnotes is always returned, but don't print anything
	// if we didn't get any actual footnotes.
	fn := data.Footnotes
	if fn[0].Code != nil {
		fmt.Println("Footnotes:")
		for _, v := range fn {
			fmt.Println("\t", *v.Code, "-", *v.Text)
		}
	}

	// Calculations are optional, so only print them if we received them.
	// TODO: some series send only net or pct change, not both.
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
