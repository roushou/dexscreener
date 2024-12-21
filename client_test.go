package dexscreener_test

import (
	"testing"

	"github.com/roushou/dexscreener"
)

func TestWithBaseURL_Default(t *testing.T) {
	client := dexscreener.NewClient()

	if client == nil {
		t.Fatalf("Expected non-nil client, got nil")
	}
	if client.BaseURL != dexscreener.DefaultApiBaseURL {
		t.Errorf("Expected baseUrl to be %s, got %s", dexscreener.DefaultApiBaseURL, client.BaseURL)
	}
}

func TestWithBaseURL(t *testing.T) {
	customURL := "https://custom-api.example.com"
	client := dexscreener.NewClient(dexscreener.WithBaseURL(customURL))
	if client == nil {
		t.Fatalf("Expected non-nil client, got nil")
	}
	if client.BaseURL != customURL {
		t.Errorf("Expected baseUrl to be %s, got %s", customURL, client.BaseURL)
	}
}

func TestGetTokenProfiles(t *testing.T) {
	client := dexscreener.NewClient()

	profiles, err := client.GetTokenProfiles()
	if err != nil {
		t.Fatalf("Failed to fetch token profiles: %v", err)
	}

	if len(profiles) == 0 {
		t.Error("Expected at least one token profile, got none.")
	}
}

func TestGetLatestBoostedTokens(t *testing.T) {
	client := dexscreener.NewClient()

	tokens, err := client.GetLatestBoostedTokens()
	if err != nil {
		t.Fatalf("Failed to fetch latest boosted tokens: %v", err)
	}

	if len(tokens) == 0 {
		t.Fatal("Expected at least one boosted tokens, got none")
	}
}

func TestGetMostActiveBoostedTokens(t *testing.T) {
	client := dexscreener.NewClient()

	tokens, err := client.GetMostActiveBoostedTokens()
	if err != nil {
		t.Fatalf("Failed to fetch most active boosted tokens: %v", err)
	}

	if len(tokens) == 0 {
		t.Fatal("Expected at least one boosted tokens, got none")
	}
}

func TestGetTokenOrders(t *testing.T) {
	client := dexscreener.NewClient()

	orders, err := client.GetTokenOrders("solana", "A55XjvzRU4KtR3Lrys8PpLZQvPojPqvnv5bJVHMYy3Jv")
	if err != nil {
		t.Fatalf("Failed to fetch token orders: %v", err)
	}

	if len(orders) == 0 {
		t.Error("Expected at least one token order, got none.")
	}
}

func TestGetTokenPairsByChain(t *testing.T) {
	client := dexscreener.NewClient()

	pairs, err := client.GetTokenPairsByChain("base", "0xf1fdc83c3a336bdbdc9fb06e318b08eaddc82ff4")
	if err != nil {
		t.Fatalf("Failed to fetch token pairs: %v", err)
	}

	if pairs == nil {
		t.Fatal("Token pairs should not be nil")
	}
}

func TestGetTokenPairs(t *testing.T) {
	client := dexscreener.NewClient()

	pairs, err := client.GetTokenPairs("0x4F9Fd6Be4a90f2620860d680c0d4d5Fb53d1A825")
	if err != nil {
		t.Fatalf("Failed to fetch token pairs: %v", err)
	}

	if pairs == nil {
		t.Fatal("Token pairs should not be nil")
	}
}

func TestSearchPairs(t *testing.T) {
	client := dexscreener.NewClient()

	pairs, err := client.SearchPairs("SOL/USDC")
	if err != nil {
		t.Fatalf("Failed to search for token pairs: %v", err)
	}

	if pairs == nil {
		t.Fatal("Token pairs should not be nil")
	}
}
