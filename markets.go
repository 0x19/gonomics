package gonomics

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

// Markets.

// MarketsRequest represents market requests parameters.
type MarketsRequest struct {
	Exchange string
	Base     []string
	Quote    []string

	// If the format is csv, please also include fileNameWithPath field value.
	// Default is json.
	Format string

	// csv file name. Required, if the format is csv.
	// example : /home/user/nomicsdata/markets_1613046296.csv.
	// Here, new "markets_1613046296.csv" file will be created on existing, /home/user/nomicsdata/, directory.
	FileNameWithPath string
}

// MarketsResponse represents markets response.
// This will be nil, if the format requested is csv.
// Fields will contain default go lang values if there is no value received from the server.
type MarketsResponse struct {
	Exchange string `json:"exchange"`
	Market   string `json:"market"`
	Base     string `json:"base"`
	Quote    string `json:"quote"`
}

// GetMarkets fetches the markets from the server and returns array of
// MarketsResponse if the requested format is json, otherwise creates a csv file on disk, format is csv.
// Note : in case of csv format, MarketsRequest.FileNameWithPath is required and the []MarketsResponse return data is nil.
func (c *Connecter) GetMarkets(mReq MarketsRequest) ([]MarketsResponse, error) {
	req, err := c.newRequest(marketsURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if mReq.Exchange != "" {
		q.Add("exchange", mReq.Exchange)
	}
	if len(mReq.Base) > 0 {
		q.Add("base", strings.Join(mReq.Base[:], ","))
	}
	if len(mReq.Quote) > 0 {
		q.Add("quote", strings.Join(mReq.Quote[:], ","))
	}
	if mReq.Format != "" {
		if mReq.Format == "csv" && mReq.FileNameWithPath == "" {
			return nil, errors.New("csv file path is required")
		}
		q.Add("format", mReq.Format)
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response, if the requested format is json.
	if mReq.Format == "" || mReq.Format == "json" {
		var mResp []MarketsResponse
		if err := json.NewDecoder(resp.Body).Decode(&mResp); err != nil {
			return nil, err
		}
		return mResp, nil
	}

	// Creates a CSV file and then copies the server's response, if the requested format is csv.
	err = c.createFile(resp.Body, mReq.FileNameWithPath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// MarketsCap History.

// MarketsCapHistoryRequest represents markets-cap history request parameters.
type MarketsCapHistoryRequest struct {
	Start   time.Time
	End     time.Time
	Convert string

	// If the format is csv, please also include fileNameWithPath field value.
	// Default is json.
	Format string

	// csv file name. Required, if the format is csv.
	// example : /home/user/nomicsdata/markets-cap_history_1613046296.csv.
	// Here, new "markets-cap_history_1613046296.csv" file will be created on existing, /home/user/nomicsdata/, directory.
	FileNameWithPath string

	IncludeTransparency bool
}

// MarketsCapHistoryResponse represents markets-cap history response.
// This will be nil, if the format requested is csv.
// Fields will contain default go lang values if there is no value received from the server.
type MarketsCapHistoryResponse struct {
	Timestamp            time.Time `json:"timestamp"`
	MarketCap            float64   `json:"market_cap,string"`
	TransparentMarketCap float64   `json:"transparent_market_cap,string"`
}

// GetMarketsCapHistory fetches the markets-cap history from the server and returns array of
// MarketsCapHistoryResponse if the requested format is json, otherwise creates a csv file on disk, format is csv.
// Note : in case of csv format, MarketsCapHistoryRequest.FileNameWithPath is required
// and the []MarketsCapHistoryResponse return data is nil.
func (c *Connecter) GetMarketsCapHistory(mchReq MarketsCapHistoryRequest) ([]MarketsCapHistoryResponse, error) {
	req, err := c.newRequest(marketsCapHistoryURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if mchReq.Start.IsZero() {
		return nil, errors.New("start is required")
	}
	q.Add("start", mchReq.Start.Format(time.RFC3339))
	if !mchReq.End.IsZero() {
		q.Add("end", mchReq.End.Format(time.RFC3339))
	}
	if mchReq.Convert != "" {
		q.Add("convert", mchReq.Convert)
	}
	if mchReq.Format != "" {
		if mchReq.Format == "csv" && mchReq.FileNameWithPath == "" {
			return nil, errors.New("csv file path is required")
		}
		q.Add("format", mchReq.Format)
	}
	if mchReq.IncludeTransparency {
		q.Add("include-transparency", strconv.FormatBool(mchReq.IncludeTransparency))
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response, if the requested format is json.
	if mchReq.Format == "" || mchReq.Format == "json" {
		var mchResp []MarketsCapHistoryResponse
		if err := json.NewDecoder(resp.Body).Decode(&mchResp); err != nil {
			return nil, err
		}
		return mchResp, nil
	}

	// Creates a CSV file and then copies the server's response, if the requested format is csv.
	err = c.createFile(resp.Body, mchReq.FileNameWithPath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Exchange Markets Ticker.

// ExchangeMarketsTickerRequest represents exchange-markets ticker request parameters.
type ExchangeMarketsTickerRequest struct {
	Interval []string
	Currency []string
	Base     []string
	Quote    []string
	Exchange []string
	Market   []string
	Convert  string
	Status   string
	Search   string
	PerPage  int
	Page     int
}

// ExchangeMarketsTickerResponse represents exchange-markets ticker history response.
type ExchangeMarketsTickerResponse struct {
	Exchange      string                                `json:"exchange"`
	Market        string                                `json:"market"`
	Type          string                                `json:"type"`
	SubType       string                                `json:"subtype"`
	Aggregated    bool                                  `json:"aggregated"`
	PriceExclude  bool                                  `json:"price_exclude"`
	VolumeExclude bool                                  `json:"volume_exclude"`
	Base          string                                `json:"base"`
	Quote         string                                `json:"quote"`
	BaseSymbol    string                                `json:"base_symbol"`
	QuoteSymbol   string                                `json:"quote_symbol"`
	Price         float64                               `json:"price,string"`
	PriceQuote    float64                               `json:"price_quote,string"`
	VolumeUSD     float64                               `json:"volume_usd,string"`
	LastUpdated   time.Time                             `json:"last_updated"`
	OneH          ExchangeMarketsTickerIntervalResponse `json:"1h"`
	OneD          ExchangeMarketsTickerIntervalResponse `json:"1d"`
	SevenD        ExchangeMarketsTickerIntervalResponse `json:"7d"`
	Three0D       ExchangeMarketsTickerIntervalResponse `json:"30d"`
	Three65D      ExchangeMarketsTickerIntervalResponse `json:"365d"`
	Ytd           ExchangeMarketsTickerIntervalResponse `json:"ytd"`
}

// ExchangeMarketsTickerIntervalResponse represents exchange-markets ticker interval response,
// Included in ExchangeMarketsTickerResponse.
// Fields will contain default go lang values if there is no value received from the server.
type ExchangeMarketsTickerIntervalResponse struct {
	Volume           float64 `json:"volume,string"`
	VolumeBase       float64 `json:"volume_base,string"`
	VolumeChange     float64 `json:"volume_change,string"`
	VolumeBaseChange float64 `json:"volume_base_change,string"`
	Trades           float64 `json:"trades,string"`
	TradesChange     float64 `json:"trades_change,string"`
	PriceChange      float64 `json:"price_change,string"`
	PriceQuoteChange float64 `json:"price_quote_change,string"`
}

// GetExchangeMarketsTicker fetches the exchange-markets ticker from the server
// and returns array of ExchangeMarketsTickerResponse.
func (c *Connecter) GetExchangeMarketsTicker(emtReq ExchangeMarketsTickerRequest) ([]ExchangeMarketsTickerResponse, error) {
	req, err := c.newRequest(exchangeMarketsTickerURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if len(emtReq.Interval) > 0 {
		q.Add("interval", strings.Join(emtReq.Interval[:], ","))
	}
	if len(emtReq.Currency) > 0 {
		q.Add("currency", strings.Join(emtReq.Currency[:], ","))
	}
	if len(emtReq.Base) > 0 {
		q.Add("base", strings.Join(emtReq.Base[:], ","))
	}
	if len(emtReq.Quote) > 0 {
		q.Add("quote", strings.Join(emtReq.Quote[:], ","))
	}
	if len(emtReq.Exchange) > 0 {
		q.Add("exchange", strings.Join(emtReq.Exchange[:], ","))
	}
	if len(emtReq.Market) > 0 {
		q.Add("market", strings.Join(emtReq.Market[:], ","))
	}
	if emtReq.Convert != "" {
		q.Add("Convert", emtReq.Convert)
	}
	if emtReq.Status != "" {
		q.Add("status", emtReq.Status)
	}
	if emtReq.Search != "" {
		q.Add("search", emtReq.Search)
	}
	if emtReq.PerPage != 0 {
		q.Add("per-page", strconv.Itoa(emtReq.PerPage))
	}
	if emtReq.Page != 0 {
		q.Add("page", strconv.Itoa(emtReq.Page))
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response.
	var emtResp []ExchangeMarketsTickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&emtResp); err != nil {
		return nil, err
	}
	return emtResp, nil
}
