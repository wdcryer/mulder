package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	quotesKey = "quotes"
)

var (
	quotes = []interface{}{ // https://www.imdb.com/title/tt0106179/quotes
		"I would never lie. I willfully participated in a campaign of misinformation.",
		"Scully, I was like you once. I didn't know who to trust. Then I... I chose another path... another life, another fate, where I found my sister. The end of my world was unrecognizable and upside down. There was one thing that remained the same. You were my friend, and you told me the truth. Even when the world was falling apart, you were my constant. My touchstone.",
		"Trust no one.",
		"You know, they say when you talk to God it's prayer, but when God talks to you, it's schizophrenia.",
		"Sorry, nobody down here but the FBI's most unwanted.",
		"I have a theory. Do you want to hear it?",
		"You have to be willing to see.",
		"Scully, you have to believe me. Nobody else on this whole damn planet does or ever will. You're my one in five billion.",
		"Scully, you are the only one I trust.",
		"Sometimes the only sane answer to an insane world is insanity.",
		"I've often felt that dreams are answers to questions we haven't yet figured out how to ask.",
		"We've both lost so much... but I believe that what we're looking for is in the X-Files. I'm more certain than ever that the truth is in there.",
		"If coincidences are coincidences, why do they feel so contrived?",
		"And all the choices would then lead to this very moment. One wrong turn, and we wouldn't be sitting here together. Well, that says a lot. That says a lot, a lot, a lot.",
		"The truth will save you, Scully. I think it'll save both of us.",
		"THE TRUTH IS OUT THERE",
		"I want to believe.",
		"TRUST NO-ONE",
		"What can I do about a Lie with an Official Seal on it?",
	}

	listenAddr          string
	redisAddr           string
	redisConnectTimeout time.Duration

	redisConn redis.Conn
)

func init() {
	rand.Seed(time.Now().UnixNano())

	flag.StringVar(&listenAddr, "listen-addr", ":8080", "host:port on which to listen")
	flag.StringVar(&redisAddr, "redis-addr", ":6379", "redis host:port to connect to")
	flag.DurationVar(&redisConnectTimeout, "redis-connect-timeout", 1*time.Minute, "timeout for connecting to redis")

	http.HandleFunc("/quote/random", randomQuoteHandler)
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
}

func main() {
	log.Println("Mulder is waking up...")
	flag.Parse()

	if err := connectToRedis(); err != nil {
		log.Fatalf("Failed to connect to The (redis) X-Files at %s after timeout %s: %v", redisAddr, redisConnectTimeout, err)
	}
	defer redisConn.Close()

	if err := insertQuotesInRedis(); err != nil {
		log.Fatalf("Failed to insert files in The X-Files: %v", err)
	}

	log.Printf("Starting HTTP server on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("Failed to listen on %s: %v", listenAddr, err)
	}
}

func connectToRedis() (err error) {
	log.Printf("Connecting to The (redis) X-Files at %s...", redisAddr)

	redisConn, err = redis.Dial("tcp", redisAddr, redis.DialConnectTimeout(redisConnectTimeout))
	if err != nil {
		return err
	}

	infos, err := redis.String(redisConn.Do("INFO", "SERVER"))
	if err != nil {
		return err
	}

	log.Printf("Connected to The (redis) X-Files:\n%s", infos)
	return nil
}

func insertQuotesInRedis() error {
	log.Println("Checking The X-Files...")
	existingQuotes, err := redis.Int(redisConn.Do("LLEN", quotesKey))
	if err != nil {
		return err
	}

	if existingQuotes == len(quotes) {
		log.Printf("All The %d X-Files are already there!", existingQuotes)
		return nil
	}

	if existingQuotes > 0 {
		log.Printf("There is a mess in The X-Files, we don't have the right number of quotes - %d instead of %d. Let's clean everything first...", existingQuotes, len(quotes))
		if _, err = redis.Int(redisConn.Do("DEL", quotesKey)); err != nil {
			return err
		}
	}

	log.Printf("Inserting %d files in The X-Files...", len(quotes))
	args := append([]interface{}{}, quotesKey)
	args = append(args, quotes...)
	insertedQuotes, err := redis.Int(redisConn.Do("RPUSH", args...))
	if err != nil {
		return err
	}

	log.Printf("Inserted %d/%d files in The X-Files", insertedQuotes, len(quotes))
	return nil
}

func getRandomQuote() (string, error) {
	quotesCount, err := redis.Int(redisConn.Do("LLEN", quotesKey))
	if err != nil {
		return "", err
	}

	randomIndex := rand.Intn(quotesCount)
	return redis.String(redisConn.Do("LINDEX", quotesKey, randomIndex))
}

func randomQuoteHandler(w http.ResponseWriter, r *http.Request) {
	quote, err := getRandomQuote()
	if err != nil {
		log.Printf("Failed to retrieve an X-File: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Handled an X-File request, returned: '%s'", quote)

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(&response{Quote: quote}); err != nil {
		log.Printf("Failed to write HTTP response: %v", err)
	}
}

type response struct {
	Quote string `json:"quote"`
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	pong, err := redis.String(redisConn.Do("PING"))
	if err != nil {
		log.Printf("Healthz handler failing: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, pong)
}
