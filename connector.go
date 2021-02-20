package gonomics

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// API endpoints for Nomics.
const (
	apiServerURL string = "https://api.nomics.com/v1"

	// Currencies.

	// Currencies Ticker (Partial Paid Plan).
	currenciesTickerURL string = apiServerURL + "/currencies/ticker"
	// Currencies Metadata.
	currenciesMetadataURL string = apiServerURL + "/currencies"
	// Currencies Sparkline.
	currenciesSparklineURL string = apiServerURL + "/currencies/sparkline"
	// Currencies Supply History (Paid Plan).
	currenciesSupplyHistoryURL string = apiServerURL + "/supplies/history"

	// Markets.

	// Markets.
	marketsURL string = apiServerURL + "/markets"
	// Markets Cap History (Partial Paid Plan).
	marketsCapHistoryURL string = apiServerURL + "/market-cap/history"
	// Exchange Markets Ticker (Paid Plan).
	exchangeMarketsTickerURL string = apiServerURL + "/exchange-markets/ticker"

	// Volume.

	// Volume History (Partial Paid Plan).
	volumeHistoryURL string = apiServerURL + "/volume/history"

	// Exchange Rates.

	// Exchange Rates.
	exchangeRatesURL string = apiServerURL + "/exchange-rates"
	// Exchange Rates History.
	exchangeRatesHistoryURL string = apiServerURL + "/exchange-rates/history"

	// Global.

	// Global Ticker (Paid Plan).
	globalTickerURL string = apiServerURL + "/global-ticker"

	// Exchanges.

	// Exchanges Ticker (Paid Plan).
	exchangesTickerURL string = apiServerURL + "/exchanges/ticker"
	// Exchanges Volume History (Paid Plan).
	exchangesVolumeHistoryURL string = apiServerURL + "/exchanges/volume/history"
	// Exchanges Metadata (Paid Plan).
	exchangesMetadataURL string = apiServerURL + "/exchanges"

	// Candles.

	// Candles (Paid Plan).
	candlesURL string = apiServerURL + "/candles"
	// Exchange Candles (Paid Plan).
	exchangeCandlesURL string = apiServerURL + "/exchange_candles"
	// Markets Candles (Paid Plan).
	marketsCandlesURL string = apiServerURL + "/markets/candles"

	// Trades.

	// Trades (Paid Plan).
	tradesURL string = apiServerURL + "/trades"

	// Orders.

	// Orders Snapshot (Paid Plan).
	ordersSnapshotURL string = apiServerURL + "/orders/snapshot"

	// Predictions.

	// Currencies Predictions Ticker (Paid Plan).
	currenciesPredictionsTickerURL string = apiServerURL + "/currencies/predictions/ticker"
	// Currencies Predictions History (Paid Plan).
	currenciesPredictionsHistoryURL string = apiServerURL + "/currencies/predictions/history"

	// Only used for Unit Testing.
	// Please check this latest demo key published in nomics doc or use private key for paid API endpoint testing.
	demoAPIKey = "demo-6410726746980cead2a17c9db9ef29af"
)

// Connecter to connect nomics server.
type Connecter struct {
	apiKey string

	HTTPClient *http.Client
}

// New creates a brand new Nomics Connector.
// Modify Connector http client to specific needs, like Timeout, MaxIdleConns etc. once this function returns.
func New(apiKey string) *Connecter {
	connector := &Connecter{
		apiKey:     apiKey,
		HTTPClient: &http.Client{},
	}
	return connector
}

// Helper methods

// newRequest creates net.http request.
func (c *Connecter) newRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// do makes the net.http request to the server.
func (c *Connecter) do(req *http.Request) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Check for user or server error
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("User or Server error. Please check. Status Code : %v, Status : %v", resp.StatusCode, resp.Status)
	}

	return resp, nil
}

// createFile copies http response body to a new csv file on disk.
func (c *Connecter) createFile(data io.ReadCloser, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, data)
	if err != nil {
		return err
	}
	return nil
}
