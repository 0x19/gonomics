# Go client for Nomics API's

The unofficial Go client for communicating with the Nomics API.  

Nomics Website : https://nomics.com/  
Nomics API Doc : https://nomics.com/docs

Nomics provides professional-grade market data APIs to institutional crypto investors and exchanges. To our luck, this awesome sleek APIs can be accessed freely as well through their free plan. So here is Go client for you to further enrich the whole experience.

**Note:** This is an unofficial library, so some new API's included / changed ones in nomics may not reflect here in this client immediately, so be always check the official doc once, before checking in this client into your application.

Copyright (c) 2021 Pavan Shetty. Licensed under the MIT License.

## Documentation

- [Gonomics Go Client documentation - GoDoc](https://godoc.org/github.com/milkywaybrain/gonomics)
- [Nomics API documentation](https://nomics.com/docs)

**Note:** For reference, you can import all the Nomics APIs through Postman Collection in following link :
https://github.com/milkywaybrain/nomics-postman-api-collection

## Installation

```
go get github.com/milkywaybrain/gonomics
```
**Note:** Gonomics uses only golang standard library, so no external dependencies. 

## Client usage

Usage is simple. Create required request struct with all the needed values and call the function on connector with this. Then you will get the response in a specific struct with all the fields filled with values from the nomics server. In case of csv request, it saves the file on provided path on disk.

```go

package main

import (
	"fmt"
	"time"

	gonomics "github.com/milkywaybrain/gonomics"
)

const apiKey string = "demo-b5d84e505a11969a7184f899fbb40ae1"

func main() {
	// Create a new Nomics connector instance.
	c := gonomics.New(apiKey)

	// Get Currencies Sparkline.
	startTime, _ := time.Parse(time.RFC3339, "2021-02-14T00:00:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2021-02-16T00:00:00Z")
	csReq := gonomics.CurrenciesSparklineRequest{
		Ids:     []string{"BTC", "ETH"},
		Start:   startTime,
		End:     endTime,
		Convert: "USD",
	}
	csResp, err := c.GetCurrenciesSparkline(csReq)
	if err != nil {
		fmt.Printf("Error getting currencies sparkline: %v", err)
	}
	fmt.Println("currencies sparkline: ", csResp)

	// Get Markets Cap-History in CSV format.
	startTime, _ = time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	endTime, _ = time.Parse(time.RFC3339, "2021-02-01T00:00:00Z")
	mchReqCSV := gonomics.MarketsCapHistoryRequest{
		Start:            startTime,
		End:              endTime,
		Convert:          "EUR",
		Format:           "csv",
		FileNameWithPath: "/home/pavan/nomicsdata/market-cap_history.csv",
	}
	_, err = c.GetMarketsCapHistory(mchReqCSV)
	if err != nil {
		fmt.Printf("Error getting markets cap-history: %v", err)
	}
}

```
**Note:**
1. Base nomics URL and all the other API endpoints are mentioned in connector.go file.
2. See *_test.go files for examples and different usage scenario.

## Run unit tests

For all functions test :
```
go test -v
```

For nomics free-plan functions test :
```
go test -plan=free -v
```

## Donate

For Pavan Shetty, original author of this client, BTC : 1LkR7QwpKqFEd6Gdueeebfun3djLocjtuu

For Nomics, https://p.nomics.com/donate 
