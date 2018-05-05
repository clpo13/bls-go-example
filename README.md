# bls-go-calc

This is a simple demo program that uses the [bls-go](https://github.com/clpo13/bls-go)
library to query BLS.gov for labor statistics.

## Requirements

- [Go](https://golang.org)
- BLS.gov API key (not strictly required, but highly recommended; you can
    request one [here](https://data.bls.gov/registrationEngine/))

Note: without an API key, the request won't return the annual average, calculations,
or catalog data, and you're subject to stricter daily query limits.

## Installation

```bash
git clone https://github.com/clpo13/bls-go-calc.git
cd bls-go-calc
go install
```

## Usage

```txt
Usage of bls-calc:
  -start year
        first year to get data for (required)
  -end year
        last year to get data for (required)
  -series string
        series to get data for (required)
  -key string
        API key to use
  -avg
        request annual average of monthly values
  -calc
        request data calculations
  -cat
        request series catalog data
```

The program will query the BLS.gov API and print out the catalog data of the
given series (if requested) and the first period of data, along with calculations,
if the `-calc` flag was used.

## Contributing

Issues and pull requests are always welcome. Please file any bug reports using
the GitHub [issues page](https://github.com/clpo13/bls-go-calc/issues).

## License

This program is available under the terms of the Apache 2.0 license, the text
of which can be found in [LICENSE](LICENSE) or at
<https://www.apache.org/licenses/LICENSE-2.0>.
