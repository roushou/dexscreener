package dexscreener

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const DefaultApiBaseURL = "https://api.dexscreener.com"

// Client represents a client for interacting with the DexScreener API
type Client struct {
	BaseURL    string
	httpClient *http.Client
}

// Option defines a function that can customize the Client
type Option func(*Client)

// NewClient creates a new DexScreener client with default settings and optional configurations
func NewClient(options ...Option) *Client {
	client := &Client{
		BaseURL:    DefaultApiBaseURL,
		httpClient: &http.Client{},
	}

	for _, option := range options {
		option(client)
	}

	return client
}

// WithBaseURL returns an Option that sets the base URL for API requests
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.BaseURL = baseURL
	}
}

// newRequest constructs a new HTTP request for the DexScreener API
func (c *Client) newRequest(method string, path string) (*http.Request, error) {
	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// doRequest sends an HTTP request and decodes the response into the provided interface
func (c *Client) doRequest(req *http.Request, out interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed (status: %s, body: %s)", resp.Status, string(body))

	}
	return json.NewDecoder(resp.Body).Decode(out)

}

type TokenProfile struct {
	URL          string `json:"url"`
	ChainId      string `json:"chainId"`
	TokenAddress string `json:"tokenAddress"`
	Icon         string `json:"icon"`
	Header       string `json:"header"`
	Description  string `json:"description"`
	Links        []struct {
		Type  string `json:"type"`
		Label string `json:"label"`
		Url   string `json:"url"`
	} `json:"links"`
}

// GetTokenProfiles fetches the latest token profiles
func (c *Client) GetTokenProfiles() ([]TokenProfile, error) {
	req, err := c.newRequest(http.MethodGet, "/token-profiles/latest/v1")
	if err != nil {
		return nil, err
	}

	var profiles []TokenProfile
	err = c.doRequest(req, &profiles)
	return profiles, err
}

type TokenBoosted struct {
	TokenProfile
	Amount      uint64 `json:"amount"`
	TotalAmount uint64 `json:"totalAmount"`
}

// getBoostedTokens is a helper function to retrieve boosted tokens from a specific endpoint
func (c *Client) getBoostedTokens(endpoint string) ([]TokenBoosted, error) {
	req, err := c.newRequest(http.MethodGet, endpoint)
	if err != nil {
		return nil, err
	}

	var boosts []TokenBoosted
	err = c.doRequest(req, &boosts)
	return boosts, err
}

// GetLatestBoostedTokens fetches the most recently boosted tokens
func (c *Client) GetLatestBoostedTokens() ([]TokenBoosted, error) {
	return c.getBoostedTokens("/token-boosts/latest/v1")
}

// GetMostActiveBoostedTokens fetches the most active boosted tokens
func (c *Client) GetMostActiveBoostedTokens() ([]TokenBoosted, error) {
	return c.getBoostedTokens("/token-boosts/top/v1")
}

type TokenOrder struct {
	Type             string `json:"type"`
	Status           string `json:"status"`
	PaymentTimestamp uint64 `json:"paymentTimestamp"`
}

// GetTokenOrders fetches token orders for a specific chain and token address.
func (c *Client) GetTokenOrders(chainId string, tokenAddress string) ([]TokenOrder, error) {
	path := fmt.Sprintf("/orders/v1/%s/%s", chainId, tokenAddress)
	req, err := c.newRequest(http.MethodGet, path)
	if err != nil {
		return nil, err
	}

	var orders []TokenOrder
	err = c.doRequest(req, &orders)
	return orders, err
}

type TokenPair struct {
	SchemaVersion string   `json:"schemaVersion"`
	ChainID       string   `json:"chainId"`
	DexID         string   `json:"dexId"`
	URL           string   `json:"url"`
	PairAddress   string   `json:"pairAddress"`
	Labels        []string `json:"labels"`
	BaseToken     struct {
		Address string `json:"address"`
		Name    string `json:"name"`
		Symbol  string `json:"symbol"`
	} `json:"baseToken"`
	QuoteToken struct {
		Address string `json:"address"`
		Name    string `json:"name"`
		Symbol  string `json:"symbol"`
	} `json:"quoteToken"`
	PriceNative string `json:"priceNative"`
	PriceUSD    string `json:"priceUsd"`
	Liquidity   struct {
		USD   float64 `json:"usd"`
		Base  float64 `json:"base"`
		Quote float64 `json:"quote"`
	} `json:"liquidity"`
	FDV           float64 `json:"fdv"`
	MarketCap     float64 `json:"marketCap"`
	PairCreatedAt int64   `json:"pairCreatedAt"`
	Info          struct {
		ImageURL string `json:"imageUrl"`
		Websites []struct {
			URL string `json:"url"`
		} `json:"websites"`
		Socials []struct {
			Platform string `json:"platform"`
			Handle   string `json:"handle"`
		} `json:"socials"`
	} `json:"info"`
	Boosts struct {
		Active int `json:"active"`
	} `json:"boosts"`
	Volume struct {
		H24 float64 `json:"h24"`
		H6  float64 `json:"h6"`
		H1  float64 `json:"h1"`
		M5  float64 `json:"m5"`
	} `json:"volume"`
	PriceChange struct {
		H24 float64 `json:"h24"`
		H6  float64 `json:"h6"`
		H1  float64 `json:"h1"`
		M5  float64 `json:"m5"`
	} `json:"priceChange"`
	Transactions struct {
		M5 struct {
			Buys  uint64 `json:"buys"`
			Sells uint64 `json:"sells"`
		} `json:"m5"`
		H6 struct {
			Buys  uint64 `json:"buys"`
			Sells uint64 `json:"sells"`
		} `json:"h6"`
		H1 struct {
			Buys  uint64 `json:"buys"`
			Sells uint64 `json:"sells"`
		} `json:"h1"`
		H24 struct {
			Buys  uint64 `json:"buys"`
			Sells uint64 `json:"sells"`
		} `json:"h24"`
	} `json:"txns"`
}

type TokenPairs struct {
	SchemaVersion string      `json:"schemaVersion"`
	Pair          TokenPair   `json:"pair"`
	Pairs         []TokenPair `json:"pairs"`
}

// GetTokenPairsByChain fetches token pairs by chain ID and pair ID.
func (c *Client) GetTokenPairsByChain(chainId string, pairId string) (*TokenPairs, error) {
	path := fmt.Sprintf("/latest/dex/pairs/%s/%s", chainId, pairId)
	req, err := c.newRequest(http.MethodGet, path)
	if err != nil {
		return nil, err
	}

	var pairs TokenPairs
	err = c.doRequest(req, &pairs)
	return &pairs, err
}

// GetTokenPairs fetches token pairs for a token address.
func (c *Client) GetTokenPairs(tokenAddress string) (*TokenPairs, error) {
	req, err := c.newRequest(http.MethodGet, "/latest/dex/tokens/"+tokenAddress)
	if err != nil {
		return nil, err
	}

	var pairs TokenPairs
	err = c.doRequest(req, &pairs)
	return &pairs, err
}

// SearchPairs fetches token pairs matching the query
func (c *Client) SearchPairs(query string) (*TokenPairs, error) {
	req, err := c.newRequest(http.MethodGet, "/latest/dex/search?q="+query)
	if err != nil {
		return nil, err
	}

	var pairs TokenPairs
	err = c.doRequest(req, &pairs)
	return &pairs, err
}
