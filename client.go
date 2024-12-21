package dexscreener

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const DefaultApiBaseURL = "https://api.dexscreener.com"

type Client struct {
	BaseURL    string
	httpClient *http.Client
}

type Option func(*Client)

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

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.BaseURL = baseURL
	}
}

func (c *Client) newRequest(method string, path string) (*http.Request, error) {
	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

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
	Url          string        `json:"url"`
	ChainId      string        `json:"chainId"`
	TokenAddress string        `json:"tokenAddress"`
	Icon         string        `json:"icon"`
	Header       string        `json:"header"`
	Description  string        `json:"description"`
	Links        []ProfileLink `json:"links"`
}

type ProfileLink struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	Url   string `json:"url"`
}

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

// Retrieves boosted tokens from the given endpoint.
func (c *Client) getBoostedTokens(endpoint string) ([]TokenBoosted, error) {
	req, err := c.newRequest(http.MethodGet, endpoint)
	if err != nil {
		return nil, err
	}

	var boosts []TokenBoosted
	err = c.doRequest(req, &boosts)
	return boosts, err
}

// Retrieves the latest boosted tokens.
func (c *Client) GetLatestBoostedTokens() ([]TokenBoosted, error) {
	return c.getBoostedTokens("/token-boosts/latest/v1")
}

// Retrieves the most active boosted tokens.
func (c *Client) GetMostActiveBoostedTokens() ([]TokenBoosted, error) {
	return c.getBoostedTokens("/token-boosts/top/v1")
}

type TokenOrder struct {
	Type             string `json:"type"`
	Status           string `json:"status"`
	PaymentTimestamp uint64 `json:"paymentTimestamp"`
}

// Retrieves token orders for a specific chain and token address.
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

type TokenInfo struct {
	Address string `json:"address"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
}

type Liquidity struct {
	Usd   float64 `json:"usd"`
	Base  float64 `json:"base"`
	Quote float64 `json:"quote"`
}

type Social struct {
	Platform string `json:"platform"`
	Handle   string `json:"handle"`
}

type Website struct {
	URL string `json:"url"`
}

type Info struct {
	ImageURL string    `json:"imageUrl"`
	Websites []Website `json:"websites"`
	Socials  []Social  `json:"socials"`
}

type Boosts struct {
	Active int `json:"active"`
}

type TimeframeStats struct {
	H24 float64 `json:"h24"`
	H6  float64 `json:"h6"`
	H1  float64 `json:"h1"`
	M5  float64 `json:"m5"`
}

type TimeframeTransactions struct {
	M5  TradeCounts `json:"m5"`
	H6  TradeCounts `json:"h6"`
	H1  TradeCounts `json:"h1"`
	H24 TradeCounts `json:"h24"`
}

type TradeCounts struct {
	Buys  uint64 `json:"buys"`
	Sells uint64 `json:"sells"`
}

type TokenPair struct {
	SchemaVersion string                `json:"schemaVersion"`
	ChainID       string                `json:"chainId"`
	DexID         string                `json:"dexId"`
	URL           string                `json:"url"`
	PairAddress   string                `json:"pairAddress"`
	Labels        []string              `json:"labels"`
	BaseToken     TokenInfo             `json:"baseToken"`
	QuoteToken    TokenInfo             `json:"quoteToken"`
	PriceNative   string                `json:"priceNative"`
	PriceUSD      string                `json:"priceUsd"`
	Liquidity     Liquidity             `json:"liquidity"`
	FDV           float64               `json:"fdv"`
	MarketCap     float64               `json:"marketCap"`
	PairCreatedAt int64                 `json:"pairCreatedAt"`
	Info          Info                  `json:"info"`
	Boosts        Boosts                `json:"boosts"`
	Volume        TimeframeStats        `json:"volume"`
	PriceChange   TimeframeStats        `json:"priceChange"`
	Transactions  TimeframeTransactions `json:"txns"`
}

type TokenPairs struct {
	SchemaVersion string      `json:"schemaVersion"`
	Pair          TokenPair   `json:"pair"`
	Pairs         []TokenPair `json:"pairs"`
}

// Retrieves token pairs by chain ID and pair ID.
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

// Retrieves token pairs for a token address.
func (c *Client) GetTokenPairs(tokenAddress string) (*TokenPairs, error) {
	req, err := c.newRequest(http.MethodGet, "/latest/dex/tokens/"+tokenAddress)
	if err != nil {
		return nil, err
	}

	var pairs TokenPairs
	err = c.doRequest(req, &pairs)
	return &pairs, err
}

// Search for pairs matching the query
func (c *Client) SearchPairs(query string) (*TokenPairs, error) {
	req, err := c.newRequest(http.MethodGet, "/latest/dex/search?q="+query)
	if err != nil {
		return nil, err
	}

	var pairs TokenPairs
	err = c.doRequest(req, &pairs)
	return &pairs, err
}
