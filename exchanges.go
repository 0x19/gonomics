package gonomics

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

// Exchanges Ticker.

// ExchangesTickerRequest represents exchanges ticker request parameters.
type ExchangesTickerRequest struct {
	Ids      []string
	Interval []string
	Convert  string
	Status   string
	Type     string
	PerPage  int
	Page     int
}

// ExchangesTickerResponse represents exchanges ticker response.
// Fields will contain default go lang values if there is no value received from the server.
type ExchangesTickerResponse struct {
	ID                string                          `json:"id"`
	Name              string                          `json:"name"`
	LogoURL           string                          `json:"logo_url"`
	Rank              int                             `json:"rank,string"`
	TransparencyGrade string                          `json:"transparency_grade"`
	CoverageType      string                          `json:"coverage_type"`
	OrderBooks        bool                            `json:"order_books,string"`
	FirstTrade        time.Time                       `json:"first_trade"`
	FirstCandle       time.Time                       `json:"first_candle"`
	FirstOrderBook    time.Time                       `json:"first_order_book"`
	LastUpdated       time.Time                       `json:"last_updated"`
	FiatCurrencies    []string                        `json:"fiat_currencies"`
	NumPairs          int                             `json:"num_pairs,string"`
	NumPairsUnmapped  int                             `json:"num_pairs_unmapped,string"`
	OneH              ExchangesTickerIntervalResponse `json:"1h"`
	OneD              ExchangesTickerIntervalResponse `json:"1d"`
	SevenD            ExchangesTickerIntervalResponse `json:"7d"`
	Three0D           ExchangesTickerIntervalResponse `json:"30d"`
	Three65D          ExchangesTickerIntervalResponse `json:"365d"`
	Ytd               ExchangesTickerIntervalResponse `json:"ytd"`
}

// ExchangesTickerIntervalResponse represents exchanges ticker interval response,
// Included in ExchangesTickerResponse.
// Fields will contain default go lang values if there is no value received from the server.
type ExchangesTickerIntervalResponse struct {
	Volume                    float64 `json:"volume,string"`
	VolumeChange              float64 `json:"volume_change,string"`
	VolumeChangePct           float64 `json:"volume_change_pct,string"`
	SpotVolume                float64 `json:"spot_volume,string"`
	SpotVolumeChange          float64 `json:"spot_volume_change,string"`
	SpotVolumeChangePct       float64 `json:"spot_volume_change_pct,string"`
	DerivativeVolume          float64 `json:"derivative_volume,string"`
	DerivativeVolumeChange    float64 `json:"derivative_volume_change,string"`
	DerivativeVolumeChangePct float64 `json:"derivative_volume_change_pct,string"`
	Trades                    int     `json:"trades,string"`
	TradesChange              int     `json:"trades_change,string"`
	TradesChangePct           float64 `json:"trades_change_pct,string"`
}

