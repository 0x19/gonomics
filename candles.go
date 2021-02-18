package gonomics

import (
	"encoding/json"
	"errors"
	"time"
)

// Candles.

// CandlesRequest represents candles request parameters.
type CandlesRequest struct {
	Interval string
	Currency string
	Start    time.Time
	End      time.Time

	// If the format is csv, please also include fileNameWithPath field value.
	// Default is json.
	Format string

	// csv file name. Required, if the format is csv.
	// example : /home/user/nomicsdata/candles_1613046296.csv.
	// Here, new "candles_1613046296.csv" file will be created on existing, /home/user/nomicsdata/, directory.
	FileNameWithPath string
}

// CandlesResponse represents candles response.
// This will be nil, if the format requested is csv.
// Fields will contain default go lang values if there is no value received from the server.
type CandlesResponse struct {
	Timestamp          time.Time                         `json:"timestamp"`
	Open               float64                           `json:"open,string"`
	High               float64                           `json:"high,string"`
	Low                float64                           `json:"low,string"`
	Close              float64                           `json:"close,string"`
	Volume             float64                           `json:"volume,string"`
	TransparentOpen    float64                           `json:"transparent_open,string"`
	TransparentHigh    float64                           `json:"transparent_high,string"`
	TransparentLow     float64                           `json:"transparent_low,string"`
	TransparentClose   float64                           `json:"transparent_close,string"`
	TransparentVolume  float64                           `json:"transparent_volume,string"`
	VolumeTransparency CandlesVolumeTransparencyResponse `json:"volume_transparency"`
}

// CandlesVolumeTransparencyResponse represents candles volume transparency response, Included in CandlesResponse.
// Fields will contain default go lang values if there is no value received from the server.
type CandlesVolumeTransparencyResponse struct {
	Others float64 `json:"?,string"`
	A      float64 `json:"A,string"`
	B      float64 `json:"B,string"`
	C      float64 `json:"C,string"`
	D      float64 `json:"D,string"`
}

