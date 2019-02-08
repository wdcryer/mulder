package main

import (
	"testing"

	"github.com/rafaeljusto/redigomock"
)

func TestGetRandomQuote(t *testing.T) {
	mock := redigomock.NewConn()
	redisConn = mock

	mock.Command("LLEN", quotesKey).Expect(int64(len(quotes)))
	for i, quote := range quotes {
		mock.Command("LINDEX", quotesKey, i).Expect(quote)
	}

	quote, err := getRandomQuote()
	if err != nil {
		t.Fatalf("Failed to get a random quote: %v", err)
	}
	for _, q := range quotes {
		if q.(string) == quote {
			return
		}
	}
	t.Errorf("Invalid random quote %s", quote)
}
