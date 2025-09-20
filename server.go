package honeybadger

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Errors returned by the backend when unable to successfully handle payload.
var (
	ErrRateExceeded    = errors.New("Rate exceeded: slow down!")
	ErrPaymentRequired = errors.New("Payment required: expired trial or credit card?")
	ErrUnauthorized    = errors.New("Unauthorized: bad API key?")
)

func newServerBackend(config *Configuration) *server {
	return &server{
		URL:    &config.Endpoint,
		APIKey: &config.APIKey,
		Client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   config.Timeout,
		},
		Timeout: &config.Timeout,
	}
}

type server struct {
	APIKey  *string
	URL     *string
	Timeout *time.Duration
	Client  *http.Client
}

func (s *server) Notify(feature Feature, payload Payload) error {
	return s.sendRequest("v1/"+feature.Endpoint, payload.toJSON(), "application/json")
}

func (s *server) Event(events []*eventPayload) error {
	var jsonl []byte
	for _, event := range events {
		jsonl = append(jsonl, event.toJSON()...)
		jsonl = append(jsonl, '\n')
	}
	return s.sendRequest("v1/events", jsonl, "application/x-ndjson")
}

func (s *server) sendRequest(path string, body []byte, contentType string) error {
	s.Client.Timeout = *s.Timeout

	url, err := url.Parse(*s.URL)
	if err != nil {
		return err
	}
	url.Path = path

	req, err := http.NewRequest("POST", url.String(), bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("X-API-Key", *s.APIKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}()

	switch resp.StatusCode {
	case 201:
		return nil
	case 429, 503:
		return ErrRateExceeded
	case 402:
		return ErrPaymentRequired
	case 403:
		return ErrUnauthorized
	default:
		return fmt.Errorf(
			"request failed status=%d expected=%d",
			resp.StatusCode,
			http.StatusCreated)
	}
}