// GetCandles fetches the candles from the server and returns array of
// CandlesResponse if the requested format is json, otherwise creates a csv file on disk, format is csv.
// Note : in case of csv format, CandlesRequest.FileNameWithPath is required
// and the []CandlesResponse return data is nil.
func (c *Connecter) GetCandles(cReq CandlesRequest) ([]CandlesResponse, error) {
	req, err := c.newRequest(candlesURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if cReq.Interval == "" {
		return nil, errors.New("interval is required")
	}
	q.Add("interval", cReq.Interval)
	if cReq.Currency == "" {
		return nil, errors.New("currency is required")
	}
	q.Add("currency", cReq.Currency)
	if !cReq.Start.IsZero() {
		q.Add("start", cReq.Start.Format(time.RFC3339))
	}
	if !cReq.End.IsZero() {
		q.Add("end", cReq.End.Format(time.RFC3339))
	}
	if cReq.Format != "" {
		if cReq.Format == "csv" && cReq.FileNameWithPath == "" {
			return nil, errors.New("csv file path is required")
		}
		q.Add("format", cReq.Format)
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response, if the requested format is json.
	if cReq.Format == "" || cReq.Format == "json" {
		var cResp []CandlesResponse
		if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
			return nil, err
		}
		return cResp, nil
	}

	// Creates a CSV file and then copies the server's response, if the requested format is csv.
	err = c.createFile(resp.Body, cReq.FileNameWithPath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Exchange Candles.

// ExchangeCandlesRequest represents exchange candles request parameters.
type ExchangeCandlesRequest struct {
	Interval string
	Exchange string
	Market   string
	Start    time.Time
	End      time.Time

	// If the format is csv, please also include fileNameWithPath field value.
	// Default is json.
	Format string

	// csv file name. Required, if the format is csv.
	// example : /home/user/nomicsdata/exchange_candles_1613046296.csv.
	// Here, new "exchange_candles_1613046296.csv" file will be created on existing, /home/user/nomicsdata/, directory.
	FileNameWithPath string
}

// ExchangeCandlesResponse represents exchange candles response.
// This will be nil, if the format requested is csv.
// Fields will contain default go lang values if there is no value received from the server.
type ExchangeCandlesResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Low       float64   `json:"low,string"`
	Open      float64   `json:"open,string"`
	Close     float64   `json:"close,string"`
	High      float64   `json:"high,string"`
	Volume    float64   `json:"volume,string"`
	NumTrades int       `json:"num_trades,string"`

	// If price_outlier from server is null, then this field contains false value instead of null.
	PriceOutlier bool `json:"price_outlier"`

	// If volume_outlier from server is null, then this field contains false value instead of null.
	VolumeOutlier bool `json:"volume_outlier"`
}

// GetExchangeCandles fetches the exchange candles from the server and returns array of
// ExchangeCandlesResponse if the requested format is json, otherwise creates a csv file on disk, format is csv.
// Note : in case of csv format, ExchangeCandlesRequest.FileNameWithPath is required
// and the []ExchangeCandlesResponse return data is nil.
func (c *Connecter) GetExchangeCandles(ecReq ExchangeCandlesRequest) ([]ExchangeCandlesResponse, error) {
	req, err := c.newRequest(exchangeCandlesURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if ecReq.Interval == "" {
		return nil, errors.New("interval is required")
	}
	q.Add("interval", ecReq.Interval)
	if ecReq.Exchange == "" {
		return nil, errors.New("exchange is required")
	}
	q.Add("exchange", ecReq.Exchange)
	if ecReq.Market == "" {
		return nil, errors.New("market is required")
	}
	q.Add("market", ecReq.Market)
	if !ecReq.Start.IsZero() {
		q.Add("start", ecReq.Start.Format(time.RFC3339))
	}
	if !ecReq.End.IsZero() {
		q.Add("end", ecReq.End.Format(time.RFC3339))
	}
	if ecReq.Format != "" {
		if ecReq.Format == "csv" && ecReq.FileNameWithPath == "" {
			return nil, errors.New("csv file path is required")
		}
		q.Add("format", ecReq.Format)
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response, if the requested format is json.
	if ecReq.Format == "" || ecReq.Format == "json" {
		var ecResp []ExchangeCandlesResponse
		if err := json.NewDecoder(resp.Body).Decode(&ecResp); err != nil {
			return nil, err
		}
		return ecResp, nil
	}

	// Creates a CSV file and then copies the server's response, if the requested format is csv.
	err = c.createFile(resp.Body, ecReq.FileNameWithPath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Markets Candles.

// MarketsCandlesRequest represents markets candles request parameters.
type MarketsCandlesRequest struct {
	Interval string
	Base     string
	Quote    string
	Start    time.Time
	End      time.Time

	// If the format is csv, please also include fileNameWithPath field value.
	// Default is json.
	Format string

	// csv file name. Required, if the format is csv.
	// example : /home/user/nomicsdata/markets_candles_1613046296.csv.
	// Here, new "markets_candles_1613046296.csv" file will be created on existing, /home/user/nomicsdata/, directory.
	FileNameWithPath string
}

// MarketsCandlesResponse represents markets candles response.
// This will be nil, if the format requested is csv.
// Fields will contain default go lang values if there is no value received from the server.
type MarketsCandlesResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Low       float64   `json:"low,string"`
	Open      float64   `json:"open,string"`
	Close     float64   `json:"close,string"`
	High      float64   `json:"high,string"`
	Volume    float64   `json:"volume,string"`
	NumTrades int       `json:"num_trades,string"`

	// If price_outlier from server is null, then this field contains false value instead of null.
	PriceOutlier bool `json:"price_outlier"`

	// If volume_outlier from server is null, then this field contains false value instead of null.
	VolumeOutlier bool `json:"volume_outlier"`
}

// GetMarketsCandles fetches the TestGetExchangeCandlesexchange candles from the server and returns array of
// MarketsCandlesResponse if the requested format is json, otherwise creates a csv file on disk, format is csv.
// Note : in case of csv format, MarketsCandlesRequest.FileNameWithPath is required
// and the []MarketsCandlesResponse return data is nil.
func (c *Connecter) GetMarketsCandles(mcReq MarketsCandlesRequest) ([]MarketsCandlesResponse, error) {
	req, err := c.newRequest(marketsCandlesURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if mcReq.Interval == "" {
		return nil, errors.New("interval is required")
	}
	q.Add("interval", mcReq.Interval)
	if mcReq.Base == "" {
		return nil, errors.New("base is required")
	}
	q.Add("base", mcReq.Base)
	if mcReq.Quote == "" {
		return nil, errors.New("quote is required")
	}
	q.Add("quote", mcReq.Quote)
	if !mcReq.Start.IsZero() {
		q.Add("start", mcReq.Start.Format(time.RFC3339))
	}
	if !mcReq.End.IsZero() {
		q.Add("end", mcReq.End.Format(time.RFC3339))
	}
	if mcReq.Format != "" {
		if mcReq.Format == "csv" && mcReq.FileNameWithPath == "" {
			return nil, errors.New("csv file path is required")
		}
		q.Add("format", mcReq.Format)
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response, if the requested format is json.
	if mcReq.Format == "" || mcReq.Format == "json" {
		var mcResp []MarketsCandlesResponse
		if err := json.NewDecoder(resp.Body).Decode(&mcResp); err != nil {
			return nil, err
		}
		return mcResp, nil
	}

	// Creates a CSV file and then copies the server's response, if the requested format is csv.
	err = c.createFile(resp.Body, mcReq.FileNameWithPath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
