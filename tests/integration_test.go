package tests

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/tidwall/gjson"
)

var (
	mulderAddr string
)

func init() {
	flag.StringVar(&mulderAddr, "addr", "localhost:8080", "mulder host:port on which to connect to run the integration tests")
}

func TestRandomQuote(t *testing.T) {
	randomQuoteURL := fmt.Sprintf("http://%s/quote/random", mulderAddr)
	t.Logf("Testing %s", randomQuoteURL)

	resp, err := http.Get(randomQuoteURL)
	if err != nil {
		t.Fatalf("Got unexpected error on %s: %v", randomQuoteURL, err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Got wrong HTTP Status Code %d - expected %d", resp.StatusCode, http.StatusOK)
	}
	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Got wrong Content-Type '%s' - expected '%s'", contentType, "application/json")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Got unexpected error while reading the HTTP response body from %s: %v", randomQuoteURL, err)
	}

	quote := gjson.ParseBytes(body).Get("quote").String()
	t.Logf("Got quote: %s", quote)
	if len(quote) == 0 {
		t.Error("Got invalid empty quote")
	}
}