// GetExchangesTicker fetches the exchanges ticker from the server and returns array of ExchangesTickerResponse.
func (c *Connecter) GetExchangesTicker(etReq ExchangesTickerRequest) ([]ExchangesTickerResponse, error) {
	req, err := c.newRequest(exchangesTickerURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if len(etReq.Ids) > 0 {
		q.Add("ids", strings.Join(etReq.Ids[:], ","))
	}
	if len(etReq.Interval) > 0 {
		q.Add("interval", strings.Join(etReq.Interval[:], ","))
	}
	if etReq.Convert != "" {
		q.Add("convert", etReq.Convert)
	}
	if etReq.Status != "" {
		q.Add("status", etReq.Status)
	}
	if etReq.Type != "" {
		q.Add("type", etReq.Type)
	}
	if etReq.PerPage != 0 {
		q.Add("per-page", strconv.Itoa(etReq.PerPage))
	}
	if etReq.Page != 0 {
		q.Add("page", strconv.Itoa(etReq.Page))
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response.
	var etResp []ExchangesTickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&etResp); err != nil {
		return nil, err
	}
	return etResp, nil
}

// Exchanges Volume History.

// ExchangesVolumeHistoryRequest represents exchanges volume history request parameters.
type ExchangesVolumeHistoryRequest struct {
	Exchange string
	Start    time.Time
	End      time.Time
	Convert  string

	// If the format is csv, please also include fileNameWithPath field value.
	// Default is json.
	Format string

	// csv file name. Required, if the format is csv.
	// example : /home/user/nomicsdata/exchanges_volume_history_1613046296.csv.
	// Here, new "exchanges_volume_history_1613046296.csv" file will be created on existing, /home/user/nomicsdata/, directory.
	FileNameWithPath string

	IncludeTransparency bool
}

// ExchangesVolumeHistoryResponse represents exchanges volume history response.
// This will be nil, if the format requested is csv.
// Fields will contain default go lang values if there is no value received from the server.
type ExchangesVolumeHistoryResponse struct {
	Timestamp         time.Time `json:"timestamp"`
	Volume            float64   `json:"volume,string"`
	TransparentVolume float64   `json:"transparent_volume,string"`
}

// GetExchangesVolumeHistory fetches the exchanges volume history from the server and returns array of
// ExchangesVolumeHistoryResponse if the requested format is json, otherwise creates a csv file on disk, format is csv.
// Note : in case of csv format, ExchangesVolumeHistoryRequest.FileNameWithPath is required
// and the []ExchangesVolumeHistoryResponse return data is nil.
func (c *Connecter) GetExchangesVolumeHistory(evhReq ExchangesVolumeHistoryRequest) ([]ExchangesVolumeHistoryResponse, error) {
	req, err := c.newRequest(exchangesVolumeHistoryURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if evhReq.Exchange == "" {
		return nil, errors.New("exchange is required")
	}
	q.Add("exchange", evhReq.Exchange)
	if evhReq.Start.IsZero() {
		return nil, errors.New("start is required")
	}
	q.Add("start", evhReq.Start.Format(time.RFC3339))
	if !evhReq.End.IsZero() {
		q.Add("end", evhReq.End.Format(time.RFC3339))
	}
	if evhReq.Convert != "" {
		q.Add("convert", evhReq.Convert)
	}
	if evhReq.Format != "" {
		if evhReq.Format == "csv" && evhReq.FileNameWithPath == "" {
			return nil, errors.New("csv file path is required")
		}
		q.Add("format", evhReq.Format)
	}
	if evhReq.IncludeTransparency {
		q.Add("include-transparency", strconv.FormatBool(evhReq.IncludeTransparency))
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response, if the requested format is json.
	if evhReq.Format == "" || evhReq.Format == "json" {
		var evhResp []ExchangesVolumeHistoryResponse
		if err := json.NewDecoder(resp.Body).Decode(&evhResp); err != nil {
			return nil, err
		}
		return evhResp, nil
	}

	// Creates a CSV file and then copies the server's response, if the requested format is csv.
	err = c.createFile(resp.Body, evhReq.FileNameWithPath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Exchanges Metadata.

// ExchangesMetadataRequest represents exchanges metadata request parameters.
type ExchangesMetadataRequest struct {
	Ids        []string
	Attributes []string

	// If the format is csv, please also include fileNameWithPath field value.
	// Default is json.
	Format string

	// csv file name. Required, if the format is csv.
	// example : /home/user/nomicsdata/exchanges_metadata_1613046296.csv.
	// Here, new "exchanges_metadata_1613046296.csv" file will be created on existing, /home/user/nomicsdata/, directory.
	FileNameWithPath string
}

// ExchangesMetadataResponse represents exchanges metadata response.
// This will be nil, if the format requested is csv.
// Fields will contain default go lang values if there is no value received from the server.
type ExchangesMetadataResponse struct {
	ID                          string `json:"id"`
	CapabilityMarkets           bool   `json:"capability_markets"`
	CapabilityTrades            bool   `json:"capability_trades"`
	CapabilityTradesByTimestamp bool   `json:"capability_trades_by_timestamp"`
	CapabilityTradesSnapshot    bool   `json:"capability_trades_snapshot"`
	CapabilityOrdersSnapshot    bool   `json:"capability_orders_snapshot"`
	CapabilityCandles           bool   `json:"capability_candles"`
	CapabilityTicker            bool   `json:"capability_ticker"`
	Integrated                  bool   `json:"integrated"`
	Name                        string `json:"name"`
	Description                 string `json:"description"`
	Location                    string `json:"location"`
	LogoURL                     string `json:"logo_url"`
	WebsiteURL                  string `json:"website_url"`
	FeesURL                     string `json:"fees_url"`
	TwitterURL                  string `json:"twitter_url"`
	FacebookURL                 string `json:"facebook_url"`
	RedditURL                   string `json:"reddit_url"`
	ChatURL                     string `json:"chat_url"`
	BlogURL                     string `json:"blog_url"`
	Year                        int    `json:"year"`
	TransparencyGrade           string `json:"transparency_grade"`
	OrderBooksInterval          int    `json:"order_books_interval"`
}

// GetExchangesMetadata fetches the exchanges metadata from the server and returns array of
// ExchangesMetadataResponse if the requested format is json, otherwise creates a csv file on disk, format is csv.
// Note : in case of csv format, ExchangesMetadataRequest.FileNameWithPath is required
// and the []ExchangesMetadataResponse return data is nil.
func (c *Connecter) GetExchangesMetadata(emReq ExchangesMetadataRequest) ([]ExchangesMetadataResponse, error) {
	req, err := c.newRequest(exchangesMetadataURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if len(emReq.Ids) > 0 {
		q.Add("ids", strings.Join(emReq.Ids[:], ","))
	}
	if len(emReq.Attributes) > 0 {
		q.Add("attributes", strings.Join(emReq.Attributes[:], ","))
	}
	if emReq.Format != "" {
		if emReq.Format == "csv" && emReq.FileNameWithPath == "" {
			return nil, errors.New("csv file path is required")
		}
		q.Add("format", emReq.Format)
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response, if the requested format is json.
	if emReq.Format == "" || emReq.Format == "json" {
		var emResp []ExchangesMetadataResponse
		if err := json.NewDecoder(resp.Body).Decode(&emResp); err != nil {
			return nil, err
		}
		return emResp, nil
	}

	// Creates a CSV file and then copies the server's response, if the requested format is csv.
	err = c.createFile(resp.Body, emReq.FileNameWithPath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
