package gonomics

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

// Currencies Ticker.

// CurrenciesTickerRequest represents currencies ticker request parameters.
type CurrenciesTickerRequest struct {
	Ids                 []string
	Interval            []string
	Convert             string
	Status              string
	Filter              string
	Sort                string
	IncludeTransparency bool
	PerPage             int
	Page                int
}

// CurrenciesTickerResponse represents currencies ticker response.
// Fields will contain default go lang values if there is no value received from the server.
type CurrenciesTickerResponse struct {
	ID                   string                           `json:"id"`
	Status               string                           `json:"status"`
	Price                float64                          `json:"price,string"`
	PriceDate            time.Time                        `json:"price_date"`
	PriceTimestamp       time.Time                        `json:"price_timestamp"`
	Symbol               string                           `json:"symbol"`
	CirculatingSupply    float64                          `json:"circulating_supply,string"`
	MaxSupply            float64                          `json:"max_supply,string"`
	Name                 string                           `json:"name"`
	LogoURL              string                           `json:"logo_url"`
	MarketCap            float64                          `json:"market_cap,string"`
	TransparentMarketCap float64                          `json:"transparent_market_cap,string"`
	NumExchanges         int                              `json:"num_exchanges,string"`
	NumPairs             int                              `json:"num_pairs,string"`
	NumPairsUnmapped     int                              `json:"num_pairs_unmapped,string"`
	FirstCandle          time.Time                        `json:"first_candle"`
	FirstTrade           time.Time                        `json:"first_trade"`
	FirstOrderBook       time.Time                        `json:"first_order_book"`
	FirstPricedAt        time.Time                        `json:"first_priced_at"`
	Rank                 int                              `json:"rank,string"`
	RankDelta            int                              `json:"rank_delta,string"`
	High                 float64                          `json:"high,string"`
	HighTimestamp        time.Time                        `json:"high_timestamp"`
	OneH                 CurrenciesTickerIntervalResponse `json:"1h"`
	OneD                 CurrenciesTickerIntervalResponse `json:"1d"`
	SevenD               CurrenciesTickerIntervalResponse `json:"7d"`
	Three0D              CurrenciesTickerIntervalResponse `json:"30d"`
	Three65D             CurrenciesTickerIntervalResponse `json:"365d"`
	Ytd                  CurrenciesTickerIntervalResponse `json:"ytd"`
}

// CurrenciesTickerIntervalResponse represents currencies ticker interval response, Included in CurrenciesTickerResponse.
// Fields will contain default go lang values if there is no value received from the server.
type CurrenciesTickerIntervalResponse struct {
	PriceChange                   float64                                              `json:"price_change,string"`
	PriceChangePct                float64                                              `json:"price_change_pct,string"`
	Volume                        float64                                              `json:"volume,string"`
	VolumeChange                  float64                                              `json:"volume_change,string"`
	VolumeChangePct               float64                                              `json:"volume_change_pct,string"`
	MarketCapChange               float64                                              `json:"market_cap_change,string"`
	MarketCapChangePct            float64                                              `json:"market_cap_change_pct,string"`
	TransparentMarketCapChange    float64                                              `json:"transparent_market_cap_change,string"`
	TransparentMarketCapChangePct float64                                              `json:"transparent_market_cap_change_pct,string"`
	VolumeTransparency            []CurrenciesTickerIntervalVolumeTransparencyResponse `json:"volume_transparency"`
	VolumeTransparencyGrade       string                                               `json:"volume_transparency_grade"`
}

// CurrenciesTickerIntervalVolumeTransparencyResponse represents currencies ticker interval volume transparency response,
// Included in CurrenciesTickerIntervalResponse.
// Fields will contain default go lang values if there is no value received from the server.
type CurrenciesTickerIntervalVolumeTransparencyResponse struct {
	Grade           string  `json:"grade"`
	Volume          float64 `json:"volume,string"`
	VolumeChange    float64 `json:"volume_change,string"`
	VolumeChangePct float64 `json:"volume_change_pct,string"`
}

