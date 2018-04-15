package main

import (
	log "github.com/sirupsen/logrus"

	"flag"
	"io"
	"net/http"
	"os"

	"github.com/zqureshi/go/redirect"
)

var (
	logger = &log.Logger{
		Out:       os.Stdout,
		Formatter: new(log.JSONFormatter),
		Hooks:     make(log.LevelHooks),
		Level:     log.InfoLevel,
	}

	listenAddr     string
	airtableAPIKey string
	airtableBaseID string
)

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "PONG")
}

func init() {
	flag.StringVar(&listenAddr, "listen", ":80", "[host]:port on which to bind server")
	flag.Parse()

	airtableAPIKey = os.Getenv("AIRTABLE_API_KEY")
	if airtableAPIKey == "" {
		logger.Fatal("AIRTABLE_API_KEY must be specified")
	}

	airtableBaseID = os.Getenv("AIRTABLE_BASE_ID")
	if airtableBaseID == "" {
		logger.Fatal("AIRTABLE_BASE_ID must be specified")
	}
}

func main() {
	c, err := redirect.NewClient(airtableAPIKey, airtableBaseID)
	if err != nil {
		logger.Fatal(err)
	}

	client, err := redirect.NewCachingClient(c)
	if err != nil {
		logger.Fatal(err)
	}

	defaultRedirect, err := client.Get(redirect.DefaultRedirectKey)
	if err != nil {
		logger.Fatal("A 'default' redirect must be specified")
	}
	logger.Info("Default redirect " + defaultRedirect.URL)

	http.HandleFunc("/_ah/health", healthHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if r, err := client.Get(req.URL.Path[1:]); req.URL.Path != "/" && err == nil {
			http.Redirect(w, req, r.URL, 302)
		} else {
			d, _ := client.Get(redirect.DefaultRedirectKey)
			http.Redirect(w, req, d.URL, 302)
		}
	})

	logger.Fatal(http.ListenAndServe(listenAddr, nil))
}