// GetCurrenciesTicker fetches the currency ticker from the server and returns array of CurrenciesTickerResponse.
func (c *Connecter) GetCurrenciesTicker(ctReq CurrenciesTickerRequest) ([]CurrenciesTickerResponse, error) {
	req, err := c.newRequest(currenciesTickerURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if len(ctReq.Ids) > 0 {
		q.Add("ids", strings.Join(ctReq.Ids[:], ","))
	}
	if len(ctReq.Interval) > 0 {
		q.Add("interval", strings.Join(ctReq.Interval[:], ","))
	}
	if ctReq.Convert != "" {
		q.Add("convert", ctReq.Convert)
	}
	if ctReq.Status != "" {
		q.Add("status", ctReq.Status)
	}
	if ctReq.Filter != "" {
		q.Add("filter", ctReq.Filter)
	}
	if ctReq.Sort != "" {
		q.Add("sort", ctReq.Sort)
	}
	if ctReq.IncludeTransparency {
		q.Add("include-transparency", strconv.FormatBool(ctReq.IncludeTransparency))
	}
	if ctReq.PerPage != 0 {
		q.Add("per-page", strconv.Itoa(ctReq.PerPage))
	}
	if ctReq.Page != 0 {
		q.Add("page", strconv.Itoa(ctReq.Page))
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response.
	var ctResp []CurrenciesTickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&ctResp); err != nil {
		return nil, err
	}
	return ctResp, nil
}

// Currencies Metadata.

// CurrenciesMetadataRequest represents currencies metadata request parameters.
type CurrenciesMetadataRequest struct {
	Ids        []string
	Attributes []string

	// If the format is csv, please also include fileNameWithPath field value.
	// Default is json.
	Format string

	// csv file name. Required, if the format is csv.
	// example : /home/user/nomicsdata/currency_metadata_1613046296.csv.
	// Here, new "currency_metadata_1613046296.csv" file will be created on existing, /home/user/nomicsdata/, directory.
	FileNameWithPath string
}

// CurrenciesMetadataResponse represents currencies metadata response.
// This will be nil, if the format requested is csv.
// Fields will contain default go lang values if there is no value received from the server.
type CurrenciesMetadataResponse struct {
	ID                      string `json:"id"`
	OriginalSymbol          string `json:"original_symbol"`
	Name                    string `json:"name"`
	Description             string `json:"description"`
	WebsiteURL              string `json:"website_url"`
	LogoURL                 string `json:"logo_url"`
	BlogURL                 string `json:"blog_url"`
	DiscordURL              string `json:"discord_url"`
	FacebookURL             string `json:"facebook_url"`
	GithubURL               string `json:"github_url"`
	MediumURL               string `json:"medium_url"`
	RedditURL               string `json:"reddit_url"`
	TelegramURL             string `json:"telegram_url"`
	TwitterURL              string `json:"twitter_url"`
	WhitepaperURL           string `json:"whitepaper_url"`
	YoutubeURL              string `json:"youtube_url"`
	LinkedinURL             string `json:"linkedin_url"`
	BitcointalkURL          string `json:"bitcointalk_url"`
	BlockExplorerURL        string `json:"block_explorer_url"`
	ReplacedBy              string `json:"replaced_by"`
	CryptocontrolCoinID     string `json:"cryptocontrol_coin_id"`
	PlatformCurrencyID      string `json:"platform_currency_id"`
	PlatformContractAddress string `json:"platform_contract_address"`
}

// GetCurrenciesMetadata fetches the currency metadata from the server and returns array of
// CurrenciesMetadataResponse if the requested format is json, otherwise creates a csv file on disk, format is csv.
// Note : in case of csv format, CurrenciesMetadataRequest.FileNameWithPath is required
// and the []CurrenciesMetadataResponse return data is nil.
func (c *Connecter) GetCurrenciesMetadata(cmReq CurrenciesMetadataRequest) ([]CurrenciesMetadataResponse, error) {
	req, err := c.newRequest(currenciesMetadataURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if len(cmReq.Ids) > 0 {
		q.Add("ids", strings.Join(cmReq.Ids[:], ","))
	}
	if len(cmReq.Attributes) > 0 {
		q.Add("attributes", strings.Join(cmReq.Attributes[:], ","))
	}
	if cmReq.Format != "" {
		if cmReq.Format == "csv" && cmReq.FileNameWithPath == "" {
			return nil, errors.New("csv file path is required")
		}
		q.Add("format", cmReq.Format)
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response, if the requested format is json.
	if cmReq.Format == "" || cmReq.Format == "json" {
		var cmResp []CurrenciesMetadataResponse
		if err := json.NewDecoder(resp.Body).Decode(&cmResp); err != nil {
			return nil, err
		}
		return cmResp, nil
	}

	// Creates a CSV file and then copies the server's response, if the requested format is csv.
	err = c.createFile(resp.Body, cmReq.FileNameWithPath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Currencies Sparkline.

// CurrenciesSparklineRequest represents currencies sparkline request parameters.
type CurrenciesSparklineRequest struct {
	Ids     []string
	Start   time.Time
	End     time.Time
	Convert string
}

// CurrenciesSparklineResponse represents currencies sparkline response.
// Fields will contain default go lang values if there is no value received from the server.
type CurrenciesSparklineResponse struct {
	Currency   string
	Timestamps []time.Time
	Prices     []float64
}

// UnmarshalJSON is needed because nomics server gives []string of prices, but we need []float64 of prices.
func (cs *CurrenciesSparklineResponse) UnmarshalJSON(data []byte) error {
	temp := struct {
		Currency   string      `json:"currency"`
		Timestamps []time.Time `json:"timestamps"`
		Prices     []string    `json:"prices"`
	}{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	cs.Currency = temp.Currency
	cs.Timestamps = temp.Timestamps
	for _, v := range temp.Prices {
		if price, err := strconv.ParseFloat(v, 64); err == nil {
			cs.Prices = append(cs.Prices, price)
		}
	}
	return nil
}

// GetCurrenciesSparkline fetches the currency sparklines from the server and returns array of CurrenciesSparklineResponse.
func (c *Connecter) GetCurrenciesSparkline(csReq CurrenciesSparklineRequest) ([]CurrenciesSparklineResponse, error) {
	req, err := c.newRequest(currenciesSparklineURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if len(csReq.Ids) > 0 {
		q.Add("ids", strings.Join(csReq.Ids[:], ","))
	}
	if csReq.Start.IsZero() {
		return nil, errors.New("start is required")
	}
	q.Add("start", csReq.Start.Format(time.RFC3339))
	if !csReq.End.IsZero() {
		q.Add("end", csReq.End.Format(time.RFC3339))
	}
	if csReq.Convert != "" {
		q.Add("convert", csReq.Convert)
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response.
	var csResp []CurrenciesSparklineResponse
	if err := json.NewDecoder(resp.Body).Decode(&csResp); err != nil {
		return nil, err
	}
	return csResp, nil
}

// Currencies Supply History.

// CurrenciesSupplyHistoryRequest represents currencies supply history request parameters.
type CurrenciesSupplyHistoryRequest struct {
	Currency string
	Start    time.Time
	End      time.Time

	// If the format is csv, please also include fileNameWithPath field value.
	// Default is json.
	Format string

	// csv file name. Required, if the format is csv.
	// example : /home/user/nomicsdata/currency_supplyhistory_1613046296.csv.
	// Here, new "currency_supplyhistory_1613046296.csv" file will be created on existing, /home/user/nomicsdata/, directory.
	FileNameWithPath string
}

// CurrenciesSupplyHistoryResponse represents currencies supply history response.
// This will be nil, if the format requested is csv.
// Fields will contain default go lang values if there is no value received from the server.
type CurrenciesSupplyHistoryResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Available float64   `json:"available,string"`
	Max       float64   `json:"max,string"`
}

// GetCurrenciesSupplyHistory fetches the currency supply history from the server and returns array of
// CurrenciesSupplyHistoryResponse if the requested format is json, otherwise creates a csv file on disk, format is csv.
// Note : in case of csv format, CurrenciesSupplyHistoryRequest.FileNameWithPath is required
// and the []CurrenciesSupplyHistoryResponse return data is nil.
func (c *Connecter) GetCurrenciesSupplyHistory(cshReq CurrenciesSupplyHistoryRequest) ([]CurrenciesSupplyHistoryResponse, error) {
	req, err := c.newRequest(currenciesSupplyHistoryURL)
	if err != nil {
		return nil, err
	}

	// Formulate query params.
	q := req.URL.Query()
	q.Add("key", c.apiKey)
	if cshReq.Currency == "" {
		return nil, errors.New("currency is required")
	}
	q.Add("currency", cshReq.Currency)
	if cshReq.Start.IsZero() {
		return nil, errors.New("start is required")
	}
	q.Add("start", cshReq.Start.Format(time.RFC3339))
	if !cshReq.End.IsZero() {
		q.Add("end", cshReq.End.Format(time.RFC3339))
	}
	if cshReq.Format != "" {
		if cshReq.Format == "csv" && cshReq.FileNameWithPath == "" {
			return nil, errors.New("csv file path is required")
		}
		q.Add("format", cshReq.Format)
	}
	req.URL.RawQuery = q.Encode()

	// Do the requset to server.
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Creates formatted response from the server's response, if the requested format is json.
	if cshReq.Format == "" || cshReq.Format == "json" {
		var cshResp []CurrenciesSupplyHistoryResponse
		if err := json.NewDecoder(resp.Body).Decode(&cshResp); err != nil {
			return nil, err
		}
		return cshResp, nil
	}

	// Creates a CSV file and then copies the server's response, if the requested format is csv.
	err = c.createFile(resp.Body, cshReq.FileNameWithPath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
